package middleware

import (
	"net/http"
	"tp-tdl/token"
)

func AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := token.Store.Get(r, "auth-token")

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ok := session.Values["authorize"].(bool)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r.Header.Set("user_id", session.Values["user_id"].(string))
		r.Header.Set("username", session.Values["username"].(string))

		next.ServeHTTP(w, r)
	})
}
