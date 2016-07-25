package static

import "net/http"

type handlerStaticContent struct {
	content    string
	returnCode int
}

func NewHandlerStaticContent(content string) *handlerStaticContent {
	return NewHandlerStaticContentReturnCode(content, http.StatusOK)
}

func NewHandlerStaticContentReturnCode(content string, returnCode int) *handlerStaticContent {
	h := new(handlerStaticContent)
	h.content = content
	h.returnCode = returnCode
	return h
}

func (t *handlerStaticContent) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.WriteHeader(t.returnCode)
	responseWriter.Write([]byte(t.content))
}
