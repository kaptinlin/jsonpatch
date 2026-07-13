package op

import "unicode/utf8"

func codePointCount(value string) int {
	return utf8.RuneCountInString(value)
}

func codePointSubstring(value string, pos, length int) (string, error) {
	points := []rune(value)
	if pos < 0 || pos > len(points) {
		return "", ErrPositionOutOfStringRange
	}
	end := pos + length
	if end > len(points) {
		return "", ErrSubstringTooLong
	}
	return string(points[pos:end]), nil
}
