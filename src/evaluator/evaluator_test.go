package evaluator

import (
	"testing"

	"monkey/src/lexer"
	"monkey/src/object"
	"monkey/src/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"42069", 42069},
		{"-5", -5},
		{"-10", -10},
		{"-42069", -42069},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, evaluated object.Object, expected int64) bool {
	result, ok := evaluated.(*object.Integer)
	if !ok {
		t.Errorf("Object is not Integer Object, got: %T", evaluated)
		return false
	}

	if result.Value != expected {
		t.Errorf("Integer Object value wrong, expected: %d, got: %d", expected, result.Value)
		return false
	}

	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			"true", true,
		},
		{
			"false", false,
		},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, evaluated object.Object, expected bool) bool {
	result, ok := evaluated.(*object.Boolean)
	if !ok {
		t.Errorf("Object is not Boolean, got: %T", evaluated)
		return false
	}

	if result.Value != expected {
		t.Errorf("Object value wrong, expected: %t, got: %t", expected, result.Value)
		return false
	}

	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		{"!5", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("Object is not NULL, got: %T (%+v)", obj, obj)
		return false
	} else {
		return true
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"if (10 > 1) { if (10 > 1) { return 10;} return 1;}", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true",
			"type missmatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5",
			"type missmatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operation: -BOOLEAN",
		},
		{
			"true + false",
			"unknown operation: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operation: BOOLEAN + BOOLEAN",
		},

		{
			"if (10 > 1) { true + false;}",
			"unknown operation: BOOLEAN + BOOLEAN",
		},
		{
			`if (10 > 1) {
        if (10 > 1) {
          return true + false;
        }

        return 1;
      }`,
			"unknown operation: BOOLEAN + BOOLEAN",
		},
		{
			"foobar;",
			"identifier not found: `foobar`",
		},
		{
			`"Hello" - "World"`,
			"unknown operation: STRING - STRING",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errorObject, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("No error object returned, got: %T (%+v)", evaluated, evaluated)
			continue
		}

		if errorObject.Message != tt.expectedMessage {
			t.Errorf("Wrong error message. Expected: %s, got: %s", tt.expectedMessage, errorObject.Message)
		}
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			"let x = 10 + 5; x;",
			15,
		},
		{
			"let x = 10 * 5; x;",
			50,
		},
		{
			"let x = 10 - 5; x;",
			5,
		},
		{
			"let x = 10 / 5; x;",
			2,
		},

		{
			"let a = 5; let b = a; b;",
			5,
		},
		{
			"let a = 5; let b = a; let c = a + b + 5; c",
			15,
		},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2;};"

	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)

	if !ok {
		t.Errorf("Evaluated input is not Function, got: %T(%+v) ", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Errorf("Function parameters wrong, expeted: %d, got: %d", 1, len(fn.Parameters))
	}

	if fn.Parameters[0].String() != "x" {
		t.Errorf("Function parameters identifier wrong, expected: %s, got: %s", "x", fn.Parameters[0].String())
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("Function body expression wrong, expected: %s, got: %s", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identify = fn(x) {return x;}; identify(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosure(t *testing.T) {
	input := `
    let newAdder = fn(x) {
      fn(y) { x + y };
    };
    let addTwo = newAdder(3);
    addTwo(2);
  `

	testIntegerObject(t, testEval(input), 5)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello Worlds!";`

	evaluated := testEval(input)

	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("Input is not a string, got: %T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello Worlds!" {
		t.Fatalf("Str value wrong. Expected: %s, got: %s", "Hello Worlds!", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " "  + "worlds!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("evaluated is not String, got: %T (%+v) ", evaluated, evaluated)
	}

	if str.Value != "Hello worlds!" {
		t.Fatalf("stf.Value wrong. Expected: %s, got: %s", "Hello worlds!", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello worlds")`, 12},
		{`len(1)`, "argument to `len` not supported: INTEGER"},
		{`len(1,2)`, "wrong number of arguments. Got: 2, take: 1"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errorObject, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("Expected Error Object, got: %T (%+v)", evaluated, evaluated)
				continue
			}
			if errorObject.Message != expected {
				t.Errorf("wrong error message, expected: %s, got: %s", expected, errorObject.Message)
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2, 2 * 2, 3, 3 + 2]"
	evaluated := testEval(input)

	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("evaluated is not Array object, got: %T (%+v)", evaluated, evaluated)
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 2)
	testIntegerObject(t, result.Elements[2], 4)
	testIntegerObject(t, result.Elements[3], 3)
	testIntegerObject(t, result.Elements[4], 5)
}

func TestArrayIndexExpression(t *testing.T) {
	tests := []struct {
		input     string
		expeceted interface{}
	}{
		{
			"[0, 1, 2][0]",
			0,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expeceted.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}
