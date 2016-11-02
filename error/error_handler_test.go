package error

import (
	"net/http"
	"testing"

	. "github.com/bborbe/assert"
	server_mock "github.com/bborbe/http/mock"
)

func TestImplementsRequestHandler(t *testing.T) {
	r := New(http.StatusNotFound)
	var i (*http.Handler) = nil
	if err := AssertThat(r, Implements(i).Message("check implements http.Handler")); err != nil {
		t.Fatal(err)
	}
}

func TestContent(t *testing.T) {
	handler := New(http.StatusNotFound)
	responseWriter := server_mock.NewHttpResponseWriterMock()
	request, err := server_mock.NewHttpRequestMock("http://www.example.com/foobar")
	if err != nil {
		t.Error(err)
	}
	handler.ServeHTTP(responseWriter, request)
	if err := AssertThat(responseWriter.Status(), Is(http.StatusNotFound).Message("check status")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(responseWriter.String(), Is("{\"status\":404,\"message\":\"Not Found\"}\n").Message("check content")); err != nil {
		t.Fatal(err)
	}
}
