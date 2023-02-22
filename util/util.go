package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"time"
)

func GetSign(secret string) (string, string) {
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
	stringToSign := fmt.Sprintf("%s\n%s", timestamp, secret)
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(stringToSign))
	b64 := base64.StdEncoding.EncodeToString(h.Sum(nil))
	sign := url.QueryEscape(b64)
	return sign, timestamp
}
