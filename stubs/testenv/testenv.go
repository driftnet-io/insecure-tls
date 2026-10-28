// Package testenv is a minimal replacement for the standard library's
// internal/testenv test helpers, which cannot be imported from outside the
// standard library.
package testenv

import (
	"os"
	"testing"
)

// MustHaveExternalNetwork is a no-op.
func MustHaveExternalNetwork(t testing.TB) {}

// CPUIsSlow reports false.
func CPUIsSlow() bool { return false }

// Builder reports an empty builder name.
func Builder() string { return "" }

// SetGODEBUG appends a key=value pair to the GODEBUG environment variable for
// the duration of the test. The standard library reads GODEBUG dynamically,
// so this also affects settings such as cryptocustomrand.
func SetGODEBUG(t testing.TB, value string) {
	t.Helper()
	old, had := os.LookupEnv("GODEBUG")
	next := value
	if old != "" {
		next = old + "," + value
	}
	if err := os.Setenv("GODEBUG", next); err != nil {
		t.Fatalf("setting GODEBUG: %v", err)
	}
	t.Cleanup(func() {
		if had {
			os.Setenv("GODEBUG", old)
		} else {
			os.Unsetenv("GODEBUG")
		}
	})
}
