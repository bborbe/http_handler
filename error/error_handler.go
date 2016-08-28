package error

import (
	"net/http"

	"encoding/json"

	"github.com/golang/glog"
)

type object struct {
	status  int
	message string
}

func NewError(status int) *object {
	return NewErrorMessage(status, http.StatusText(status))
}

func NewErrorMessage(status int, message string) *object {
	o := new(object)
	o.status = status
	o.message = message
	return o
}

func (o *object) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	glog.V(2).Info("handle error")

	var data struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}
	data.Message = o.message
	data.Status = o.status
	glog.V(2).Infof("set status: %d", o.status)
	responseWriter.WriteHeader(o.status)
	responseWriter.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(responseWriter).Encode(&data); err != nil {
		glog.Warningf("render failureRenderer failed! %v", err)
	}
}
