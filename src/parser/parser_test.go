package parser

import (
	"monkey/src/ast"
	"monkey/src/lexer"
	"testing"
)

func TestLetStatement(t *testing.T) {
	input := `
    let a = 0;
    let b = 10;
    let foobar = 6969;
  `

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserError(t, p)

	if program == nil {
		t.Fatalf("parseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("Got %d statements, expected: %d", len(program.Statements), 3)
	}

	tests := []struct {
		exprectedIdentifier string
	}{
		{"a"},
		{"b"},
		{"foobar"},
	}

	for i, tt := range tests {
		stm := program.Statements[i]
		if !testLetStatement(t, stm, tt.exprectedIdentifier) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
    return a; 
    return b;
    return foobar;
  `

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserError(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("Got %d statements, expected %d", len(program.Statements), 3)
	}

	for _, stm := range program.Statements {
		rtStm, ok := stm.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Statement is not return statement, got %T instead", stm)
			continue
		}
		if rtStm.TokenLiteral() != "return" {
			t.Errorf("rtStm.TokenLiteral() not 'return', got %s instead", rtStm.TokenLiteral())
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

func TestIdentifierExpression(t *testing.T) {
	input := "foobar"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements not 1, got: %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not an Expression, got %T", program.Statements[0])
	}

	iden, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("not identifier, got %T", stmt.Expression)
	}

	if iden.Value != "foobar" {
		t.Fatalf("iden.Value not `foobar`, got %s", iden.Value)
	}

	if iden.TokenLiteral() != "foobar" {
		t.Fatalf("iden.TokenLiteral() not `foobar`, got %s", iden.TokenLiteral())
	}

}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("len(program.Statements) not 1, got %d", len(program.Statements))
	}

	stm, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not expression statement, got %T", program.Statements[0])
	}

	literal, ok := stm.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("stm.Expression not IntegerLiteral, got %T", stm.Expression)
	}

	if literal.Value != 5 {
		t.Fatalf("literal.Value not 5, got %d", literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Fatalf("literal.TokenLiteral() not 5, got %s", literal.TokenLiteral())
	}
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
