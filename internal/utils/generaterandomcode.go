package utils

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

func GenerateRandomCode(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	code := base64.URLEncoding.EncodeToString(bytes)
	code = strings.ReplaceAll(code, "-", "")
	code = strings.ReplaceAll(code, "_", "")
	if len(code) > length {
		code = code[:length]
	}
	return code
}
