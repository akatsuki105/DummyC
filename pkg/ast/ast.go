package ast

import (
	"bytes"
	"dummyc/pkg/token"
	"strings"
)

// Node - The base Node interface
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement - All statement nodes implement this
type Statement interface {
	Node
	statementNode()
}

// Expression - All expression nodes implement this
type Expression interface {
	Node
	expressionNode()
}

// TranslationUnit - Root node
type TranslationUnit struct {
	Prototypes []Prototype
	Functions  []FunctionLiteral
}

// TokenLiteral - Return TokenLiteral in first statement
func (tu *TranslationUnit) TokenLiteral() string {
	if len(tu.Prototypes) > 0 {
		return tu.Prototypes[0].TokenLiteral()
	} else {
		return ""
	}
}

// String - Print all statement as string
func (tu *TranslationUnit) String() string {
	var out bytes.Buffer

	for _, p := range tu.Prototypes {
		out.WriteString(p.String())
	}

	for _, f := range tu.Functions {
		out.WriteString(f.String())
	}

	return out.String()
}

// InfixExpression - Infix node e.g. +,-,*,/
type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode()      {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")

	return out.String()
}

// CallExpression - call expression node
type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Identifier  // Function Identifier
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// ReturnStatement - Return Statement
type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// FunctionStatement - function statement fn {...}
type FunctionStatement struct {
	Token        token.Token // the { token
	Declarations []DeclarationStatement
	Statements   []Statement
}

func (fs *FunctionStatement) statementNode()       {}
func (fs *FunctionStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *FunctionStatement) String() string {
	var out bytes.Buffer

	for _, d := range fs.Declarations {
		out.WriteString(d.String())
	}

	for _, s := range fs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// ExpressionStatement - Expression node
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
