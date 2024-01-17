package gzip

import (
	"net/http"
	"strings"
)

// GzipMiddleware - определяет запрос на сжатие и при необходимости возвращает сжатый ответ.
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ow             http.ResponseWriter
			acceptEncoding = r.Header.Get("Accept-Encoding")
			supportsGzip   = strings.Contains(acceptEncoding, "gzip")
		)

		if supportsGzip {

			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		} else {
			ow = w
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		if strings.Contains(contentEncoding, "gzip") {

			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(ow, r)
	})
}
