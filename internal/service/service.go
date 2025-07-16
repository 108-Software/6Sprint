package service

import (
	"errors"
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func AutoDetectAndConvert(input string) (string, error) {
	if strings.TrimSpace(input) == "" {
		return "", errors.New("empty input")
	}

	if isMorseCode(input) {
		return morse.ToText(input), nil
	}
	return morse.ToMorse(input), nil
}

func isMorseCode(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}

	allowedChars := map[rune]bool{
		'.': true,
		'-': true,
		' ': true,
		'/': true,
	}

	for _, r := range s {
		if !allowedChars[r] {
			return false
		}
	}
	return true
}
