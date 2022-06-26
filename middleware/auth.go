package middleware

import (
	"net/http"
	"tp-tdl/token"
)

func AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := token.GetContent(r, []string{"user_id", "username"})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r.Header.Set("user_id", data["user_id"])
		r.Header.Set("username", data["username"])

		next.ServeHTTP(w, r)
	})
}
