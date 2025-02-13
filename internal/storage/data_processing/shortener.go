package fnc

import (
	"crypto/sha1"
)

var ALPHABET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func URLShortener(url string) string {
	hash := sha1.New()
	hash.Write([]byte(url))
	s := hash.Sum(nil)
	result := make([]byte, 10)
	k := 9
	for k >= 0 {
		for i := range s {
			if k < 0 {
				break
			}
			num := s[i]
			for num > 0 && k >= 0 {
				digit := num % 63
				result[k] = ALPHABET[digit]
				k--
				num = num / 63
			}
		}
	}
	return string(result)
}
