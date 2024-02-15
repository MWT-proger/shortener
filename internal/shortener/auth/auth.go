package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/MWT-proger/shortener/configs"
	"github.com/MWT-proger/shortener/internal/shortener/logger"
)

// claims сущность пользователя для внутренних операций.
type claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

// nameCookie константа nameCookie - ключ в cookie для токена авторизации.
const nameCookie = "token"
const gRPCnameKey = "token"

// buildJWTString(UserID uuid.UUID) (string, error) создаёт токен для пользователя с UserID
// и возвращает его в виде строки в случае успеха
func buildJWTString(conf configs.Config, UserID uuid.UUID) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{},
		UserID:           UserID,
	})

	tokenString, err := token.SignedString([]byte(conf.Auth.SecretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// getUserID(tokenString string) (uuid.UUID, error) Проверяет токен
// и в случае успеха возвращает из полезной нагрузки UserID
func getUserID(conf configs.Config, tokenString string) uuid.UUID {

	claims := &claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(conf.Auth.SecretKey), nil
		})

	if err != nil {
		logger.Log.Error(err.Error())
		return uuid.Nil
	}

	if !token.Valid {
		logger.Log.Debug("Token is not valid")
		return uuid.Nil
	}

	logger.Log.Debug("Token is valid")
	return claims.UserID
}
