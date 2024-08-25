package parser

import (
	"monkey/src/ast"
	"monkey/src/lexer"
	"testing"
)

func TestLetStatement(t *testing.T) {
	input := `
    let a 0;
    let = 10;
    let 696969;
  `

	l := lexer.New(input)
	p := New(l)

	program := p.parseProgram()
	checkParserError(t, p)

	if program == nil {
		t.Fatalf("parseProgram() returned nil")
	}

	if len(program.Statements) != 4 {
		t.Fatalf("Got %d statements, expected: %d", len(program.Statements), 4)
	}

	tests := []struct {
		exprectedIdentifier string
	}{
		{"a"},
		{"b"},
		{"c"},
		{"you_suck"},
	}

	for i, tt := range tests {
		stm := program.Statements[i]
		if !testLetStatement(t, stm, tt.exprectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf("Let token literal not `let`, got %s", statement.TokenLiteral())
		return false
	}

	letStm, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("Statement not LetStatement, got: %T", statement)
		return false
	}

	if letStm.Name.Value != name {
		t.Errorf("Let statement value is not %s, got: %s", name, letStm.Name.Value)
		return false
	}
	if letStm.Name.TokenLiteral() != name {
		t.Errorf("Let statement literal is not %s, got: %s", name, letStm.Name.TokenLiteral())
		return false
	}

	return true
}

func checkParserError(t *testing.T, p *Parser) {
	error := p.Errors()

	if len(error) == 0 {
		return
	}

	t.Errorf("Parser has %d errors", len(error))

	for _, error := range error {
		t.Errorf("parser error: %q", error)
	}

	t.FailNow()
}
