package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	if len(s) == 0 {
		return s, nil
	}
	var builder strings.Builder
	var prev rune
	for i, c := range s {
		if i == 0 && unicode.IsDigit(c) {
			return s, ErrInvalidString
		}
		if unicode.IsLetter(c) {
			if unicode.IsLetter(prev) {
				_, err := builder.WriteRune(prev)
				if err != nil {
					return s, err
				}
			}
			prev = c
		} else if unicode.IsDigit(c) {
			if unicode.IsDigit(prev) {
				return s, ErrInvalidString
			}
			n, _ := strconv.Atoi(string(c))
			_, err := builder.WriteString(strings.Repeat(string(prev), n))
			if err != nil {
				return s, err
			}
			prev = c
		}
		if i == (len(s)-1) && unicode.IsLetter(c) {
			_, err := builder.WriteRune(c)
			if err != nil {
				return s, err
			}
		}
	}
	return builder.String(), nil
}
