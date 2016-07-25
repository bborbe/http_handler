package json

import (
	"net/http"

	"encoding/json"
	"reflect"

	error_handler "github.com/bborbe/http_handler/error"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type jsonHandler struct {
	m interface{}
}

func NewJsonHandler(m interface{}) *jsonHandler {
	h := new(jsonHandler)
	h.m = m
	return h
}

func (m *jsonHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	logger.Debug("write json")
	logger.Debugf("object to convert %v", m.m)
	b, err := json.Marshal(m.m)
	if err != nil {
		logger.Debugf("Marshal json failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(responseWriter, request)
		return
	}
	logger.Debugf("json string %s", string(b))
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)

	logger.Debugf("object type %v", reflect.TypeOf(m.m).Kind())
	if reflect.TypeOf(m.m).Kind() == reflect.Slice && string(b) == "null" {
		responseWriter.Write([]byte("[]"))
	} else {
		responseWriter.Write(b)
	}

}
