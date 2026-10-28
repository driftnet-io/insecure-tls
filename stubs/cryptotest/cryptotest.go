// Package cryptotest is a minimal replacement for the standard library's
// crypto/internal/cryptotest test helpers, which cannot be imported from
// outside the standard library.
package cryptotest

import "testing"

// SetGlobalRandom is a no-op: the global crypto/rand source cannot be replaced
// from outside the standard library. Tests that need deterministic randomness
// use Config.Rand together with GODEBUG=cryptocustomrand=1 instead.
func SetGlobalRandom(t testing.TB, seed int64) {}

// MustMinimumFIPS140ModuleVersion is a no-op: this fork never runs with the
// FIPS 140-3 module.
func MustMinimumFIPS140ModuleVersion(t testing.TB, version string) {}
