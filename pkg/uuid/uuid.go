package uuid

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

// NewString generates a version 4 UUID string.
func NewString() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	// Set version (4) and variant bits per RFC 4122.
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	hexStr := hex.EncodeToString(b[:])
	var builder strings.Builder
	builder.Grow(36)
	builder.WriteString(hexStr[0:8])
	builder.WriteByte('-')
	builder.WriteString(hexStr[8:12])
	builder.WriteByte('-')
	builder.WriteString(hexStr[12:16])
	builder.WriteByte('-')
	builder.WriteString(hexStr[16:20])
	builder.WriteByte('-')
	builder.WriteString(hexStr[20:])
	return builder.String(), nil
}
