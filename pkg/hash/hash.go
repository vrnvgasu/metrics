package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func PrepareHeaderHashSHA256(secret string, data []byte) (string, error) {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(data)

	return hex.EncodeToString(h.Sum(nil)), nil
}
