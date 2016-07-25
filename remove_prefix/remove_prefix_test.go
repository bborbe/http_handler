package remove_prefix

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/bborbe/assert"
	server_mock "github.com/bborbe/http/mock"
)

func TestImplementsHandler(t *testing.T) {
	r := New("", nil)
	var i *http.Handler
	if err := AssertThat(r, Implements(i)); err != nil {
		t.Fatal(err)
	}
}

func TestSubhandlerCalledOneTime(t *testing.T) {
	prefix := ""
	counter := 0
	handler := New(prefix, func(http.ResponseWriter, *http.Request) {
		counter++
	})
	responseWriter := server_mock.NewHttpResponseWriterMock()
	request, err := server_mock.NewHttpRequestMock(fmt.Sprintf("http://www.example.com%s/bla.json", prefix))
	if err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(counter, Is(0)); err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(responseWriter, request)
	if err := AssertThat(counter, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestSubhandlerRemovePrefixFromRequestURI(t *testing.T) {
	prefix := "/test123456"
	var request *http.Request
	handler := New(prefix, func(resp http.ResponseWriter, req *http.Request) {
		request = req
	})
	responseWriter := server_mock.NewHttpResponseWriterMock()
	request, err := server_mock.NewHttpRequestMock(fmt.Sprintf("http://www.example.com%s/bla.json", prefix))
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(responseWriter, request)
	if err := AssertThat(request, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(request.RequestURI, Is("/bla.json")); err != nil {
		t.Fatal(err)
	}
}

func TestSubhandlerRemovePrefixFromURL(t *testing.T) {
	prefix := "/test123456"
	var request *http.Request
	handler := New(prefix, func(resp http.ResponseWriter, req *http.Request) {
		request = req
	})
	responseWriter := server_mock.NewHttpResponseWriterMock()
	request, err := server_mock.NewHttpRequestMock(fmt.Sprintf("http://www.example.com%s/bla.json", prefix))
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(responseWriter, request)
	if err := AssertThat(request, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(request.URL.String(), Is("http://www.example.com/bla.json")); err != nil {
		t.Fatal(err)
	}
}
