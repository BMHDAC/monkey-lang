package parser

import (
	"fmt"
	"testing"

	"monkey/src/ast"
	"monkey/src/lexer"
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

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTest := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTest {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 elements, got %d", len(program.Statements))
		}

		stm, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ExpressionStatement, got %T", program.Statements[0])
		}

		pExp, ok := stm.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stm.Expression not PrefixExpression, got %T", stm.Expression)
		}

		if pExp.Operator != tt.operator {
			t.Fatalf("pExp.Operator not `%s`, got `%s` instead", tt.operator, pExp.Operator)
		}

		if !testIntegerLiteral(t, pExp.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfOperator(t *testing.T) {
	infixTests := []struct {
		input    string
		leftVal  int64
		rightVal int64
		operator string
	}{
		{"5 + 5", 5, 5, "+"},
		{"5 - 5", 5, 5, "-"},
		{"5 * 5", 5, 5, "*"},
		{"5 / 5", 5, 5, "/"},
		{"5 > 5", 5, 5, ">"},
		{"5 < 5", 5, 5, "<"},
		{"5 == 5", 5, 5, "=="},
		{"5 == 5", 5, 5, "=="},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("len(program.Statements) not 1, got: %d", len(program.Statements))
		}

		stm, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ExpressionStatement, got %T", program.Statements[0])
		}

		exp, ok := stm.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stm.Expression not InfixExpression, got %T", stm.Expression)
		}

		if !testIntegerLiteral(t, exp.Left, tt.leftVal) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator not %s, got %s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.rightVal) {
			return
		}

	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("il not IntegerLiteral, got %T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("il.Value not %d, got %d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("il.TokenLiteral() not %d, got %s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4;-5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserError(t, p)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("Expected: %q, got: %q", tt.expected, actual)
		}
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

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)

	if !ok {
		t.Errorf("epx is not identifier, got %T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value is not: %s, got: %s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() is not %s, got: %s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		{
			return testIntegerLiteral(t, exp, int64(v))
		}
	case int64:
		{
			return testIntegerLiteral(t, exp, v)
		}
	case string:
		{
			return testIdentifier(t, exp, v)
		}
	}

	t.Errorf("type of exp not handled. Got: %T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp not ast.OperatorExpression, got: %T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not `%s`, got: %q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}
