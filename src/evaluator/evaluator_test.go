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
	return Eval(program)
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
