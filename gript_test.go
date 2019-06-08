package gript

import "testing"

type testCase struct {
	expression string
	variables  map[string]interface{}
	expected   interface{}
}

func testEval(t *testing.T, testCases []testCase) {

	for _, testCase := range testCases {
		result, err := Eval(testCase.expression, testCase.variables)

		if err != nil {
			t.Errorf("%s : error: %s", testCase.expression, err)
			continue
		}

		if result != testCase.expected {
			t.Errorf("%s : invalid result. Got %+v, expected %+v", testCase.expression, result, testCase.expected)
		}
	}
}

func TestEvalConstants(t *testing.T) {

	testEval(t, []testCase{
		{"0", nil, 0},
		{"1", nil, 1},
		{"-1", nil, -1},
		{"12.3", nil, 12.3},
		{"'abc'", nil, "abc"},
		{"''", nil, ""},
		{"`'`", nil, `'`},
	})
}

func TestEvalIntComparison(t *testing.T) {

	testEval(t, []testCase{
		{"2 == 3", nil, false},
		{"2 != 3", nil, true},
		{"2 >= 3", nil, false},
		{"2 > 3", nil, false},
		{"2 <= 3", nil, true},
		{"2 < 3", nil, true},
		{"2 == 2", nil, true},
		{"2 != 2", nil, false},
		{"2 >= 2", nil, true},
		{"2 > 2", nil, false},
		{"2 <= 2", nil, true},
		{"2 < 2", nil, false},
	})
}

func TestEvalFloatComparison(t *testing.T) {

	testEval(t, []testCase{
		{"2.2 == 3.2", nil, false},
		{"2.2 != 3.2", nil, true},
		{"2.2 >= 3.2", nil, false},
		{"2.2 > 3.2", nil, false},
		{"2.2 <= 3.2", nil, true},
		{"2.2 < 3.2", nil, true},
		{"2.2 == 2.2", nil, true},
		{"2.2 != 2.2", nil, false},
		{"2.2 >= 2.2", nil, true},
		{"2.2 > 2.2", nil, false},
		{"2.2 <= 2.2", nil, true},
		{"2.2 < 2.2", nil, false},
	})
}

func TestEvalStringComparison(t *testing.T) {

	testEval(t, []testCase{
		{"'abc' == 'abd'", nil, false},
		{"'abc' != 'abd'", nil, true},
		{"'abc' >= 'abd'", nil, false},
		{"'abc' > 'abd'", nil, false},
		{"'abc' <= 'abd'", nil, true},
		{"'abc' < 'abd'", nil, true},
		{"'abc' == 'abc'", nil, true},
		{"'abc' != 'abc'", nil, false},
		{"'abc' >= 'abc'", nil, true},
		{"'abc' > 'abc'", nil, false},
		{"'abc' <= 'abc'", nil, true},
		{"'abc' < 'abc'", nil, false},
	})
}

func TestEvalArithmetic(t *testing.T) {

	testEval(t, []testCase{
		{"1 + 1", nil, 2},
		{"1.2 + 1.3", nil, 2.5},
		{"'ab' + 'c'", nil, "abc"},
		{"1 - 3", nil, -2},
		{"1.3 -3.1", nil, -1.8},
		{"3 * 4", nil, 12},
		{"3. * 4.", nil, 12.},
		{"6 / 2", nil, 3},
		{"5 / 2", nil, 2},
		{"7 / 2", nil, 3},
		{"-7 / 2", nil, -3},
		{"6. / 2.", nil, 3.},
		{"5. / 2.", nil, 2.5},
		{"6 % 2", nil, 0},
		{"6 % 5", nil, 1},
	})
}

func TestEvalAccessObject(t *testing.T) {
	testEval(t, []testCase{
		{"payload.a", map[string]interface{}{"payload": map[string]interface{}{"a": 1}}, 1},
	})
}

func TestEvalComplex(t *testing.T) {

	testEval(t, []testCase{
		{"a", map[string]interface{}{"a": 1}, 1},
		{"a", map[string]interface{}{"a": nil}, nil},
		{"a == nil", map[string]interface{}{"a": nil}, true},
		{"true", nil, true},
		{"false", nil, false},
		{"a > 4 || (a < 2 && a > 0) || a == 6", map[string]interface{}{"a": 1}, true},
		{"ab > 3+1   ||	(ab < 4-2 && ab > 6%2) ", map[string]interface{}{"ab": 1}, true},
	})
}

func TestEvalInvalidSyntax(t *testing.T) {

	testCases := []struct {
		expression string
		variables  map[string]interface{}
		err        string
	}{
		{"", nil, "invalid syntax"},
		{"1)", nil, "Unbalanced right parenthesis"},
		{"(1", nil, "invalid expression"},
		{"#", nil, "Illegal token: '#'"},
		{"1 >> 2", nil, "Unsupported operator '>>'"},
		{"1+", nil, "invalid expression"},
		{"(1+)", nil, "invalid expression"},
		{"9223372036854775808", nil, "strconv.Atoi: parsing \"9223372036854775808\": value out of range"},
		{"1.1.", nil, "Illegal token: '1.1.'"},
		{"'a", nil, "Illegal token: 'a'"},
		{"a", nil, "undefined variable 'a'"},
		{"a &&  || b", nil, "invalid expression"},
		{"a.b", nil, "undefined variable 'a.b'"},
		{"a.b", map[string]interface{}{"a":1}, "undefined variable 'a.b'"},
		{"a.b", map[string]interface{}{"a":map[string]interface{}{"c":1}}, "undefined variable 'a.b'"},
	}

	for _, testCase := range testCases {
		result, err := Eval(testCase.expression, testCase.variables)

		if err == nil || err.Error() != testCase.err {
			t.Errorf("%s : expecting error %s, got %+v", testCase.expression, testCase.err, err)
			continue
		}

		if result != nil {
			t.Errorf("%s : invalid result. Got %+v, expected %+v", testCase.expression, result, nil)
		}
	}
}
func TestEvalInvalidTypes(t *testing.T) {

	testCases := []struct {
		expression string
		variables  map[string]interface{}
		err        string
	}{
		{"'a' || 0>1", nil, "boolean expected in OR expression"},
		{"0>1 || 'a'", nil, "boolean expected in OR expression"},
		{"'a' && 0>1", nil, "boolean expected in AND expression"},
		{"0>1 && 'a'", nil, "boolean expected in AND expression"},
		{"'a' > 0", nil, "incompatible types in comparison"},
		{"'a' + 0", nil, "incompatible types in sum"},
		{"'a' - 0", nil, "incompatible types in difference"},
		{"'a' * 0", nil, "incompatible types in product"},
		{"'a' / 0", nil, "incompatible types in quotient"},
		{"'a' % 0", nil, "incompatible types in modulo"},
		{"'a' > 0 || 1>0", nil, "incompatible types in comparison"},
		{"1>0 && 'a' > 0", nil, "incompatible types in comparison"},
	}

	for _, testCase := range testCases {
		_, err := Eval(testCase.expression, testCase.variables)

		if err == nil || err.Error() != testCase.err {
			t.Errorf("%s : expecting error %s, got %+v", testCase.expression, testCase.err, err)
			continue
		}
	}
}

var result interface{}

func BenchmarkEvalBasic(b *testing.B) {
	var r interface{}
	for n := 0; n < b.N; n++ {
		r, _ = Eval("1", nil)
	}
	result = r
}
func BenchmarkEvalComplex(b *testing.B) {
	var r interface{}
	for n := 0; n < b.N; n++ {
		r, _ = Eval("a > 4 || (a < 2 && a > 0)", map[string]interface{}{"a": 1})
	}
	result = r
}
