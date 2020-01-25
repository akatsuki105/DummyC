package ast

import (
	"../token"
	"bytes"
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
func (pt *Prototype) GetName() string { return pt.Name.Name() }
func (pt *Prototype) GetParamNum() int { return len(pt.Parameters) }

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
func (ft *FunctionLiteral) GetName() string { return ft.Prototype.GetName() }
func (ft *FunctionLiteral) GetParamNum() int { return ft.Prototype.GetParamNum() }