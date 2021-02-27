package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/halvfigur/webstuff/controller"
	"github.com/halvfigur/webstuff/model"
	"github.com/halvfigur/webstuff/testutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockHandler struct {
	mock.Mock
}

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func TestAuthenticate(t *testing.T) {
	now := time.Date(2021, time.February, 7, 0, 0, 0, 0, time.UTC)

	toJson := func(session model.Session) []byte {
		b, _ := json.Marshal(session)
		return b
	}

	tests := []struct {
		name         string
		cookie       *http.Cookie
		decodedValue []byte
		decoderErr   error

		expectCode       int
		expectNextCalled bool
	}{
		{
			name:   "Session cookie missing",
			cookie: nil,

			expectCode:       http.StatusUnauthorized,
			expectNextCalled: false,
		},
		{
			name: "Failed to decode cookie",
			cookie: &http.Cookie{
				Name:  controller.CookieAttrSession,
				Value: "",
			},
			decoderErr: errors.New("failed to decode"),

			expectCode:       http.StatusBadRequest,
			expectNextCalled: false,
		},
		{
			name: "Non JSON cookie payload",
			cookie: &http.Cookie{
				Name:  controller.CookieAttrSession,
				Value: "",
			},

			expectCode:       http.StatusBadRequest,
			expectNextCalled: false,
		},
		{
			name: "Cookie expired",
			cookie: &http.Cookie{
				Name:  controller.CookieAttrSession,
				Value: "",
			},
			decodedValue: toJson(model.Session{
				Username: "testuser",
				Expires:  now.Add(-1 * time.Second),
			}),

			expectCode:       http.StatusUnauthorized,
			expectNextCalled: false,
		},
		{
			name: "Cookie accepted",
			cookie: &http.Cookie{
				Name:  controller.CookieAttrSession,
				Value: "",
			},
			decodedValue: toJson(model.Session{
				Username: "testuser",
				Expires:  now,
			}),

			expectCode:       http.StatusOK,
			expectNextCalled: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			timeNow = func() time.Time {
				return now
			}

			codec := new(testutil.MockSessionCodec)
			codec.On(testutil.MockSessionCodecDecode, mock.Anything).Return(test.decodedValue, test.decoderErr)

			authenticate := Authenticate(codec)
			m := new(mockHandler)
			m.On("ServeHTTP", mock.Anything, mock.Anything).Return()
			handler := authenticate(m.ServeHTTP)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://url.com", nil)
			if test.cookie != nil {
				req.AddCookie(test.cookie)
			}

			handler(w, req)
			resp := w.Result()

			require.Equal(t, resp.StatusCode, test.expectCode)

			if test.expectNextCalled {
				m.AssertCalled(t, "ServeHTTP")
			} else {
				m.AssertNotCalled(t, "ServeHTTP")
			}
		})
	}
}
