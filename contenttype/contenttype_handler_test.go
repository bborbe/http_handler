package contenttype

import (
	"net/http"
	"testing"

	. "github.com/bborbe/assert"
	server_mock "github.com/bborbe/http/mock"
	"github.com/bborbe/http_handler/static"
)

func TestImplementsHandler(t *testing.T) {
	r := NewContentTypeHandler(nil)
	var i *http.Handler
	err := AssertThat(r, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}

func TestIgnoreUnkownExtention(t *testing.T) {
	subHandler := static.NewHandlerStaticContent("foo bar")
	handler := NewContentTypeHandler(subHandler)
	responseWriter := server_mock.NewHttpResponseWriterMock()
	request, err := server_mock.NewHttpRequestMock("http://www.example.com/bla")
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(responseWriter, request)
	err = AssertThat(responseWriter.Header().Get("Content-Type"), Is(""))
	if err != nil {
		t.Fatal(err)
	}
}

func TestKownExtention(t *testing.T) {
	subHandler := static.NewHandlerStaticContent("foo bar")
	handler := NewContentTypeHandler(subHandler)
	responseWriter := server_mock.NewHttpResponseWriterMock()
	request, err := server_mock.NewHttpRequestMock("http://www.example.com/bla.json")
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(responseWriter, request)
	err = AssertThat(responseWriter.Header().Get("Content-Type"), Is("application/json"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetExtension(t *testing.T) {
	var err error
	err = AssertThat(getExtension(""), Is(""))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(getExtension("bla"), Is(""))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(getExtension("bla."), Is(""))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(getExtension("bla.jpg"), Is("jpg"))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(getExtension("....bla.gif"), Is("gif"))
	if err != nil {
		t.Fatal(err)
	}
}
