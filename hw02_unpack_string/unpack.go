package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var result strings.Builder
	n := len(s)
	i := 0

	for i < n {
		char := s[i]

		if unicode.IsDigit(rune(char)) {
			return "", ErrInvalidString
		}

		i++

		if i < n && unicode.IsDigit(rune(s[i])) {
			countStr := string(s[i])
			i++
			for i < n && unicode.IsDigit(rune(s[i])) {
				countStr += string(s[i])
				i++
			}
			count, err := strconv.Atoi(countStr)
			if err != nil || count < 0 || count > 9 {
				return "", ErrInvalidString
			}
			if count > 0 {
				result.WriteString(strings.Repeat(string(char), count))
			}
		} else {
			result.WriteByte(char)
		}
	}

	return result.String(), nil
}
