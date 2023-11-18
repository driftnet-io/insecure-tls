# Insecure TLS for Go

This is an **insecure** fork of golang's `crypto/tls` library.

It is motivated by ongoing efforts to deprecate and eventually remove weak or broken protocols and ciphers from the standard library (e.g. issues [32716](https://github.com/golang/go/issues/32716), [45428](https://github.com/golang/go/issues/45428), [63413](https://github.com/golang/go/issues/63413)), and by the need for [driftnet](https://driftnet.io) to maintain support for those protocols and ciphers.

This fork aims for minimal deviation from the official `crypto/tls`, whilst supporting protocols and cipher suites which have either already been removed from the standard library, will be removed in an upcoming release, or which were never included in the first place.

The current version is based on go1.21.4.


## When to use this library

This library is inherently insecure and should not be used in any situation where security is a requirement.

It might be suitable in the rare case where

  * you must connect to a client or server which only supports a broken version of TLS, and
  * you cannot upgrade that client or server to a non-broken version, and
  * you are willing to completely give up the confidentiality, integrity and authentication that TLS provides.

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

As a convenience, the fork re-uses `(crypto/tls).ConnectionState` as-is, and makes it accessible as `(github.com/driftnet-io/insecure-tls).ConnectionState.ConnectionState`.


## Changes with respect to the standard library

### SSLv3 server support

Server-side SSLv3 support is re-introduced, and enabled by default.

### SSLv3 client support

SSLv3 client support is introduced, and enabled by default.

### Additional cipher suites

The cipher suite `TLS_RSA_WITH_RC4_128_MD5` is added to `InsecureCipherSuites()`. It may be used with protocol versions SSLv3 to TLS 1.2.

### Maximum size of server RSA certificates

This is usually configurable by setting `GODEBUG=tlsmaxrsasize=n`. However, as `GODEBUG` internals are not available to us outside the standard library, this is set to a fixed value of 16384 bytes.
