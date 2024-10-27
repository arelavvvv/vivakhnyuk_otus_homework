package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var result strings.Builder
	n := len(s)
	i := 0

	for i < n {
		char, size := utf8.DecodeRuneInString(s[i:])
		if size == 0 {
			break
		}

		if unicode.IsDigit(char) {
			return "", ErrInvalidString
		}

		i += size
		count := 1

		if i < n && unicode.IsDigit(rune(s[i])) {
			var err error
			count, err = getCount(s, &i)
			if err != nil {
				return "", err
			}
		}

		result.WriteString(strings.Repeat(string(char), count))
	}

	return result.String(), nil
}

func getCount(s string, index *int) (int, error) {
	countStr := ""
	n := len(s)

	for *index < n && unicode.IsDigit(rune(s[*index])) {
		countStr += string(s[*index])
		(*index)++
	}

	count, err := strconv.Atoi(countStr)
	if err != nil || count < 0 || count > 9 {
		return 0, ErrInvalidString
	}

	return count, nil
}
