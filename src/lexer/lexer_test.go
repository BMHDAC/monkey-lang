package lexer

import (
	"monkey/src/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
    let five = 5;
    let ten = 10;
    let add = fn(x,y) {
      x + y;
    };
    
    let result = add(five, ten);
  `

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
	}

	lexer := New(input)

	for i, token := range tests {
		tok := lexer.NextToken()

		if tok.Type != token.expectedType {
			t.Fatalf("Test [%d] type failed. Expected: %q, got: %q", i, token.expectedType, tok.Type)
		}
		if tok.Literal != token.expectedLiteral {
			t.Fatalf("Test [%d] literal failed. Expected: %q, got: %q", i, token.expectedLiteral, tok.Literal)
		}
	}
}
