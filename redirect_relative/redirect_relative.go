package redirect_relative

import (
	"net/http"

	"fmt"

	"github.com/golang/glog"
)

type handler struct {
	path   string
	status int
}

func New(target string) *handler {
	h := new(handler)
	h.path = target
	h.status = http.StatusMovedPermanently
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	target := fmt.Sprintf("%v%v", req.URL.Path, h.path)
	glog.V(4).Infof("redirect to %s %d", target, h.status)
	http.Redirect(resp, req, target, h.status)
}
