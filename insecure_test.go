// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tls

import (
	"testing"
)

// The cipher suites and the protocol version added by this fork are exercised
// with a loopback handshake between the fork's own client and server: OpenSSL
// builds with RC4, MD5 and SSLv3 support are not generally available, so no
// transcripts are recorded for them.
func TestInsecureCipherSuitesLoopback(t *testing.T) {
	tests := []struct {
		suite    uint16
		versions []uint16
	}{
		{TLS_RSA_WITH_RC4_128_MD5, []uint16{VersionSSL30, VersionTLS10, VersionTLS11, VersionTLS12}},
		{TLS_RSA_WITH_RC4_128_SHA, []uint16{VersionSSL30, VersionTLS10}},
		{TLS_RSA_WITH_3DES_EDE_CBC_SHA, []uint16{VersionSSL30, VersionTLS10}},
		{TLS_RSA_WITH_AES_128_CBC_SHA, []uint16{VersionSSL30, VersionTLS10}},
		{TLS_RSA_WITH_AES_256_CBC_SHA256, []uint16{VersionTLS12}},
		{TLS_ECDHE_ECDSA_WITH_3DES_EDE_CBC_SHA, []uint16{VersionTLS10, VersionTLS11, VersionTLS12}},
		{TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384, []uint16{VersionTLS12}},
		{TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384, []uint16{VersionTLS12}},
	}
	for _, tc := range tests {
		for _, vers := range tc.versions {
			t.Run(CipherSuiteName(tc.suite)+"-"+VersionName(vers), func(t *testing.T) {
				clientConfig := testConfigClient.Clone()
				clientConfig.MinVersion = vers
				clientConfig.MaxVersion = vers
				clientConfig.CipherSuites = []uint16{tc.suite}
				serverConfig := testConfigServer.Clone()
				serverConfig.MinVersion = vers
				serverConfig.MaxVersion = vers
				serverConfig.CipherSuites = []uint16{tc.suite}

				serverState, clientState, err := testHandshake(t, clientConfig, serverConfig)
				if err != nil {
					t.Fatalf("handshake failed: %v", err)
				}
				for _, state := range []ConnectionState{serverState, clientState} {
					if state.Version != vers {
						t.Errorf("negotiated %s, want %s", VersionName(state.Version), VersionName(vers))
					}
					if state.CipherSuite != tc.suite {
						t.Errorf("negotiated %s, want %s", CipherSuiteName(state.CipherSuite), CipherSuiteName(tc.suite))
					}
				}
			})
		}
	}
}

// SSLv3 client authentication uses the SSLv3 CertificateVerify construction,
// which mixes the master secret into the handshake hash.
func TestSSLv3ClientAuth(t *testing.T) {
	clientConfig := testConfigClient.Clone()
	clientConfig.MinVersion = VersionSSL30
	clientConfig.MaxVersion = VersionSSL30
	clientConfig.Certificates = []Certificate{testClientRSA2048Cert}
	serverConfig := testConfigServer.Clone()
	serverConfig.MinVersion = VersionSSL30
	serverConfig.MaxVersion = VersionSSL30
	serverConfig.ClientAuth = RequireAndVerifyClientCert

	serverState, clientState, err := testHandshake(t, clientConfig, serverConfig)
	if err != nil {
		t.Fatalf("handshake failed: %v", err)
	}
	if serverState.Version != VersionSSL30 || clientState.Version != VersionSSL30 {
		t.Errorf("negotiated %s / %s, want SSLv3", VersionName(serverState.Version), VersionName(clientState.Version))
	}
	if len(serverState.PeerCertificates) != 1 {
		t.Errorf("server saw %d client certificates, want 1", len(serverState.PeerCertificates))
	}
}

// SSLv3 is supported but, like TLS 1.0 and 1.1, not enabled by default.
func TestSSLv3DisabledByDefault(t *testing.T) {
	clientConfig := testConfigClient.Clone()
	clientConfig.MinVersion = VersionSSL30
	clientConfig.MaxVersion = VersionSSL30
	serverConfig := testConfigServer.Clone()
	if _, _, err := testHandshake(t, clientConfig, serverConfig); err == nil {
		t.Fatal("SSLv3 handshake succeeded against a server with the default minimum version")
	}
}
