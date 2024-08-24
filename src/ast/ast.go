package ast

import "monkey/src/token"

// A token node of each parsed input
type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) == 0 {
		return ""
	} else {
		return p.Statements[0].TokenLiteral()
	}
}

// Node of a let statement
type LetStatement struct {
	Token token.Token // The LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// Node of an identifier
type Identifier struct {
	Token token.Token // The IDEN token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
