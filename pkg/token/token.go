package token

import "strconv"

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT = "IDENT" // add, foobar, x, y, ...
	DIGIT = "DIGIT" // 1343456

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"
	COLON    = ":"

	// Keywords
	INTTYPE = "INT"
	RETURN  = "RETURN"
)

type Token struct {
	Type    TokenType
	Literal string
	Number  int
	Line    int
}

func New(tokenType TokenType, literal string, line int) *Token {
	number := 0x7fffffff
	var err error
	if tokenType == DIGIT {
		number, err = strconv.Atoi(literal)
		if err != nil {
			number = 0x7fffffff
		}
	}

	return &Token{
		Type:    tokenType,
		Literal: literal,
		Number:  number,
		Line:    line,
	}
}
