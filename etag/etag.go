// Package etag implements an etag handler middleware for Negroni.
package etag

import (
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
)

type etagResponseWriter struct {
	nrw negroni.ResponseWriter
	r   *http.Request
}

func (erw etagResponseWriter) Write(b []byte) (int, error) {
	etag := fmt.Sprintf("%x", md5.Sum(b))
	erw.Header().Add("etag", etag)
	if erw.r.Header.Get("If-None-Match") == etag {
		erw.nrw.WriteHeader(304)
		return erw.nrw.Write(nil)
	}
	return erw.nrw.Write(b)
}
func (erw etagResponseWriter) Header() http.Header {
	return erw.nrw.Header()
}
func (erw etagResponseWriter) WriteHeader(code int) {
	erw.nrw.WriteHeader(code)
}

type handler struct{}

func Etag() *handler {
	return &handler{}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	nrw := negroni.NewResponseWriter(w)
	erw := etagResponseWriter{nrw, r}

	next(erw, r)
}
