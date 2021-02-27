package controller

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/halvfigur/webstuff/model"
	"github.com/halvfigur/webstuff/testutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLoginController(t *testing.T) {
	tests := []struct {
		name   string
		method string
		body   io.Reader

		storeUser  model.User
		storeError error

		encoderResult string
		encoderError  error

		expectCode int
	}{
		{
			name:   "Invalid method GET",
			method: http.MethodGet,

			expectCode: http.StatusMethodNotAllowed,
		},
		{
			name:   "Invalid JSON",
			method: http.MethodPost,
			body:   bytes.NewBufferString(`{"key": "value"`),

			expectCode: http.StatusBadRequest,
		},
		{
			name:   "Incorrect JSON",
			method: http.MethodPost,
			body:   bytes.NewBufferString("{..."),

			expectCode: http.StatusBadRequest,
		},
		{
			name:       "Store failure",
			method:     http.MethodPost,
			body:       bytes.NewBufferString(`{"username": "thedude", "password": "whiterussian"}`),
			storeError: errors.New("store error"),

			expectCode: http.StatusUnauthorized,
		},
		{
			name:      "Authentication failure",
			method:    http.MethodPost,
			body:      bytes.NewBufferString(`{"username": "thedude", "password": "whiterussian"}`),
			storeUser: model.User{Password: "otherpassword"},

			expectCode: http.StatusUnauthorized,
		},
		{
			name:         "Encoder failure",
			method:       http.MethodPost,
			body:         bytes.NewBufferString(`{"username": "thedude", "password": "whiterussian"}`),
			storeUser:    model.User{Password: "whiterussian"},
			encoderError: errors.New("encoder error"),

			expectCode: http.StatusInternalServerError,
		},
		{
			name:          "Authentication succeeded",
			method:        http.MethodPost,
			body:          bytes.NewBufferString(`{"username": "thedude", "password": "whiterussian"}`),
			storeUser:     model.User{Password: "whiterussian"},
			encoderResult: "encryptedblob",

			expectCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, "http://url.com", test.body)
			w := httptest.NewRecorder()

			store := new(testutil.MockStore)
			store.On(testutil.MockStoreUser, mock.Anything).Return(test.storeUser, test.storeError)
			codec := new(testutil.MockSessionCodec)
			codec.On(testutil.MockSessionCodecEncode, mock.Anything).Return(test.encoderResult, test.encoderError)

			handler := Login(store, codec)
			handler(w, req)
			resp := w.Result()

			require.Equal(t, test.expectCode, resp.StatusCode)
		})
	}
}
