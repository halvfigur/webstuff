package webstuff

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockHandler struct {
	mock.Mock
}

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

type mockAuthorizer struct {
	mock.Mock
}

const mockAuthorizerAuthorize = "Authorize"

func (m *mockAuthorizer) Authorize(req *http.Request) (interface{}, error) {
	args := m.Called(req)

	return args.Get(0), args.Error(1)
}

func TestAuthorize(t *testing.T) {
	tests := []struct {
		name    string
		session interface{}
		err     error

		expectCode       int
		expectNextCalled bool
	}{
		{
			name: "Bad request",
			err:  ErrBadRequest,

			expectCode:       http.StatusBadRequest,
			expectNextCalled: false,
		},
		{
			name: "Unautorized b/c invalid credentials",
			err:  ErrUnauthorized,

			expectCode:       http.StatusUnauthorized,
			expectNextCalled: false,
		},
		{
			name: "Unautorized b/c expired",
			err:  ErrUnauthorized,

			expectCode:       http.StatusUnauthorized,
			expectNextCalled: false,
		},
		{
			name: "Internal server error",
			err:  errors.New("some error"),

			expectCode:       http.StatusInternalServerError,
			expectNextCalled: false,
		},
		{
			name:    "Authorization success",
			session: 1,

			expectCode:       http.StatusOK,
			expectNextCalled: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			authorizer := new(mockAuthorizer)
			authorizer.On(mockAuthorizerAuthorize, mock.Anything).Return(test.session, test.err)
			authenticate := Authorize(authorizer)

			m := new(mockHandler)
			m.On("ServeHTTP", mock.Anything, mock.Anything).Return()
			handler := authenticate(m)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://url.com", nil)

			handler.ServeHTTP(w, req)
			resp := w.Result()

			require.Equal(t, resp.StatusCode, test.expectCode)

			if test.expectNextCalled {
				m.AssertCalled(t, "ServeHTTP", w, mock.Anything)
			} else {
				m.AssertNotCalled(t, "ServeHTTP")
			}
		})
	}
}
