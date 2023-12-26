package auth

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/MWT-proger/shortener/configs"
	"github.com/MWT-proger/shortener/internal/shortener/logger"
	"github.com/MWT-proger/shortener/internal/shortener/request"
)

// Claims сущность пользователя
// BUG(Андрей): мб можно сделать приватным
type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

// константа NameCookie - ключ в cookie для токена авторизации
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

			ctx        = r.Context()
			newCookie  = http.Cookie{Name: NameCookie}
			token, err = r.Cookie(NameCookie)
		)

		if err == nil {
			tokenString = token.Value
			UserID = GetUserID(tokenString)

			if UserID == uuid.Nil {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

		} else {
			UserID = uuid.New()

			tokenString, err = BuildJWTString(UserID)

			if err != nil {
				logger.Log.Error(err.Error())
				http.Error(w, "", http.StatusInternalServerError)
				return
			}

			newCookie.Value = tokenString
			http.SetCookie(w, &newCookie)
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
		RegisteredClaims: jwt.RegisteredClaims{},
		UserID:           UserID,
	})

	tokenString, err := token.SignedString([]byte(conf.Auth.SecretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetUserID(tokenString string) (uuid.UUID, error) Проверяет токен
// и в случае успеха возвращает из полезной нагрузки UserID
func GetUserID(tokenString string) uuid.UUID {

	claims := &Claims{}
	conf := configs.GetConfig()

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
