package parser

import (
	"monkey/src/ast"
	"monkey/src/lexer"
	"testing"
)

func TestLetStatement(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 8472;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	if program == nil {
		t.Fatalf("ParseProgram() returned: nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statement: Wrong number of statement. Expected 3, got: %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]

		if !testLetStatement(t, statement, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral() returned: %q, expected: `let`", s.TokenLiteral())
		return false
	}

	stm, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement, got %T", s)
		return false
	}

	if stm.Name.Value != name {
		t.Errorf("stm.Name.Value returned: %s. Expected: %s", stm.Name.Value, name)
		return false
	}

	if stm.Name.TokenLiteral() != name {
		t.Errorf("stm.Name.TokenLiteral() returned: %s. Expected: %s", stm.Name, name)
		return false
	}

	return true
}
