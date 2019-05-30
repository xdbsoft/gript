package gript

type token int

const (
	tokIllegal token = iota
	tokEOF
	tokWhitespace

	tokLeftParenthesis
	tokRightParenthesis

	tokIdentifier

	tokOperator

	tokInt
	tokFloat
	tokString
)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isOperator(ch rune) bool {
	return ch == '>' || ch == '<' || ch == '=' || ch == '!' || ch == '+' || ch == '-' || ch == '/' || ch == '*' || ch == '|' || ch == '&' || ch == '%'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '.'
}

func isMinusOrPlus(ch rune) bool {
	return ch == '-' || ch == '+'
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

func isDot(ch rune) bool {
	return ch == '.'
}

func isQuote(ch rune) bool {
	return ch == '\'' || ch == '"' || ch == '`'
}

var eof = rune(0)
