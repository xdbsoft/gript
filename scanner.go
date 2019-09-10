package gript

import (
	"bufio"
	"bytes"
	"io"
)

// scanner represents a lexical scanner.
type scanner struct {
	r    *bufio.Reader
	last token
}

// newScanner returns a new instance of Scanner.
func newScanner(r io.Reader) *scanner {
	return &scanner{r: bufio.NewReader(r)}
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *scanner) unread() { _ = s.r.UnreadRune() }

// Scan returns the next token and literal value.
func (s *scanner) Scan() (tok token, lit string) {

	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	if isWhitespace(ch) {
		s.unread()
		tok, lit = s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()
		tok, lit = s.scanIdent()
	} else if isDigit(ch) || (isMinusOrPlus(ch) && (s.last == tokIllegal || s.last == tokLeftParenthesis || s.last == tokOperator)) {
		s.unread()
		tok, lit = s.scanNumber()
	} else if isOperator(ch) {
		s.unread()
		tok, lit = s.scanOperator()
	} else if isQuote(ch) {
		s.unread()
		tok, lit = s.scanString()
	} else {

		// Otherwise read the individual character.
		switch ch {
		case eof:
			tok, lit = tokEOF, ""
		case '(':
			tok, lit = tokLeftParenthesis, "("
		case ')':
			tok, lit = tokRightParenthesis, ")"
		default:
			tok, lit = tokIllegal, string(ch)
		}
	}

	if tok != tokWhitespace {
		s.last = tok
	}
	return tok, lit
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *scanner) scanWhitespace() (tok token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return tokWhitespace, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *scanner) scanIdent() (tok token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// Otherwise return as a regular identifier.
	v := buf.String()
	switch v {
	case "in":
		return tokOperator, v
	}
	return tokIdentifier, v
}

// scanNumber consumes the current rune and all contiguous runes creating an int or a float.
func (s *scanner) scanNumber() (tok token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	tok = tokInt

	// Read every subsequent ident character into the buffer.
	// Non-int characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isDigit(ch) && !isDot(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
			if isDot(ch) {
				if tok == tokInt {
					tok = tokFloat
				} else {
					return tokIllegal, buf.String()
				}
			}
		}
	}

	// Otherwise return as a regular identifier.
	return tok, buf.String()
}

// scanOperator consumes the current rune and all contiguous operator.
func (s *scanner) scanOperator() (tok token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent operator character into the buffer.
	// Non-operator characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isOperator(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return tokOperator, buf.String()
}

// scanString consumes the current rune (quote) and all rune until next quote
func (s *scanner) scanString() (tok token, lit string) {

	var buf bytes.Buffer
	q := s.read() //Skip first quote

	// Read every subsequent character into the buffer.
	// Quote and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			return tokIllegal, buf.String()
		} else if ch == q {
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return tokString, buf.String()

}
