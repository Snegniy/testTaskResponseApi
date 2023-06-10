package jwt

import (
	"github.com/Snegniy/testTaskResponseApi/pkg/logger"
	"github.com/go-chi/jwtauth/v5"
	"go.uber.org/zap"
)

func NewJWT() *jwtauth.JWTAuth {
	tokenAuth := jwtauth.New("HS256", []byte("june2023"), nil)
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": 123})
	logger.Info("test jwt", zap.String("token", tokenString))
	return tokenAuth
}
