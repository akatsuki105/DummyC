package lexer

import (
	"dummyc/pkg/token"
	"fmt"
)

type TokenStream struct {
	tokens   []*token.Token
	curIndex int
}

func New() *TokenStream {
	return &TokenStream{}
}

func (ts *TokenStream) GetCurType() token.TokenType {
	tok := ts.tokens[ts.curIndex]
	return tok.GetTokenType()
}

func (ts *TokenStream) GetCurString() string {
	tok := ts.tokens[ts.curIndex]
	return tok.GetTokenString()
}

func (ts *TokenStream) GetCurNumVal() int {
	tok := ts.tokens[ts.curIndex]
	return tok.GetNumberValue()
}

func (ts *TokenStream) ApplyTokenIndex(index int) bool {
	if index >= 0 && index+1 < len(ts.tokens) {
		ts.curIndex = index
		return true
	}
	return false
}

func (ts *TokenStream) PrintTokens() {
	fmt.Println(ts.tokens)
}

func (ts *TokenStream) GetToken() token.Token {
	return *ts.tokens[ts.curIndex]
}

func (ts *TokenStream) GetNextToken() bool {
	if ts.curIndex+1 >= len(ts.tokens) {
		return false
	}
	ts.curIndex++
	return true
}

func (ts *TokenStream) UngetToken(times int) bool {
	for i := 0; i < times; i++ {
		if ts.curIndex == 0 {
			return false
		}
		ts.curIndex--
	}
	return true
}
