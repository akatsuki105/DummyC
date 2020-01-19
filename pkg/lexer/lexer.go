package lexer

import (
	"bytes"
	"dummyc/pkg/token"
)

type Lexer struct {
	tokens   []*token.Token
	curIndex int
}

func (lexer *Lexer) UngetToken(times int) bool {
	for i := 0; i < times; i++ {
		if lexer.curIndex == 0 {
			return false
		}
		lexer.curIndex--
	}

	return true
}

func (lexer *Lexer) GetNextToken() bool {
	if lexer.curIndex+1 >= len(lexer.tokens) {
		return false
	}

	lexer.curIndex++
	return true
}

func (lexer *Lexer) PushToken(tok *token.Token) bool {
	lexer.tokens = append(lexer.tokens, tok)
	return true
}

func (lexer *Lexer) GetToken() token.Token {
	return *lexer.tokens[lexer.curIndex]
}

func (lexer *Lexer) GetCurType() token.TokenType {
	return lexer.tokens[lexer.curIndex].Type
}

func (lexer *Lexer) GetCurString() string {
	return lexer.tokens[lexer.curIndex].Literal
}

func (lexer *Lexer) GetCurNumVal() int {
	return lexer.tokens[lexer.curIndex].Number
}

func (lexer *Lexer) GetCurIndex() int {
	return lexer.curIndex
}

func (lexer *Lexer) ApplyTokenIndex(index int) bool {
	lexer.curIndex = index
	return true
}

func (lexer *Lexer) String() string {
	var out bytes.Buffer

	out.WriteString("[")
	for _, token := range lexer.tokens {
		out.WriteString(token.Literal + ", ")
	}
	out.WriteString("]")
	return out.String()
}

// --------------------------------------------------------- constructor --------------------------------------------------------

func New(source string) *Lexer {
	lexer := &Lexer{}
	line := 0
	skip := 0

	for i, char := range source {

		if skip > 0 {
			skip--
			continue
		}
		if isWhitespace(char) {
			continue
		}

		switch char {
		case '=':
			lexer.PushToken(token.New(token.ASSIGN, string(char), line))
		case '+':
			lexer.PushToken(token.New(token.PLUS, string(char), line))
		case '-':
			lexer.PushToken(token.New(token.MINUS, string(char), line))
		case '*':
			lexer.PushToken(token.New(token.ASTERISK, string(char), line))
		case '/':
			lexer.PushToken(token.New(token.SLASH, string(char), line))
		case ',':
			lexer.PushToken(token.New(token.COMMA, string(char), line))
		case ';':
			lexer.PushToken(token.New(token.SEMICOLON, string(char), line))
		case '(':
			lexer.PushToken(token.New(token.LPAREN, string(char), line))
		case ')':
			lexer.PushToken(token.New(token.RPAREN, string(char), line))
		case '{':
			lexer.PushToken(token.New(token.LBRACE, string(char), line))
		case '}':
			lexer.PushToken(token.New(token.RBRACE, string(char), line))
		case '[':
			lexer.PushToken(token.New(token.LBRACKET, string(char), line))
		case ']':
			lexer.PushToken(token.New(token.RBRACKET, string(char), line))
		case ':':
			lexer.PushToken(token.New(token.COLON, string(char), line))
		case '\n':
			line++
		default:
			switch {
			case isLetter(char):
				identifier := string(char)
				for j := 1; ; j++ {
					char := rune(source[i+j])
					if isLetter(char) {
						identifier += string(char)
						skip++
					} else {
						break
					}
				}

				switch identifier {
				case "int":
					lexer.PushToken(token.New(token.INTTYPE, identifier, line))
				case "return":
					lexer.PushToken(token.New(token.RETURN, identifier, line))
				default:
					lexer.PushToken(token.New(token.IDENT, identifier, line))
				}

			case isDigit(char):
				number := string(char)
				for j := 1; ; j++ {
					char := rune(source[i+j])
					if isDigit(char) {
						number += string(char)
						skip++
					} else {
						break
					}
				}
				lexer.PushToken(token.New(token.DIGIT, number, line))
			default:
				lexer.PushToken(token.New(token.ILLEGAL, string(char), line))
			}
		}
	}

	return lexer
}

func isWhitespace(char rune) bool {
	return char == ' ' || char == '\t' || char == '\r'
}

func isLetter(char rune) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func isDigit(char rune) bool {
	return '0' <= char && char <= '9'
}
