package webstuff

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockAuthenticator struct {
	mock.Mock
}

const mockAuthenticatorAuthenticate = "Authenticate"

func (m *mockAuthenticator) Authenticate(r io.ReadCloser) (*http.Cookie, error) {
	args := m.Called(r)

	return args.Get(0).(*http.Cookie), args.Error(1)
}

func TestAuthenticate(t *testing.T) {
	tests := []struct {
		name   string
		method string
		body   io.Reader

		cookie *http.Cookie
		err    error

		expectCode int
	}{
		{
			name:   "Invalid method GET",
			method: http.MethodGet,

			expectCode: http.StatusMethodNotAllowed,
		},
		{
			name:   "Bad request",
			method: http.MethodPost,
			err:    ErrBadRequest,

			expectCode: http.StatusBadRequest,
		},
		{
			name:   "Unauthorized",
			method: http.MethodPost,
			err:    ErrUnauthorized,

			expectCode: http.StatusUnauthorized,
		},
		{
			name:   "Internal server error",
			method: http.MethodPost,
			err:    errors.New("generic error"),

			expectCode: http.StatusInternalServerError,
		},
		{
			name:   "Authentication succeeded",
			method: http.MethodPost,
			cookie: &http.Cookie{},

			expectCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, "http://url.com", test.body)
			w := httptest.NewRecorder()

			auth := new(mockAuthenticator)
			auth.On(mockAuthenticatorAuthenticate, mock.Anything).Return(test.cookie, test.err)

			handler := Authenticate(auth)
			handler(w, req)
			resp := w.Result()

			require.Equal(t, test.expectCode, resp.StatusCode)
		})
	}
}
