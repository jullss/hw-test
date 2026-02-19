package hw02unpackstring

import (
	"errors"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	if unicode.IsDigit(rune(s[0])) {
		return "", ErrInvalidString
	}

	res := make([]rune, 0, len(s))
	var prevRune rune
	var isEscaped bool
	var isPrevDigit bool

	for _, r := range s {
		if isEscaped {
			if !unicode.IsDigit(r) && r != '\\' {
				return "", ErrInvalidString
			}
			res = append(res, r)
			prevRune = r

			isEscaped = false
			isPrevDigit = false

			continue
		}

		if r == '\\' {
			isEscaped = true
			isPrevDigit = false

			continue
		}

		if unicode.IsDigit(r) {
			if isPrevDigit {
				return "", ErrInvalidString
			}

			count := int(r - '0')

			if count == 0 {
				res = res[:len(res)-1]
			} else if count > 1 {
				res = append(res, multiplyRunes(count, prevRune)...)
			}

			isPrevDigit = true
		} else {
			res = append(res, r)
			prevRune = r

			isPrevDigit = false
		}
	}

	if isEscaped {
		return "", ErrInvalidString
	}

	return string(res), nil
}

func multiplyRunes(count int, r rune) []rune {
	if count < 1 {
		return nil
	}

	res := make([]rune, count-1)

	for i := range res {
		res[i] = r
	}

	return res
}
