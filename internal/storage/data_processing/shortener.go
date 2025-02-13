package fnc

import (
	"crypto/sha1"
)

var ALPHABET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_" //63

func URLShortener(url string) string {
	hash := sha1.New()
	hash.Write([]byte(url))
	s := hash.Sum(nil)[:10]
	res := ""
	for i := range s {
		res += EncodeTo63(s[i])
	}
	//shortURL := base64.URLEncoding.EncodeToString(hash.Sum(nil))[:8]
	//return fmt.Sprintf("%x", s)
	return res[:10]
}

func EncodeTo63(num byte) string {
	res := []string{}
	for num > 0 {
		digit := num % 63
		res = append(res, string(ALPHABET[digit]))
		num = num / 63
	}
	str := ""
	for i := len(res) - 1; i >= 0; i-- {
		str += res[i]
	}
	return str
}
