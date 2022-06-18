package middleware

import (
	"fmt"
	"net/http"
	"tp-tdl/token"

	"github.com/golang-jwt/jwt/v4"
)

func AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tkn := r.Header.Get("Auth-Token")

		if len(tkn) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, err := token.GetContent(tkn)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		fmt.Println(claims.(jwt.MapClaims)["user_id"].(string))
		r.Header.Set("user_id", claims.(jwt.MapClaims)["user_id"].(string))

		next.ServeHTTP(w, r)
	})
}
