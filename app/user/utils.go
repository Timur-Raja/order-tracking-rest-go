package user

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// generates an URL safe token of "Len" bytes
func GenerateSessionToken(Len int) (string, error) {
	b := make([]byte, Len)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("could not generate session token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
