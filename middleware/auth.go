package middleware

import (
	"net/http"
	"tp-tdl/token"
)

func AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/* 		tkn := r.Header.Get("Auth-Token")

		   		if len(tkn) == 0 {
		   			w.WriteHeader(http.StatusUnauthorized)
		   			return
		   		}

		   		claims, err := token.GetContent(tkn)

		   		if err != nil {
		   			w.WriteHeader(http.StatusUnauthorized)
		   			return
		   		}

		   		fmt.Println(claims.(jwt.MapClaims))
		   		r.Header.Set("user_id", claims.(jwt.MapClaims)["user_id"].(string)) */
		session, err := token.Store.Get(r, "auth-token")

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !session.Values["authorize"].(bool) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r.Header.Set("user_id", session.Values["user_id"].(string))

		next.ServeHTTP(w, r)
	})
}
