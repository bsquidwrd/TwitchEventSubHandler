package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
)

func ValidateSignature(secret []byte, messageID string, messageTimestamp string, body []byte, messageSignature string) bool {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(messageID))
	mac.Write([]byte(messageTimestamp))
	mac.Write(body)
	expectedMAC := mac.Sum(nil)

	signatureBytes, err := hex.DecodeString(messageSignature)
	if err != nil {
		slog.Error("Error decoding signature:", "error", err)
		return false
	}

	return hmac.Equal(signatureBytes, expectedMAC)
}
