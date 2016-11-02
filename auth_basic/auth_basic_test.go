package auth_basic

import (
	"net/http"
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsHandler(t *testing.T) {
	r := New(nil, nil, "realm")
	var i *http.Handler
	if err := AssertThat(r, Implements(i)); err != nil {
		t.Fatal(err)
	}
}
