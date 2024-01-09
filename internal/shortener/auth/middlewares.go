package auth

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/MWT-proger/shortener/internal/shortener/logger"
)

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
			newCookie  = http.Cookie{Name: nameCookie}
			token, err = r.Cookie(nameCookie)
		)

		if err == nil {
			tokenString = token.Value
			UserID = getUserID(tokenString)

			if UserID == uuid.Nil {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

		} else {
			UserID = uuid.New()

			tokenString, err = buildJWTString(UserID)

			if err != nil {
				logger.Log.Error(err.Error())
				http.Error(w, "", http.StatusInternalServerError)
				return
			}

			newCookie.Value = tokenString
			http.SetCookie(w, &newCookie)
		}

		ctx = WithUserID(ctx, UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
