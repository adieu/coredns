package hostname

import (
	"testing"

	"github.com/mholt/caddy"
)

func TestSetupWhoami(t *testing.T) {
	c := caddy.NewTestController("dns", `hostname`)
	if err := setupHostname(c); err != nil {
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	c = caddy.NewTestController("dns", `hostname example.org`)
	if err := setupHostname(c); err == nil {
		t.Fatalf("Expected errors, but got: %v", err)
	}
}
