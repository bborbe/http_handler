package debug

import (
	"net/http"
	"time"

	"github.com/golang/glog"
)

type handler struct {
	subhandler http.Handler
}

func New(subhandler http.Handler) *handler {
	m := new(handler)
	m.subhandler = subhandler
	return m
}

func (m *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	start := time.Now()
	defer glog.V(2).Infof("%s %s takes %dms", request.Method, request.RequestURI, time.Now().Sub(start)/time.Millisecond)

	glog.V(2).Infof("request %v: ", request)
	m.subhandler.ServeHTTP(responseWriter, request)
	glog.V(2).Infof("response %v: ", responseWriter)
}
