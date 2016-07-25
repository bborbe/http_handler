package mux

import (
	"net/http"

	"github.com/bborbe/http_handler_finder"
)

type muxHandler struct {
	handlerFinder handler_finder.HandlerFinder
	errorHandler  http.Handler
}

func NewMuxHandler(handlerFinder handler_finder.HandlerFinder, errorHandler http.Handler) *muxHandler {
	m := new(muxHandler)
	m.handlerFinder = handlerFinder
	m.errorHandler = errorHandler
	return m
}

func (m *muxHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	handler := m.handlerFinder.FindHandler(request)
	if handler != nil {
		handler.ServeHTTP(responseWriter, request)
	} else {
		m.errorHandler.ServeHTTP(responseWriter, request)
	}
}
