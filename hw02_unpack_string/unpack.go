package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	// Place your code here
	if str == "" {
		return "", nil
	}

	var (
		result     strings.Builder
		runeStr    = []rune(str)
		escape     bool
		escapeChar = '\\'
	)

	for i := 0; i < len(runeStr); i++ {
		cur := runeStr[i]

		if i == 0 {
			if unicode.IsDigit(cur) {
				return "", ErrInvalidString
			}
			if cur == escapeChar {
				escape = true
				continue
			}
			result.WriteRune(cur)
			continue
		}

		pre := runeStr[i-1]

		if unicode.IsDigit(cur) && !escape {
			s, err := unpackHelper(cur, pre)
			if err != nil {
				return "", err
			}
			result.WriteString(s)
			continue
		}

		if unicode.IsDigit(cur) && escape {
			result.WriteRune(cur)
			escape = false
			continue
		}

		if cur == escapeChar && !escape {
			escape = true
			continue
		}

		if cur == escapeChar && escape {
			result.WriteRune(cur)
			escape = false
			continue
		}

		result.WriteRune(cur)
	}

	return result.String(), nil
}

func unpackHelper(cur, pre rune) (string, error) {
	const one = 1
	n, err := strconv.Atoi(string(cur))
	if err != nil {
		return "", err
	}
	if n < one {
		return "", ErrInvalidString
	}
	return strings.Repeat(string(pre), n-one), nil
}
