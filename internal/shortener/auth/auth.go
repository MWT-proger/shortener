package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/MWT-proger/shortener/configs"
	"github.com/MWT-proger/shortener/internal/shortener/logger"
	"github.com/MWT-proger/shortener/internal/shortener/request"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

const NameCookie = "token"

// AuthCookieMiddleware(next http.Handler) http.Handler — middleware-для входящих HTTP-запросов.
// Выдаёт пользователю симметрично подписанную куку, содержащую уникальный идентификатор пользователя,
// если такой куки не существует или она не проходит проверку подлинности.
func AuthCookieMiddleware(next http.Handler) http.Handler {
	// получаем Handler приведением типа http.HandlerFunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			UserID      uuid.UUID
			tokenString string
			isNewCookie bool

			ctx        = r.Context()
			newCookie  = http.Cookie{Name: NameCookie}
			token, err = r.Cookie(NameCookie)
		)

		if err == nil {
			tokenString = token.Value
			UserID, _ = GetUserID(tokenString)
		}

		if UserID == uuid.Nil {
			isNewCookie = true
			UserID = uuid.New()

			tokenString, err = BuildJWTString(UserID)

			if err != nil {
				logger.Log.Error(err.Error())
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		}

		if isNewCookie {
			newCookie.Value = tokenString
			http.SetCookie(w, &newCookie)
		}

		ctx = request.WithUserID(ctx, UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// ValidAuthCookieMiddleware(next http.Handler) http.Handler — middleware-для входящих HTTP-запросов.
// Проверяет у пользователя куку, содержащую уникальный идентификатор пользователя,
// если такой куки не существует или она не проходит проверку подлинности то возвращает статус неавторизованного
func ValidAuthCookieMiddleware(next http.Handler) http.Handler {
	// получаем Handler приведением типа http.HandlerFunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			UserID      uuid.UUID
			tokenString string

			ctx        = r.Context()
			token, err = r.Cookie(NameCookie)
		)

		if err == nil {
			tokenString = token.Value
			UserID, _ = GetUserID(tokenString)
		}

		if UserID == uuid.Nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		ctx = request.WithUserID(ctx, UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// BuildJWTString(UserID uuid.UUID) (string, error) создаёт токен для пользователя с UserID
// и возвращает его в виде строки в случае успеха
func BuildJWTString(UserID uuid.UUID) (string, error) {
	conf := configs.GetConfig()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(conf.Auth.TokenExp)),
		},

		UserID: UserID,
	})

	tokenString, err := token.SignedString([]byte(conf.Auth.SecretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetUserID(tokenString string) (uuid.UUID, error) Проверяет токен
// и в случае успеха возвращает из полезной нагрузки UserID
func GetUserID(tokenString string) (uuid.UUID, error) {

	claims := &Claims{}
	conf := configs.GetConfig()

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(conf.Auth.SecretKey), nil
		})

	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		logger.Log.Debug("Token is not valid")
		return uuid.Nil, err
	}

	logger.Log.Debug("Token is valid")
	return claims.UserID, err
}