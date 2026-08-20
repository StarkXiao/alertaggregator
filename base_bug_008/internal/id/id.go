package id

import (
	"crypto/rand"
	"encoding/hex"
)

func New(p string) string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return p + "_" + hex.EncodeToString(b)
}
