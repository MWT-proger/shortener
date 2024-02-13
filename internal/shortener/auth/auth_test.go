package auth

import (
	"testing"

	"github.com/MWT-proger/shortener/configs"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBuildJWTString(t *testing.T) {

	testCases := []struct {
		name   string
		userID uuid.UUID
	}{
		{name: "Тест 1 - успешный тест", userID: uuid.New()},
	}
	conf := configs.Config{Auth: configs.AuthConfig{SecretKey: "TestKey"}}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			claims := &claims{}

			tokenString, err := buildJWTString(conf, tt.userID)

			assert.NoError(t, err, "Ошибка при генерации токена")
			assert.NotNil(t, tokenString, "Токен пустой")

			token, err := jwt.ParseWithClaims(tokenString, claims,
				func(t *jwt.Token) (interface{}, error) {
					return []byte(conf.Auth.SecretKey), nil
				})

			assert.NoError(t, err, "Ошибка при чтение токена")
			assert.True(t, token.Valid, "Не валидный токен")

			assert.Equal(t, tt.userID, claims.UserID, "ID пользователя не совпадают")

		})
	}

}

func BenchmarkBuildJWTString(b *testing.B) {
	conf := configs.Config{Auth: configs.AuthConfig{SecretKey: "TestKey"}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buildJWTString(conf, uuid.New())
	}

}

func TestGetUserID(t *testing.T) {

	testCases := []struct {
		name   string
		userID uuid.UUID
	}{
		{name: "Тест 1 - успешный тест", userID: uuid.New()},
	}
	conf := configs.Config{Auth: configs.AuthConfig{SecretKey: "TestKey"}}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			tokenString, _ := buildJWTString(conf, tt.userID)

			userID := getUserID(conf, tokenString)

			assert.Equal(t, tt.userID, userID, "ID пользователя не совпадают")
		})
	}
}

func BenchmarkGetUserID(b *testing.B) {
	conf := configs.Config{Auth: configs.AuthConfig{SecretKey: "TestKey"}}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer() // останавливаем таймер
		tokenString, _ := buildJWTString(conf, uuid.New())
		b.StartTimer() // возобновляем таймер

		getUserID(conf, tokenString)
	}

}
