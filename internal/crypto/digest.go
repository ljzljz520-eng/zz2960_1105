package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

func Digest(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }

func Seal(label, payload string) string { return Digest([]byte(label + "|" + payload)) }

func Verify(data []byte, expected string) bool { return Digest(data) == expected }
