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

		
		if i < n {
			nextChar, _ := utf8.DecodeRuneInString(s[i:])
			if unicode.IsDigit(nextChar) {
				countStr := string(nextChar)
				i += utf8.RuneLen(nextChar)

				for i < n {
					nextChar, _ := utf8.DecodeRuneInString(s[i:])
					if unicode.IsDigit(nextChar) {
						countStr += string(nextChar)
						i += utf8.RuneLen(nextChar)
					} else {
						break
					}
				}

				count, err := strconv.Atoi(countStr)
				if err != nil || count < 0 || count > 9 {
					return "", ErrInvalidString
				}

				if count > 0 {
					result.WriteString(strings.Repeat(string(char), count))
				}
			} else {
				result.WriteRune(char)
			}
		} else {
			result.WriteRune(char)
		}
	}

	return result.String(), nil
}
