package lexer

import (
	"monkey/src/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `={}()+,;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.PLUS, "+"},
		{token.COMMA, ","},
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
