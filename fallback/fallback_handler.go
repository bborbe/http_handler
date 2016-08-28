package fallback

import (
	"net/http"

	"github.com/bborbe/http_handler_finder"
	"github.com/golang/glog"
)

type fallback struct {
	handlerFinder handler_finder.HandlerFinder
	fallback      http.Handler
}

func NewFallback(handlerFinder handler_finder.HandlerFinder, fallbackHandler http.Handler) *fallback {
	m := new(fallback)
	m.handlerFinder = handlerFinder
	m.fallback = fallbackHandler
	return m
}

func (m *fallback) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	handler := m.handlerFinder.FindHandler(request)
	if handler != nil {
		glog.V(2).Info("handler found, use handler")
		handler.ServeHTTP(responseWriter, request)
		return
	}
	if m.fallback != nil {
		glog.V(2).Info("no handler found, use fallback")
		m.fallback.ServeHTTP(responseWriter, request)
		return
	}
	glog.Info("no handler found and no fallback found")
}
