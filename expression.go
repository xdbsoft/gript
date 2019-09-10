package gript

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type intExpression int

func (e intExpression) Eval(c Context) (interface{}, error) {
	return int(e), nil
}

type floatExpression float64

func (e floatExpression) Eval(c Context) (interface{}, error) {
	return float64(e), nil
}

type stringExpression string

func (e stringExpression) Eval(c Context) (interface{}, error) {
	return string(e), nil
}

type identExpression string

func (e identExpression) Eval(c Context) (interface{}, error) {

	switch e {
	case "true":
		return true, nil
	case "false":
		return false, nil
	case "nil":
		return nil, nil
	}

	r, found := c.Value(string(e))
	if !found {
		return nil, fmt.Errorf("undefined variable '%s'", e)
	}
	return r, nil
}

type binaryExpression struct {
	operator string
	left     Expression
	right    Expression
}

func or(l, r interface{}) (bool, error) {

	vl, ok := l.(bool)
	if !ok {
		return false, errors.New("boolean expected in OR expression")
	}
	vr, ok := r.(bool)
	if !ok {
		return false, errors.New("boolean expected in OR expression")
	}
	return vl || vr, nil
}
func and(l, r interface{}) (bool, error) {

	vl, ok := l.(bool)
	if !ok {
		return false, errors.New("boolean expected in AND expression")
	}
	vr, ok := r.(bool)
	if !ok {
		return false, errors.New("boolean expected in AND expression")
	}
	return vl && vr, nil
}

func less(l, r interface{}) (bool, error) {

	switch vl := l.(type) {
	case int:
		if vr, ok := r.(int); ok {
			return vl < vr, nil
		}
	case float64:
		if vr, ok := r.(float64); ok {
			return vl < vr, nil
		}
	case string:
		if vr, ok := r.(string); ok {
			return vl < vr, nil
		}
	}
	return false, errors.New("incompatible types in comparison")
}

func sum(l, r interface{}) (interface{}, error) {

	switch vl := l.(type) {
	case int:
		if vr, ok := r.(int); ok {
			return vl + vr, nil
		}
	case float64:
		if vr, ok := r.(float64); ok {
			return vl + vr, nil
		}
	case string:
		if vr, ok := r.(string); ok {
			return vl + vr, nil
		}
	}
	return nil, errors.New("incompatible types in sum")
}
func difference(l, r interface{}) (interface{}, error) {

	switch vl := l.(type) {
	case int:
		if vr, ok := r.(int); ok {
			return vl - vr, nil
		}
	case float64:
		if vr, ok := r.(float64); ok {
			return vl - vr, nil
		}
	}
	return nil, errors.New("incompatible types in difference")
}
func product(l, r interface{}) (interface{}, error) {
	switch vl := l.(type) {
	case int:
		if vr, ok := r.(int); ok {
			return vl * vr, nil
		}
	case float64:
		if vr, ok := r.(float64); ok {
			return vl * vr, nil
		}
	}
	return nil, errors.New("incompatible types in product")
}
func quotient(l, r interface{}) (interface{}, error) {

	switch vl := l.(type) {
	case int:
		if vr, ok := r.(int); ok {
			return vl / vr, nil
		}
	case float64:
		if vr, ok := r.(float64); ok {
			return vl / vr, nil
		}
	}
	return nil, errors.New("incompatible types in quotient")
}
func modulo(l, r interface{}) (interface{}, error) {

	switch vl := l.(type) {
	case int:
		if vr, ok := r.(int); ok {
			return vl % vr, nil
		}
	}
	return nil, errors.New("incompatible types in modulo")
}

func in(l, r interface{}) (interface{}, error) {

	rValue := reflect.ValueOf(r)
	lValue := reflect.ValueOf(l)

	switch rValue.Kind() {
	case reflect.Array, reflect.Slice:
		if !lValue.Type().AssignableTo(rValue.Type().Elem()) {
			return nil, errors.New("invalid type in operator in")
		}
		for i := 0; i < rValue.Len(); i++ {
			if rValue.Index(i).Interface() == lValue.Interface() {
				return true, nil
			}
		}
		return false, nil
	case reflect.Map:
		if !lValue.Type().AssignableTo(rValue.Type().Key()) {
			return nil, errors.New("invalid key type in operator in")
		}
		return rValue.MapIndex(lValue).IsValid(), nil
	case reflect.Struct:
		found := rValue.FieldByNameFunc(func(name string) bool {
			return strings.ToLower(name) == strings.ToLower(lValue.String())
		})
		return found.IsValid(), nil
	}
	return nil, errors.New("unsupported types in operator in")
}

func (e binaryExpression) Eval(c Context) (interface{}, error) {

	l, err := e.left.Eval(c)
	if err != nil {
		return nil, err
	}

	//Fast exit: in some cases, no need to compute right part of the expression
	switch e.operator {
	case "&&":
		vl, ok := l.(bool)
		if ok && !vl {
			return false, nil
		}
	case "||":
		vl, ok := l.(bool)
		if ok && vl {
			return true, nil
		}
	}
	r, err := e.right.Eval(c)
	if err != nil {
		return nil, err
	}

	switch e.operator {
	case "==":
		return l == r, nil
	case "!=":
		return l != r, nil
	case ">":
		return less(r, l)
	case ">=":
		if l == r {
			return true, nil
		}
		return less(r, l)
	case "<":
		return less(l, r)
	case "<=":
		if l == r {
			return true, nil
		}
		return less(l, r)
	case "||":
		return or(l, r)
	case "&&":
		return and(l, r)
	case "+":
		return sum(l, r)
	case "-":
		return difference(l, r)
	case "*":
		return product(l, r)
	case "/":
		return quotient(l, r)
	case "%":
		return modulo(l, r)
	case "in":
		return in(l, r)
	}
	return nil, fmt.Errorf("Unsupported operator '%s'", e.operator)
}
