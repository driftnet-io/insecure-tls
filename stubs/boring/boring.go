package boring

import "crypto/cipher"

// Enabled marks the boring package as disabled
var Enabled = false

// Unreachable is a no-op
func Unreachable() {
	return
}

// NewGCMTLS should never be called
func NewGCMTLS(cipher cipher.Block) (cipher.AEAD, error) {
	panic("boring.NewGCMTLS() called")
}
