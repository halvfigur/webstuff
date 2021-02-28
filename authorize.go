package webstuff

import (
	"context"
	"net/http"
)

type Authorizer interface {
	Authorize(*http.Request) (interface{}, error)
}

// ContextAttrSession is the key associated with the current session attached
// to an authenticated request as a context value.
const ContextAttrSession = 0

// Authorize attemps to authorize a request and forwards the request to the next
// handler if successful.  If the request could not be authorized an
// appropriate error code is sent to the client.
func Authorize(auth Authorizer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, req *http.Request) {

			session, err := auth.Authorize(req)
			if err != nil {
				switch err {
				case ErrBadRequest:
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				case ErrUnauthorized:
					http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
					return
				case ErrExpired:
					http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
					return
				default:
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			}

			// All checks passed, forward request to next handler
			ctx := context.WithValue(req.Context(), ContextAttrSession, session)
			next.ServeHTTP(w, req.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
