package middleware

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// проверяем, что клиент отправил сжатые данные
			contentEncoding := r.Header.Get("Content-Encoding")
			if strings.Contains(contentEncoding, "gzip") {
				gr, err := gzip.NewReader(r.Body)
				if err != nil {
					fmt.Println("ERROR: " + err.Error())
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				defer func(gr *gzip.Reader) {
					_ = gr.Close()
				}(gr)
				r.Body = gr
			}

			// проверяем, что клиент умеет принимать сжатые данные
			acceptEncoding := r.Header.Get("Accept-Encoding")
			if strings.Contains(acceptEncoding, "gzip") {
				gw := NewGzipWriter(w)
				defer gw.Close()
				next.ServeHTTP(gw, r)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}(next)
}

type GzipWriter struct {
	gw *gzip.Writer
	http.ResponseWriter
}

func NewGzipWriter(w http.ResponseWriter) *GzipWriter {
	gw := gzip.NewWriter(w)
	w.Header().Set("Content-Encoding", "gzip")
	return &GzipWriter{gw, w}
}

func (w *GzipWriter) Write(b []byte) (int, error) {
	return w.gw.Write(b)
}

func (w *GzipWriter) Close() {
	_ = w.gw.Close()
}
