package adapter

import (
	"net/http"
)

type handler struct {
	handler http.HandlerFunc
}

func New(subhandler http.HandlerFunc) *handler {
	h := new(handler)
	h.handler = subhandler
	return h
}

func (h *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	h.handler(responseWriter, request)
}
