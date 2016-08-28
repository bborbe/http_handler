package json

import (
	"net/http"

	"encoding/json"
	"reflect"

	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type jsonHandler struct {
	m interface{}
}

func NewJsonHandler(m interface{}) *jsonHandler {
	h := new(jsonHandler)
	h.m = m
	return h
}

func (m *jsonHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	glog.V(2).Info("write json")
	glog.V(2).Infof("object to convert %v", m.m)
	b, err := json.Marshal(m.m)
	if err != nil {
		glog.V(2).Infof("Marshal json failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(responseWriter, request)
		return
	}
	glog.V(2).Infof("json string %s", string(b))
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)

	glog.V(2).Infof("object type %v", reflect.TypeOf(m.m).Kind())
	if reflect.TypeOf(m.m).Kind() == reflect.Slice && string(b) == "null" {
		responseWriter.Write([]byte("[]"))
	} else {
		responseWriter.Write(b)
	}

}
