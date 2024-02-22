package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

// isModernDigit является ли руна современной арабской цифрой от 0 до 9. В рамках задачи цифры только от 0 до 9,
// unicode.IsDigit не подошел из-за некоторых крайних случаев в интересных языках, подробнее в последних тесткейсах.
func isModernDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

// checkString проверить является ли строка корректной для распаковки.
// Корректная строка начиается не с цифры и может содержать только цифры, но не числа.
func checkString(s string) error {
	if len(s) == 0 {
		return nil
	}
	runes := []rune(s)
	//  не может начинаться с цифры
	if isModernDigit(runes[0]) {
		return ErrInvalidString
	}
	// не может содержать чисел
	for i := 1; i < len(runes); i++ {
		if isModernDigit(runes[i-1]) && isModernDigit(runes[i]) {
			return ErrInvalidString
		}
	}
	return nil
}

var ErrInvalidString = errors.New("invalid string")

// Unpack распоковать сжатую стркоу в полную строку.
func Unpack(s string) (string, error) {
	if len(s) == 0 {
		return s, nil
	}
	err := checkString(s)
	if err != nil {
		return s, err
	}
	// корректная строка из 1 символа, распоковывать нечего
	if len(s) == 1 {
		return s, nil
	}
	var builder strings.Builder
	runes := []rune(s)
	for i := 0; i < len(runes); {
		next := i + 1
		if next < len(runes) && isModernDigit(runes[next]) {
			n, _ := strconv.Atoi(string(runes[next]))
			if n > 0 {
				builder.WriteString(strings.Repeat(string(runes[i]), n))
			}
			i += 2
		} else {
			builder.WriteRune(runes[i])
			i++
		}
	}
	return builder.String(), nil
}
