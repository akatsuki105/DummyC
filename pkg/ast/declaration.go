package ast

import (
	"bytes"
	"dummyc/pkg/token"
)

// DeclarationStatement - Varaiable declaration statement
type DeclarationStatement struct {
	Token token.Token // the token.INTTYPE token
	Name  Identifier
}

func (ls *DeclarationStatement) statementNode()       {}
func (ls *DeclarationStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *DeclarationStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())

	out.WriteString(";")

	return out.String()
}
