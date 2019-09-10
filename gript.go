package gript

import (
	"bytes"
	"reflect"
	"strings"
)

//Context is an interface allowing access to variable values
type Context interface {
	Value(identifier string) (value interface{}, found bool)
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

func (vm *vm) Value(ident string) (interface{}, bool) {

	parts := strings.Split(ident, ".")
	var current interface{}
	current = vm.values

	for i := 0; i < len(parts); i++ {

		if currentMap, ok := current.(map[string]interface{}); ok {
			next, found := currentMap[parts[i]]
			if !found {
				return nil, false
			}
			current = next
		} else {
			currentValue := reflect.ValueOf(current)
			if currentValue.Kind() != reflect.Struct {
				return nil, false
			}
			nextValue := currentValue.FieldByNameFunc(func(name string) bool {
				return strings.ToLower(name) == strings.ToLower(parts[i])
			})
			if (nextValue == reflect.Value{}) {
				return nil, false
			}
			current = nextValue.Interface()
		}
	}

	return current, true
}

func (vm *vm) Eval(s string) (interface{}, error) {

	exp, err := Parse(s)
	if err != nil {
		return nil, err
	}
	return exp.Eval(vm)
}
