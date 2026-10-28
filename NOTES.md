# Update notes

The repository is a flat copy of `src/crypto/tls` from a Go release with the
changes described in README.md applied on top; there is no upstream git
history to rebase onto.

To move to a new Go release:

1. Install that release and run, in this repository:

       rm -f *.go && rm -rf testdata internal fipsonly
       cp "$(go env GOROOT)"/src/crypto/tls/*.go .
       cp -R "$(go env GOROOT)"/src/crypto/tls/{testdata,internal,fipsonly} .

2. Replace the imports of packages that cannot be used outside of the
   standard library (`crypto/internal/boring`, `crypto/internal/fips140/...`,
   `crypto/tls/internal/fips140tls`, `internal/godebug`, `internal/cpu`,
   `internal/byteorder`, `internal/testenv`, `crypto/internal/cryptotest`)
   with the packages under `internal/` and `stubs/`, or with public
   equivalents, and refresh those copies from the new release if needed.

3. Re-apply the changes listed in README.md. `git diff` against the previous
   tag shows exactly what was changed last time.

4. Drop the tests that need the standard library's test infrastructure
   (`bogo_shim_*_test.go`, `link_test.go`, `fips140_test.go`,
   `certificates_generator_test.go`) and keep `fips140_stub_test.go` and
   `insecure_test.go`.

5. Update `go.mod`, run `gofmt -l .`, `go vet ./...` and `go test ./...`.
