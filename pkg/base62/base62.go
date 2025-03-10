package base62

import "strings"

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Encode(n int64) string {
	if n == 0 {
		return "0"
	}

	var encoded strings.Builder
	for n > 0 {
		remainder := n % 62
		encoded.WriteByte(base62Chars[remainder])
		n = n / 62
	}

	encodedStr := reverse(encoded.String())
	for len(encodedStr) < 6 {
		encodedStr = "0" + encodedStr
	}
	return encodedStr
}
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
