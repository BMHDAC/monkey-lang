package parser

import (
	"monkey/src/ast"
	"monkey/src/lexer"
	"monkey/src/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {}

func (p *Parser) parseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stm := p.parseStatement()
		if stm != nil {
			program.Statements = append(program.Statements, stm)
		}
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() ast.Statement {
	stm := &ast.LetStatement{
		Token: p.curToken,
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stm.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stm
}

func (p *Parser) peekTokenIs(token token.TokenType) bool {
	return p.peekToken.Type == token
}

func (p *Parser) curTokenIs(token token.TokenType) bool {
	return p.curToken.Type == token
}

func (p *Parser) expectPeek(token token.TokenType) bool {
	if p.peekTokenIs(token) {
		p.nextToken()
		return true
	} else {
		return false
	}
}
