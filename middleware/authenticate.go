package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/halvfigur/webstuff/controller"
	"github.com/halvfigur/webstuff/model"
)

// ContextAttrSession is the key associated with the current session attached
// to an authenticated request as a context value.
const ContextAttrSession = 0

// timeNow is an alias for the current time function, it's useful when writing unit tests
var timeNow = time.Now

// Authenticate authenticates a session by verifying that the session cookie
// "session" is valid.  If the session is valid the session is attached to the
// request context with the key middleware.ContextAttrSession and the request
// is forwareded to the next handler.
// If the session was not valid an appropriate error code is returned.
func Authenticate(codec controller.SessionCodec) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {

			// Do we have a session cookie?
			sessionCookie, err := req.Cookie(controller.CookieAttrSession)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			decodedCookie, err := codec.Decode(sessionCookie.Value)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			// Can we turn the cookie into a Session?
			var session model.Session
			if err := json.Unmarshal(decodedCookie, &session); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			// Check if session has expired
			if timeNow().After(session.Expires) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// All checks passed, forward request to next handler
			ctx := context.WithValue(req.Context(), ContextAttrSession, session)
			next(w, req.WithContext(ctx))
		}
	}
}
