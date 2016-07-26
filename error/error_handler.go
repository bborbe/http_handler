package error

import (
	"net/http"

	"encoding/json"

	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

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
	logger.Debug("handle error")

	var data struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}
	data.Message = o.message
	data.Status = o.status
	logger.Debugf("set status: %d", o.status)
	responseWriter.WriteHeader(o.status)
	responseWriter.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(responseWriter).Encode(&data); err != nil {
		logger.Warnf("render failureRenderer failed! %v", err)
	}
}
