package redirect_relative

import (
	"net/http"
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsHandler(t *testing.T) {
	object := New("/target")
	var expected *http.Handler
	if err := AssertThat(object, Implements(expected)); err != nil {
		t.Fatal(err)
	}
}
