package log

import (
	"net/http"
	"time"

	"github.com/golang/glog"
)

type logHandler struct {
	handler http.Handler
}

func NewLogHandler(handler http.Handler) *logHandler {
	m := new(logHandler)
	m.handler = handler
	return m
}

func (m *logHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	start := time.Now()
	glog.V(2).Infof("%s %s", request.Method, request.RequestURI)
	m.handler.ServeHTTP(responseWriter, request)
	end := time.Now()
	glog.V(2).Infof("%s %s takes %dms", request.Method, request.RequestURI, end.Sub(start)/time.Millisecond)
}
