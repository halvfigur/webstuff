package controller

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/halvfigur/webstuff/model"
)

const CookieAttrSession = "session"
const SessionTimeout = 30 * time.Minute

type Store interface {
	User(username string) (model.User, error)
}

type Authenticate func(username, password string) error

func readJson(b io.ReadCloser, v interface{}) error {
	// The Server will close the request body. The ServeHTTP
	// Handler does not need to.

	// TODO consider using http.MaxBytesReader in order to limit the size of
	// incoming requests and save server resources.
	data, err := ioutil.ReadAll(b)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

// Login accepts POST requests with a valid JSON credentials model. Reply codes
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
func Login(store Store, codec SessionCodec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		// Extract credentials from request
		var cred model.Credentials
		err := readJson(r.Body, &cred)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Test password against stored user
		user, err := store.User(cred.Username)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if user.Password != cred.Password {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Create session
		session := model.Session{
			Username: cred.Username,
			Expires:  time.Now().Add(SessionTimeout),
		}

		// JSON encode session
		blob, err := json.Marshal(session)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Encrypt session
		encoded, err := codec.Encode(blob)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Create session cookie
		c := http.Cookie{
			Name:   CookieAttrSession,
			Value:  encoded,
			MaxAge: int(SessionTimeout.Seconds()),
		}

		// Set cookie, reply and we're done
		http.SetCookie(w, &c)
		w.WriteHeader(http.StatusOK)
	}
}
