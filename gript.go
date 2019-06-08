package gript

import (
	"bytes"
	"strings"
)

//Context is an interface allowing access to variable values
type Context interface {
	Value(identifier string) interface{}
}

//Expression (boolean, numerical, ...) is an obect that can be evaluated against a context
type Expression interface {
	Eval(c Context) (interface{}, error)
}

//Parse parses a string to create an expression
//
//Examples
// a < 2
// (a=='abc') || (b<c && (b+c >= 3.14))
func Parse(s string) (Expression, error) {
	b := bytes.NewBufferString(s)

	parser := newParser(b)
	return parser.Parse()
}

//Eval evaluates a string representing an expression against a set of variables
func Eval(s string, values map[string]interface{}) (interface{}, error) {
	vm := vm{values}
	return vm.Eval(s)
}

type vm struct {
	values map[string]interface{}
}

func (vm *vm) Value(ident string) interface{} {

	parts := strings.Split(ident, ".")
	current := vm.values

	for i := 0; i < len(parts)-1; i++ {
		v, found := current[parts[i]]
		if !found {
			return nil
		}

		next, ok := v.(map[string]interface{})
		if !ok {
			return nil
		}
		current = next
	}

	v, found := current[parts[len(parts)-1]]
	if !found {
		return nil
	}

	return v
}

func (vm *vm) Eval(s string) (interface{}, error) {

	exp, err := Parse(s)
	if err != nil {
		return nil, err
	}
	return exp.Eval(vm)
}
