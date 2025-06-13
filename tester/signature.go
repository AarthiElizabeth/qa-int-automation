package tester

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func GenerateSignature(secret string, body []byte) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(body)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
