package util

import (
	"unicode"
	//"unicode/utf8"
)

func ToChars(s string) []rune {
	var chars = make([]rune, 1, 10)
	for _, c := range s {
		chars = append(chars, c)
	}
	return chars
}

func ReadUntilToChar(start int, cs []rune, char rune) string {
	var chars = make([]rune, 0, 10)
	for _, c := range cs[start:] {
		if c == char {
			break
		}
		chars = append(chars, c)
	}
	return string(chars)
}

func ReadUntilToString(start int, cs []rune, s string) string {
	var chars = make([]rune, 0, 10)
	var innerChars = make([]rune, 0, 10)
	for _, c := range cs[start:] {
		if string(chars) == s {
			break
		}
		innerChars = append(innerChars, c)
		chars = append(chars, c)
		if len(chars) > len(s) {
			chars = make([]rune, 0, 10)
		}
	}
	return string(innerChars)
}

func ReadUntil(text string, char rune) string {
	count := 0
	var chars = make([]rune, 1, 10)
	for _, c := range text {
		if c == char {
			break
		}
		chars = append(chars, c)
		count++
	}
	return string(chars)
}

func ReadSkipBlank(text string, blankCount int) string {
	count := 0
	lastchar := ' '
	var chars = make([]rune, 1, 10)
	for _, c := range text {
		if lastchar != ' ' && c == ' ' {
			count++
		}
		chars = append(chars, c)
		if count > blankCount {
			break
		}
		lastchar = c
	}
	return string(chars)
}

func IsNewLine(char rune) bool {
	return char == '\r' || char == '\n' || char == '\u0085' || char == '\u2028' || char == '\u2029'
}

func IsNewLineString(str string) bool {
	cs := ToChars(str)
	return (len(cs) == 1 && (IsNewLine(cs[0]))) || (str == "\r\n")
}

func IsWhiteSpace(char rune) bool {
	return char == ' ' || char == '\f' || char == '\t' || char == '\u000B'
}

func IsNewLineOrWhiteSpace(char rune) bool {
	return IsNewLine(char) || IsWhiteSpace(char)
}
func IsLetter(char rune) bool {
	return char == unicode.LowerCase || char == unicode.UpperCase || char == unicode.TitleCase
}

func IsNumber(char rune) bool {
	return unicode.IsNumber(char)
}

func IsNumberOrLetter(char rune) bool {
	return IsLetter(char) || IsNumber(char)
}

func IsHexNumber(char rune) bool {
	return (char < '9' && char > '0') || (char >= 'A' && char <= 'F') || (char >= 'a' && char <= 'f')
}

func IsIdentifierStart(char rune) bool {
	return char == '_' || IsLetter(char)
}
