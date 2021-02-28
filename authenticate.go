package webstuff

import (
	"io"
	"net/http"
	"time"
)

const CookieAttrSession = "session"
const SessionTimeout = 30 * time.Minute

type Authenticator interface {
	// Authenticate a session based on the payload available in r.
	// On success return a session and the session timeout.
	Authenticate(io.ReadCloser) (*http.Cookie, error)
}

// Authenticate accepts POST requests with a valid login request. Reply codes
// and conditions as follows.
//
// 405 Method not allowed
//	- If the method is not "POST"
// 400 Bad request
//	- If the request didn't contain a valid credentials model
// 401 Unauthorized
//	- If the credentials could not be authenticated
// 200 OK
//	- If the credentials were authenticated
//
// If successful a session cookie named "session" is attached to the reply.
func Authenticate(auth Authenticator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only POST requests are accepted
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		// Authenticate requests
		c, err := auth.Authenticate(r.Body)
		if err != nil {
			switch err {
			case ErrBadRequest:
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			case ErrUnauthorized:
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		// Set cookie, reply and we're done
		http.SetCookie(w, c)
		w.WriteHeader(http.StatusOK)
	}
}
