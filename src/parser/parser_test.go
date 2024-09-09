package parser

import (
	"fmt"
	"testing"

	"monkey/src/ast"
	"monkey/src/lexer"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserError(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}
		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
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
		integerValue interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
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

		if !testLiteralExpression(t, pExp.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input    string
		leftVal  interface{}
		rightVal interface{}
		operator string
	}{
		{"5 + 5", 5, 5, "+"},
		{"5 - 5", 5, 5, "-"},
		{"5 * 5", 5, 5, "*"},
		{"5 / 5", 5, 5, "/"},
		{"5 > 5", 5, 5, ">"},
		{"5 < 5", 5, 5, "<"},
		{"5 == 5", 5, 5, "=="},
		{"5 != 5", 5, 5, "!="},
		{"true == true", true, true, "=="},
		{"true != false", true, false, "!="},
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

		if !testLiteralExpression(t, exp.Left, tt.leftVal) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator not %s, got %s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.rightVal) {
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
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
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
	case bool:
		{
			return testBooleanLiteral(t, exp, v)
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

func TestBooleanExpression(t *testing.T) {
	input := "true"

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

	iden, ok := stmt.Expression.(*ast.Boolean)
	if !ok {
		t.Fatalf("not identifier, got %T", stmt.Expression)
	}

	if iden.Value != true {
		t.Fatalf("iden.Value not `true`, got %v", iden.Value)
	}

	if iden.TokenLiteral() != "true" {
		t.Fatalf("iden.TokenLiteral() not `true`, got %s", iden.TokenLiteral())
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserError(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("Expected: %s, got: %s", actual, tt.expected)
		}
	}
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not Boolean, got: %T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not `%v`, got: %v", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%v", value) {
		t.Errorf("bo.Literal() not %t, got: %s", value, bo.TokenLiteral())
		return false
	}

	return true
}

func TestIfExpression(t *testing.T) {
	input := "if (x < y) {x}"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("len(program.Statement) is not 1, got: %d", len(program.Statements))
	}

	stm, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement, got %T", program.Statements[0])
	}

	exp, ok := stm.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stm.Expression is not IfExpression, got %T", stm.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1, got: %d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("exp.Consequence.Statements[0] is not ExpressionStatement, got: %T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative not nil, got %T", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
			len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := "fn(x,y) { x + y; }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("len(program.Statements) is not 1, got: %d", len(program.Statements))
	}

	stm, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not ExpressionStatement, got: %T", program.Statements[0])
	}

	function, ok := stm.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stm.Expression is not FunctionLiteral, got: %T", stm.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function.Parameters len is not 2, got: %d\n", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements is not 1, got %d\n", len(function.Body.Statements))
	}

	bodyStm, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function.Body.Statements[0] is not ExpressionStatement, got %T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStm.Expression, "x", "+", "y")
}

func TestFunctionParamsParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{
			"fn() {}",
			[]string{},
		},
		{
			"fn(x) {}",
			[]string{"x"},
		},
		{
			"fn(x,y,z) {}",
			[]string{"x", "y", "z"},
		},
	}

	for _, tt := range tests {
		input := tt.input

		l := lexer.New(input)
		p := New(l)
		program := p.ParseProgram()
		checkParserError(t, p)
		stm := program.Statements[0].(*ast.ExpressionStatement)
		function := stm.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("function.Parameters len wrong, expected: %d, got: %d", len(tt.expectedParams), len(function.Parameters))
		}

		for i, iden := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], iden)
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("len(program.Statements) is not 1, got: %d\n", len(program.Statements))
	}

	stm, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement, got: %T", program.Statements[0])
	}

	exp, ok := stm.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stm.Expression is not CallExpression, got %T", stm.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("exp.Arguments not 3, got: %d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello worlds!";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserError(t, p)

	stm, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Statement is not Expression, got: %T", program.Statements[0])
	}

	literal, ok := stm.Expression.(*ast.StringLiteral)

	if !ok {
		t.Fatalf("stm.Expression is not StringLiteral, got: %T", stm.Expression)
	}

	if literal.Value != "hello worlds!" {
		t.Fatalf("wrong string literal value, expected: %s, got: %s", "hello worlds!", literal.Value)
	}
}

func TestArrayLiteralParsing(t *testing.T) {
	input := "[1, 2, 2 * 2, 3, 3 + 2]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)

	stm, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not ExpressionStatement, got: %T", program.Statements[0])
	}

	array, ok := stm.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("stm.Expression is not ArrayLiteral, got: %T", stm.Expression)
	}

	if len(array.Elements) != 5 {
		t.Fatalf("array.Elements length wrong, expected: 5, got: %d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[2], 2, "*", 2)
	testInfixExpression(t, array.Elements[4], 3, "+", 2)
}

func TestIndexExpressionParsing(t *testing.T) {
	input := `myArray[1 + 1]`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)

	stm, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not ExpressionStatement, got: %T", program.Statements[0])
	}

	indexExp, ok := stm.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("stm.Expression is not IndexExpression, got: %T", stm.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestHashingLiteral(t *testing.T) {
	input := `{"one" : 1, "two": 2, "three": 3 }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)

	stm := program.Statements[0].(*ast.ExpressionStatement)

	hash, ok := stm.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("stm.Expression is not HashLiteral, got: %T", stm.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Fatalf("hash.Pairs len is not 3, got: %d", len(hash.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Fatalf("hash.Pairs element key not StringLiteral, got: %T (%+v)", key, key)
		}

		expectedValue := expected[literal.String()]

		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingEmptyHash(t *testing.T) {
	input := "{}"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)

	stm := program.Statements[0].(*ast.ExpressionStatement)

	hash, ok := stm.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("stm.Expression is not HashLiteral, got: %T", stm.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Fatalf("hash.Pairs len is not 0, got: %d", len(hash.Pairs))
	}
}

func TestParsingLiteralWithTypes(t *testing.T) {
	input := `{"one" : 1 + 0, "two": 10 -8, "three": 15 / 5}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)

	stm := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stm.Expression.(*ast.HashLiteral)

	if !ok {
		t.Fatalf("stm.Expression is not HashLiteral, got: %T", stm.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs len is not 3, got: %d", len(hash.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 1, "+", 0)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, val := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not StringLiteral, got: %T", key)
			continue
		}

		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key: %s", key.String())
			continue
		}

		testFunc(val)
	}
}
