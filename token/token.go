package token

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

type UserData map[string]string

var store = sessions.NewCookieStore([]byte("secret-key"))

const (
	TokenName = "auth-token"
	AuthKey   = "authorize"
)

func CreateToken(w http.ResponseWriter, r *http.Request, otherData UserData) error {
	session, err := store.Get(r, TokenName)

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

func GetContent(r *http.Request, keys []string) (UserData, error) {
	data := UserData{}
	session, err := store.Get(r, TokenName)

	if err != nil {
		return data, err
	}

	ok := session.Values[AuthKey].(bool)
	if !ok {
		return data, errors.New("Unauthorized")
	}

	for _, key := range keys {
		value := session.Values[key].(string)
		fmt.Println(key, value)
		data[key] = value
	}

	return data, nil
}

func DestroyToken(w http.ResponseWriter, r *http.Request) error {
	session, err := store.Get(r, TokenName)

	if err != nil {
		return err
	}

	session.Values[AuthKey] = false

	session.Save(r, w)

	return nil
}
