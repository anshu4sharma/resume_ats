package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func VerifyGitHubSignature(secret string, body []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(signature))
}
