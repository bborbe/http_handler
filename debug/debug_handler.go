package debug

import (
	"net/http"
	"time"

	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

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
	defer logger.Debugf("%s %s takes %dms", request.Method, request.RequestURI, time.Now().Sub(start)/time.Millisecond)

	logger.Debugf("request %v: ", request)
	m.subhandler.ServeHTTP(responseWriter, request)
	logger.Debugf("response %v: ", responseWriter)
}
