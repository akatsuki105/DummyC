package ast

import (
	"bytes"
	"dummyc/pkg/token"
	"strings"
)

// Prototype - Prototype declaration
type Prototype struct {
	Token      token.Token // the token.INTTYPE token
	Name       *Identifier
	Parameters []*Identifier
}

func (pt *Prototype) expressionNode()      {}
func (pt *Prototype) TokenLiteral() string { return pt.Name.TokenLiteral() }
func (pt *Prototype) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range pt.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("int ")
	out.WriteString(pt.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")

	return out.String()
}

// FunctionLiteral - function node
type FunctionLiteral struct {
	Token     token.Token
	Prototype Prototype
	Body      FunctionStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	out.WriteString(fl.Prototype.String())
	out.WriteString(fl.Body.String())

	return out.String()
}
