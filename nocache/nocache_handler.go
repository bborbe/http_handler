package nocache

import (
	"net/http"
)

type noCacheHandler struct {
	handler http.Handler
}

func New(handler http.Handler) *noCacheHandler {
	m := new(noCacheHandler)
	m.handler = handler
	return m
}

func (m *noCacheHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	responseWriter.Header().Set("Pragma", "no-cache")
	responseWriter.Header().Set("Expires", "0")
	m.handler.ServeHTTP(responseWriter, request)
}
