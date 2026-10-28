# Insecure TLS for Go

This is an **insecure** fork of golang's `crypto/tls` library.

It is motivated by ongoing efforts to deprecate and eventually remove weak or broken protocols and ciphers from the standard library (e.g. issues [32716](https://github.com/golang/go/issues/32716), [45428](https://github.com/golang/go/issues/45428), [63413](https://github.com/golang/go/issues/63413)), and by the need for [driftnet](https://driftnet.io) to maintain support for those protocols and ciphers.

This fork aims for minimal deviation from the official `crypto/tls`, whilst supporting protocols and cipher suites which have either already been removed from the standard library, will be removed in an upcoming release, or which were never included in the first place.

The current version is based on go1.27.1 and requires Go 1.27 or newer.

## When to use this library

This library is inherently insecure and should not be used in any situation where security is a requirement.

It might be suitable in the rare case where

- you must connect to a client or server which only supports a broken version of TLS, and
- you cannot upgrade that client or server to a non-broken version, and
- you are willing to completely give up the confidentiality, integrity and authentication that TLS provides.

In all other cases, stick to `crypto/tls`.

## How to use this library

Instead of

```
import "crypto/tls"
```

simply

```
import tls "github.com/driftnet-io/insecure-tls"
```

As a convenience, the fork re-uses `(crypto/tls).ConnectionState` as-is, and makes it accessible as `(github.com/driftnet-io/insecure-tls).ConnectionState.ConnectionState`. Note that its `CurveID` field therefore has the `(crypto/tls).CurveID` type.

## Changes with respect to the standard library

### SSLv3 support

Server-side and client-side SSLv3 support is re-introduced. Like TLS 1.0 and TLS 1.1 in `crypto/tls`, it is not enabled by default: set `Config.MinVersion` to `VersionSSL30` explicitly. (Versions up to v1.1.1 of this fork enabled SSLv3 by default.)

### Additional cipher suites

The following additional cipher suites are added to `InsecureCipherSuites()`:

- `TLS_RSA_WITH_RC4_128_MD5` (protocol versions SSLv3 to TLS 1.2),
- `TLS_RSA_WITH_AES_256_CBC_SHA256` (TLS 1.2 only),
- `TLS_ECDHE_ECDSA_WITH_3DES_EDE_CBC_SHA` (SSLv3 to TLS1.2),
- `TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384` (TLS 1.2 only), and
- `TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384` (TLS 1.2 only).

### GODEBUG settings

The standard library's `internal/godebug` is not available outside of the standard library, so the fork reads `GODEBUG` settings (`tlsmaxrsasize`, `tlssha1`, `tlsmlkem`, `tlssecpmlkem`, `fips140ems`) from the `GODEBUG` environment variable only, on every use. The `//go:debug` directives and the default `GODEBUG` of the main module's `go.mod` are not consulted. Defaults are the same as in `crypto/tls`; in particular the maximum RSA key size in certificates is 8192 bits unless `GODEBUG=tlsmaxrsasize=n` is set. (Versions up to v1.1.1 of this fork used a fixed limit of 16384 bits.)

### RSA key exchange and 3DES cipher suites

As in `crypto/tls` since Go 1.22 and Go 1.23, RSA key exchange and 3DES cipher suites are implemented but not used unless they are listed explicitly in `Config.CipherSuites`; the `tlsrsakex` and `tls3des` `GODEBUG` settings no longer exist upstream. (Versions up to v1.1.1 of this fork always acted as if `GODEBUG=tlsrsakex=1` were set.)

### Exporting key material

This code always acts as if `GODEBUG=tlsunsafeekm=1` were set: `ExportKeyingMaterial` works without Extended Master Secret.

### Internal packages

The standard library packages that `crypto/tls` imports but that cannot be imported from outside of it are replaced by small copies or stubs: `internal/tls13` (the TLS 1.3 key schedule on top of `crypto/hkdf`), `internal/fips140tls`, `stubs/godebug`, `stubs/boring` (BoringCrypto permanently disabled), `stubs/testenv` and `stubs/cryptotest` (test helpers). AES-GCM uses `crypto/cipher` instead of the FIPS 140-3 module, and the TLS 1.2 PRF is implemented locally.

## Testing

The handshake tests replay transcripts recorded against OpenSSL, as in the standard library. Transcripts for the cipher suites and the protocol version added by this fork are not recorded: the test harness only records transcripts with OpenSSL 1.1.1 built from source, and recording SSLv3 and RC4-MD5 handshakes additionally needs that build to be configured with `enable-weak-ssl-ciphers enable-ssl3 enable-ssl3-method`. They are covered by loopback handshakes between the fork's own client and server in `insecure_test.go` instead, and the SSLv3 key derivation is checked against a known-answer vector in `prf_test.go`. The BoGo shim, the linker, the FIPS 140-3 and the certificate generator tests of the standard library are not included.
