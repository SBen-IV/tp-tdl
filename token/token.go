package token

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("secret-key"))

const (
	TokenName = "auth-token"
	AuthKey   = "authorize"
)

func CreateToken(w http.ResponseWriter, r *http.Request, otherData map[string]string) error {
	session, err := Store.Get(r, TokenName)

	if err != nil {
		return err
	}

	session.Values[AuthKey] = true

	for key, value := range otherData {
		session.Values[key] = value
	}

	session.Save(r, w)

	return nil
}

func DestroyToken(w http.ResponseWriter, r *http.Request) error {
	session, err := Store.Get(r, TokenName)

	if err != nil {
		return err
	}

	session.Values[AuthKey] = false

	session.Save(r, w)

	return nil
}
