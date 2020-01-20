package ast

import "dummyc/pkg/token"

// Number - Interger literal node
type Number struct {
	Token token.Token
	Value int64
}

func (num *Number) expressionNode()      {}
func (num *Number) TokenLiteral() string { return num.Token.Literal }
func (num *Number) String() string       { return num.Token.Literal }

func (num *Number) Val() int {
	return int(num.Value)
}
