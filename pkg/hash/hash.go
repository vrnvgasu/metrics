// Package hash предоставляет утилиты для вычисления подписей.
package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// PrepareHeaderHashSHA256 вычисляет HMAC-SHA256 от data с ключом secret и возвращает hex-строку.
func PrepareHeaderHashSHA256(secret string, data []byte) (string, error) {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(data)

	return hex.EncodeToString(h.Sum(nil)), nil
}
