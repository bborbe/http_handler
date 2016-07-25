package auth_basic

import (
	"fmt"
	"net/http"

	"github.com/bborbe/http/header"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type Check func(username string, password string) (bool, error)

type handler struct {
	handler http.HandlerFunc
	check   Check
	realm   string
}

func New(subhandler http.HandlerFunc, check Check, realm string) *handler {
	h := new(handler)
	h.handler = subhandler
	h.check = check
	h.realm = realm
	return h
}

func (h *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	logger.Debugf("check basic auth")
	if err := h.serveHTTP(responseWriter, request); err != nil {
		responseWriter.Header().Add("WWW-Authenticate", fmt.Sprintf("Basic realm=\"%s\"", h.realm))
		responseWriter.WriteHeader(http.StatusUnauthorized)
	}
}

func (h *handler) serveHTTP(responseWriter http.ResponseWriter, request *http.Request) error {
	logger.Debugf("check basic auth")
	user, pass, err := header.ParseAuthorizationBasisHttpRequest(request)
	if err != nil {
		logger.Debugf("parse header failed: %v", err)
		return err
	}
	result, err := h.check(user, pass)
	if err != nil {
		logger.Debugf("check auth failed: %v", err)
		return err
	}
	if !result {
		logger.Debugf("auth failed")
		return fmt.Errorf("auth failed")
	}
	h.handler(responseWriter, request)
	return nil
}
