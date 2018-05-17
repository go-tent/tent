package hash

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
)

// Verify calculates and compares the sha256 hash.
func Verify(r io.Reader, hash []byte) error {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return err
	}
	if exp := h.Sum(nil); !bytes.Equal(exp, hash) {
		return fmt.Errorf("conflict (expected %x, got %x)", exp, hash)
	}
	return nil
}
