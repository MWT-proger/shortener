package auth

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/MWT-proger/shortener/configs"
	"github.com/MWT-proger/shortener/internal/shortener/logger"
)

// AuthCookieMiddleware(next http.Handler) http.Handler — middleware-для входящих HTTP-запросов.
// Выдаёт пользователю симметрично подписанную куку, содержащую уникальный идентификатор пользователя,
// если такой куки не существует или она не проходит проверку подлинности.
func AuthCookieMiddleware(next http.Handler, conf configs.Config) http.Handler {
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
			UserID = getUserID(conf, tokenString)

			if UserID == uuid.Nil {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

		} else {
			UserID = uuid.New()

			tokenString, err = buildJWTString(conf, UserID)

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

// CheckIncludedSubNetMiddleware — middleware-для входящих HTTP-запросов.
// Проверяет что переданный в заголовке запроса X-Real-IP IP-адрес клиента входит в доверенную подсеть,
// в противном случае возвращает статус ответа 403 Forbidden.
func CheckIPIncludedSubNetMiddleware(next http.Handler, conf configs.Config) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if conf.Auth.TrustedSubNet == "" {
			http.Error(w, "IP-адрес клиента не входит в доверенную подсеть", http.StatusForbidden)
			return
		}

		ip, err := resolveIP(r)
		fmt.Println(ip)
		if err != nil {
			http.Error(w, "IP-адрес клиента не входит в доверенную подсеть", http.StatusForbidden)
			return
		}

		if !strings.Contains(conf.Auth.TrustedSubNet, ip.String()) {
			http.Error(w, "IP-адрес клиента не входит в доверенную подсеть", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func resolveIP(r *http.Request) (net.IP, error) {
	ipStr := r.Header.Get("X-Real-IP")
	ip := net.ParseIP(ipStr)

	if ip != nil {
		return ip, nil
	}

	ips := r.Header.Get("X-Forwarded-For")
	ipStrs := strings.Split(ips, ",")
	ipStr = ipStrs[0]
	ip = net.ParseIP(ipStr)

	if ip != nil {
		return nil, fmt.Errorf("failed parse ip from http header")
	}
	return ip, nil
}
