package redirect

import (
	"net/http"

	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type handler struct {
	target string
	status int
}

func New(target string) *handler {
	h := new(handler)
	h.target = target
	h.status = http.StatusMovedPermanently
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger.Debugf("redirect to %s %d", h.target, h.status)
	http.Redirect(resp, req, h.target, h.status)
}
