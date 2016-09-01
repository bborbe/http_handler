package cacheheader

import (
	"net/http"
	"strings"
	"time"
)

var modtime = time.Now()
var extensionToCachingHeader = map[string]bool{
	"json": true,
	"gif":  true,
	"jpg":  true,
	"png":  true,
	"js":   true,
	"css":  true,
}

type cachingHeaderHandler struct {
	handler http.Handler
}

func New(handler http.Handler) *cachingHeaderHandler {
	h := new(cachingHeaderHandler)
	h.handler = handler
	return h
}

func (h *cachingHeaderHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	_, found := extensionToCachingHeader[getExtension(request.RequestURI)]
	if found {
		SetHeaderToResponse(responseWriter)
	}
	h.handler.ServeHTTP(responseWriter, request)
}

func getExtension(uri string) string {
	pos := strings.LastIndex(uri, ".")
	if pos == -1 {
		return ""
	}
	return uri[pos+1:]
}

// SetHeaderToResponse for cache
func SetHeaderToResponse(responseWriter http.ResponseWriter) {
	responseWriter.Header().Set("Cache-Control", "max-age=864000")
	responseWriter.Header().Set("Vary", "Accept-Encoding")
	responseWriter.Header().Set("Last-Modified", modtime.UTC().Format(http.TimeFormat))
}
