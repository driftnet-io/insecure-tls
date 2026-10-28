// Package godebug is a minimal replacement for the standard library's
// internal/godebug package, which cannot be imported from outside the
// standard library.
//
// Settings are read from the GODEBUG environment variable on every call to
// Value, so they can be changed at run time (for example with t.Setenv in
// tests). Unlike the standard library, the //go:debug directives of the main
// package and the go.mod default GODEBUG are not consulted.
package godebug

import (
	"os"
	"strings"
)

// defaults are the values used when a setting is not present in GODEBUG.
// An empty string means the standard library default.
var defaults = map[string]string{}

// A Setting is a single setting in the $GODEBUG environment variable.
type Setting struct {
	name string
}

// New returns a new Setting for the $GODEBUG setting with the given name.
func New(name string) *Setting {
	return &Setting{name: strings.TrimPrefix(name, "#")}
}

// Name returns the name of the setting.
func (s *Setting) Name() string {
	return s.name
}

// String returns a printable form for the setting: name=value.
func (s *Setting) String() string {
	return s.name + "=" + s.Value()
}

// Value returns the current value for the GODEBUG setting s.
func (s *Setting) Value() string {
	if v, ok := lookup(os.Getenv("GODEBUG"), s.name); ok {
		return v
	}
	return defaults[s.name]
}

// IncNonDefault is a no-op: the standard library uses it to count non-default
// behaviors in runtime/metrics, which are not available here.
func (s *Setting) IncNonDefault() {}

// lookup returns the last value of name in a GODEBUG string, since later
// settings override earlier ones.
func lookup(godebug, name string) (string, bool) {
	var value string
	var found bool
	for _, kv := range strings.Split(godebug, ",") {
		k, v, ok := strings.Cut(kv, "=")
		if ok && k == name {
			value, found = v, true
		}
	}
	return value, found
}
