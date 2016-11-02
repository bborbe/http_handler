package auth_html

import (
	"net/http"
	"testing"

	"time"

	. "github.com/bborbe/assert"
)

func TestImplementsHandler(t *testing.T) {
	r := New(nil, nil, nil)
	var i *http.Handler
	if err := AssertThat(r, Implements(i)); err != nil {
		t.Fatal(err)
	}
}

func TestCreateExpireNotConstant(t *testing.T) {
	t1 := createExpires()
	time.Sleep(100 * time.Millisecond)
	t2 := createExpires()
	if err := AssertThat(t1, Not(Is(t2))); err != nil {
		t.Fatal(err)
	}
}
