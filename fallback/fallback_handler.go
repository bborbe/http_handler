package fallback

import (
	"net/http"

	"github.com/bborbe/http_handler_finder"
	"github.com/golang/glog"
)

type handler struct {
	handlerFinder handler_finder.HandlerFinder
	fallback      http.Handler
}

func New(handlerFinder handler_finder.HandlerFinder, fallbackHandler http.Handler) *handler {
	h := new(handler)
	h.handlerFinder = handlerFinder
	h.fallback = fallbackHandler
	return h
}

func (h *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	handler := h.handlerFinder.FindHandler(request)
	if handler != nil {
		glog.V(4).Info("handler found, use handler")
		handler.ServeHTTP(responseWriter, request)
		return
	}
	if h.fallback != nil {
		glog.V(4).Info("no handler found, use fallback")
		h.fallback.ServeHTTP(responseWriter, request)
		return
	}
	glog.V(4).Info("no handler found and no fallback found")
}
