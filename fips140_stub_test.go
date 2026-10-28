// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Helpers kept from the standard library's fips140_test.go, which is not
// included in this fork because it needs the crypto/internal/cryptotest
// test infrastructure.

package tls

func isFIPSVersion(v uint16) bool {
	return v == VersionTLS12 || v == VersionTLS13
}
