package ast

import (
	"../token"
	"bytes"
)

const (
	Param = "param"
	Local = "local"
)

// DeclarationStatement - Varaiable declaration statement
type DeclarationStatement struct {
	Token    token.Token // the token.INTTYPE token
	Name     Identifier
	declType string
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
func (ls *DeclarationStatement) GetDeclType() string {
	return ls.declType
}
func (ls *DeclarationStatement) SetDeclType(dt string) {
	if dt == Param || dt == Local {
		ls.declType = dt
	}
}
