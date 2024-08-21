package hmac

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
)

func CalculateHMAC(message, secretKey string) string {
	//message = normalizeMessage(message)
	key := []byte(secretKey)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func ValidateHMAC(message, secretKey, digest string) bool {
	//message = normalizeMessage(message)
	expectedMAC := CalculateHMAC(message, secretKey)
	return hmac.Equal([]byte(digest), []byte(expectedMAC))
}

//func normalizeMessage(message string) string {
//	whitespaceRegex := regexp.MustCompile(`\s+`)
//	newlineRegex := regexp.MustCompile(`[\n\r]`)
//
//	message = strings.TrimSpace(message)
//	message = whitespaceRegex.ReplaceAllString(message, " ")
//	message = newlineRegex.ReplaceAllString(message, "")
//	return message
//}
