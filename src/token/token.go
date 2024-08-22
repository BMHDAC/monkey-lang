package token

/*
  Let's take this source code as example
  ```cpp
    let x = 0 + 5;
  ```
  It will come out of the lexer like this:
  [
    LET,
    IDENTIFIER("x"),
    EQUAL_SIGN,
    INTERGER(0),
    PLUS,
    INTERGER(5),
    SEMICOLON
  ]
*/

// Let's just take our token's type as a string for now
type TokenType string

// Each token will have a type, and their respective literal
type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

func LookUpIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// Define our basic token types
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifier
	IDENT = "ident"
	INT   = "INT"

	// Operators
	ASSIGN = "="
	PLUS   = "+"

	// Delimiter
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)
