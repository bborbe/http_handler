package nocacheheader

import (
	"net/http"
)

type handler struct {
	subhandler http.Handler
}

func New(subhandler http.Handler) *handler {
	h := new(handler)
	h.subhandler = subhandler
	return h
}

func (h *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	SetHeaderToResponse(responseWriter)
	h.subhandler.ServeHTTP(responseWriter, request)
}

// SetHeaderToResponse for nocache
func SetHeaderToResponse(responseWriter http.ResponseWriter) {
	responseWriter.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	responseWriter.Header().Set("Pragma", "no-cache")
	responseWriter.Header().Set("Expires", "0")
}
