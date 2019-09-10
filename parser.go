package gript

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

// parser represents a parser.
type parser struct {
	s   *scanner
	buf struct {
		tok token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func newParser(r io.Reader) *parser {
	return &parser{s: newScanner(r)}
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *parser) scanIgnoreWhitespace() (tok token, lit string) {
	tok, lit = p.s.Scan()
	if tok == tokWhitespace {
		tok, lit = p.s.Scan()
	}
	return
}

func isRightAssociative(o string) bool {
	return o == "^"
}
func precedence(o string) int {
	switch o {
	case "*", "/", "%":
		return 6
	case "+", "-":
		return 5
	case "<", "<=", ">", ">=", "in":
		return 4
	case "==", "!=":
		return 3
	case "&&":
		return 2
	case "||":
		return 1
	}
	return 0
}

type stack []Expression

func (s *stack) Push(v Expression) {
	*s = append(*s, v)
}
func (s *stack) Pop() Expression {
	l := len(*s)
	r := (*s)[l-1]
	*s = (*s)[:l-1]
	return r
}

type opStack []string

func (s *opStack) Push(v string) {
	*s = append(*s, v)
}
func (s *opStack) Pop() string {
	l := len(*s)
	r := (*s)[l-1]
	*s = (*s)[:l-1]
	return r
}
func (s *opStack) Peek() string {
	l := len(*s)
	return (*s)[l-1]
}

func addNode(s *stack, v string) error {
	if len(*s) < 2 {
		return errors.New("invalid expression")
	}
	r := s.Pop()
	l := s.Pop()
	s.Push(binaryExpression{
		operator: v,
		left:     l,
		right:    r,
	})
	return nil
}

func (p *parser) Parse() (Expression, error) {

	var operatorStack opStack
	var operandStack stack

main:
	for {
		tok, lit := p.scanIgnoreWhitespace()

		switch tok {
		case tokEOF:
			break main
		case tokIllegal:
			return nil, fmt.Errorf("Illegal token: '%s'", lit)
		case tokLeftParenthesis:
			operatorStack.Push(lit)
		case tokRightParenthesis:
			for len(operatorStack) != 0 {
				popped := operatorStack.Pop()
				if popped == "(" {
					continue main
				} else {
					err := addNode(&operandStack, popped)
					if err != nil {
						return nil, err
					}
				}
			}
			return nil, errors.New("Unbalanced right parenthesis")
		case tokOperator:
			o1 := lit
			for len(operatorStack) > 0 {
				o2 := operatorStack.Peek()

				if (!isRightAssociative(o1) && precedence(o1) == precedence(o2)) || precedence(o1) < precedence(o2) {
					operatorStack.Pop()
					err := addNode(&operandStack, o2)
					if err != nil {
						return nil, err
					}
				} else {
					break
				}
			}
			operatorStack.Push(o1)
		case tokString:
			operandStack.Push(stringExpression(lit))
		case tokInt:
			i, err := strconv.Atoi(lit)
			if err != nil {
				return nil, err
			}
			operandStack.Push(intExpression(i))
		case tokFloat:
			f, err := strconv.ParseFloat(lit, 64)
			if err != nil {
				return nil, err
			}
			operandStack.Push(floatExpression(f))
		case tokIdentifier:
			operandStack.Push(identExpression(lit))
		}
	}

	for len(operatorStack) > 0 {
		err := addNode(&operandStack, operatorStack.Pop())
		if err != nil {
			return nil, err
		}
	}

	if len(operandStack) == 0 || len(operandStack) > 1 {
		return nil, errors.New("invalid syntax")
	}

	return operandStack.Pop(), nil
}
