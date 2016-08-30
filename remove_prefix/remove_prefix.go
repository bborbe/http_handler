package remove_prefix

import (
	"net/http"
	"strings"
	"github.com/golang/glog"
)

type handler struct {
	prefix     string
	subhandler http.HandlerFunc
}

func New(prefix string, subhandler http.HandlerFunc) *handler {
	h := new(handler)
	h.prefix = prefix
	h.subhandler = subhandler
	return h
}

func (h *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	glog.V(2).Infof("remove prefix %v from request", h.prefix)
	if strings.HasPrefix(request.RequestURI, h.prefix) {
		request.RequestURI = request.RequestURI[len(h.prefix):]
	}
	if strings.HasPrefix(request.URL.Path, h.prefix) {
		request.URL.Path = request.URL.Path[len(h.prefix):]
	}
	h.subhandler(responseWriter, request)
}
