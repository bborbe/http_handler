package log

import (
	"net/http"
	"time"

	"github.com/golang/glog"
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
	start := time.Now()
	glog.V(4).Infof("%s %s", request.Method, request.URL.Path)
	h.subhandler.ServeHTTP(responseWriter, request)
	end := time.Now()
	glog.V(4).Infof("%s %s takes %dms", request.Method, request.URL.Path, end.Sub(start)/time.Millisecond)
}
