package remove_prefix

import (
	"net/http"
	"strings"
)

type handler struct {
	prefix  string
	handler http.HandlerFunc
}

func New(prefix string, subhandler http.HandlerFunc) *handler {
	h := new(handler)
	h.prefix = prefix
	h.handler = subhandler
	return h
}

func (h *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if strings.HasPrefix(request.RequestURI, h.prefix) {
		request.RequestURI = request.RequestURI[len(h.prefix):]
	}
	if strings.HasPrefix(request.URL.Path, h.prefix) {
		request.URL.Path = request.URL.Path[len(h.prefix):]
	}
	h.handler(responseWriter, request)
}
