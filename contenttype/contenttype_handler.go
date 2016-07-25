package contenttype

import (
	"net/http"
	"strings"
)

var ExtensionToContentType = map[string]string{
	"json": "application/json",
	"html": "text/html",
	"gif":  "image/gif",
	"jpg":  "image/jpeg",
	"png":  "image/png",
	"js":   "application/javascript",
	"css":  "text/css",
}

type contentTypeHandler struct {
	handler http.Handler
}

func NewContentTypeHandler(handler http.Handler) *contentTypeHandler {
	h := new(contentTypeHandler)
	h.handler = handler
	return h
}

func (h *contentTypeHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	contentType, found := ExtensionToContentType[getExtension(request.RequestURI)]
	if found {
		responseWriter.Header().Set("Content-Type", contentType)
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
