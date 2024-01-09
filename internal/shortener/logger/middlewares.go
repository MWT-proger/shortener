package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
)

// RequestLoggerMiddleware — middleware-логер для входящих HTTP-запросов.
func RequestLoggerMiddleware(next http.Handler) http.Handler {
	// получаем Handler приведением типа http.HandlerFunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {

			Log.Info("got incoming HTTP request",
				StringField("method", r.Method),
				StringField("path", r.URL.Path),
				IntField("status", ww.Status()),
				IntField("length", ww.BytesWritten()),
				DurationField("time", time.Since(timeStart)),
			)
		}()

		next.ServeHTTP(ww, r)
	})
}
