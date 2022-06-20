package token

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/sessions"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("secret-key")

var Store = sessions.NewCookieStore([]byte("secret-key"))

func CreateToken(user_id string) (string, error) {
	claims := Claims{
		UserID: user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "The Blues",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)

	return tokenStr, err
}

func GetContent(token string) (jwt.Claims, error) {
	tkn, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	return tkn.Claims, err
}
