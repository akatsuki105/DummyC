package token

import "strconv"

type TokenType string

const (
	tokIllegal    = "illegal" // トークンタイプ不明
	tokIdentifier = "identifier"
	tokDigit      = "digit"
	tokSymbol     = "symbol"
	tokInt        = "int"
	tokReturn     = "return"
	tokEOF        = "eof"
)

type Token struct {
	tokenType   TokenType
	tokenString string // トークンの文字列表現
	number      int    // トークンの数値表現
	line        int    // トークンが出てきた行数
}

func New(tokenType TokenType, tokenString string, line int) *Token {
	var number int
	var err error
	if tokenType == tokDigit {
		number, err = strconv.Atoi(tokenString)
		if err != nil {
			number = 0x7fffffff
		}
	} else {
		number = 0x7fffffff
	}

	return &Token{
		tokenType:   tokenType,
		tokenString: tokenString,
		number:      number,
		line:        line,
	}
}

func (tok *Token) GetTokenType() TokenType {
	return tok.tokenType
}

func (tok *Token) GetTokenString() string {
	return tok.tokenString
}

func (tok *Token) GetNumberValue() int {
	return tok.number
}

func (tok *Token) GetLine() int {
	return tok.line
}
