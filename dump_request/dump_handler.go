package dump

import (
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
	"net/http"
	"net/http/httputil"
	"fmt"
)

type handler struct {
}

func New() *handler {
	h := new(handler)
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("dump request started")
	if err := h.serveHTTP(resp, req); err != nil {
		glog.V(1).Infof("dump request failed: %v", err)
		e := error_handler.NewMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
		return
	}
	glog.V(2).Infof("dump request finished")
}

func (h *handler) serveHTTP(responseWriter http.ResponseWriter, request *http.Request) error {
	content, err := httputil.DumpRequest(request, true)
	if err != nil {
		glog.V(2).Infof("dump request failed: %v", err)
		return err
	}
	responseWriter.Write(content)
	fmt.Fprintf(responseWriter, "RemoteAddr: %v\n", request.RemoteAddr)
	glog.V(1).Info(string(content))
	return nil
}
