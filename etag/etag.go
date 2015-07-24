//Package etag implements an etag handler middleware for Negroni.
package etag

import (
	"crypto/md5"
	"fmt"
	"net/http"
)

type etagResponseWriter struct {
	writer http.ResponseWriter
	req    *http.Request
	code   int
}

func (erw *etagResponseWriter) Write(b []byte) (int, error) {
	etag := fmt.Sprintf("%x", md5.Sum(b))
	erw.Header().Set("ETag", etag)
	if erw.r.Header.Get("If-None-Match") == etag {
		erw.writer.WriteHeader(304)
		return erw.writer.Write(nil)
	}
	erw.writer.WriteHeader(erw.code)
	return erw.writer.Write(b)
}

func (erw *etagResponseWriter) Header() http.Header {
	return erw.writer.Header()
}

func (erw *etagResponseWriter) WriteHeader(code int) {
	erw.code = code
}

type handler struct{}

func Etag() *handler {
	return &handler{}
}

func (h *handler) ServeHTTP(writer http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	erw := &etagResponseWriter{writer, 200, req}
	next(erw, req)
}
