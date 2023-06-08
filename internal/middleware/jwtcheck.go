package middleware

import (
	"github.com/go-chi/jwtauth/v5"
	"log"
)

func NewJWT() *jwtauth.JWTAuth {
	tokenAuth := jwtauth.New("HS256", []byte("june2023"), nil)
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": 123})
	log.Printf("\t\twarning: test jwt is: %s\n", tokenString)
	return tokenAuth
}
