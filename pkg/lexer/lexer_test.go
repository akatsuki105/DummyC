package lexer

import (
	"../token"
	"io/ioutil"
	"testing"
)

func TestNew(t *testing.T) {
	input, err := readFile("../../test/test.c")
	if err != nil {
		t.Fatalf(err.Error())
	}

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INTTYPE, "int"},
		{token.IDENT, "test"},
		{token.LPAREN, "("},
		{token.INTTYPE, "int"},
		{token.IDENT, "j"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.INTTYPE, "int"},
		{token.IDENT, "i"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.ASSIGN, "="},
		{token.IDENT, "j"},
		{token.ASTERISK, "*"},
		{token.DIGIT, "10"},
		{token.SEMICOLON, ";"},
		{token.RETURN, "return"},
		{token.IDENT, "i"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INTTYPE, "int"},
		{token.IDENT, "main"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.INTTYPE, "int"},
		{token.IDENT, "i"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.ASSIGN, "="},
		{token.DIGIT, "10"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "test"},
		{token.LPAREN, "("},
		{token.IDENT, "i"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.RETURN, "return"},
		{token.DIGIT, "0"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := New(input)

	index := 0
	for {
		token := l.GetToken()
		if token.Type != tests[index].expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", index, tests[index].expectedType, token.Type)
		}
		if token.Literal != tests[index].expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", index, tests[index].expectedLiteral, token.Literal)
		}
		index++
		if !l.GetNextToken() {
			break
		}
	}
}

func readFile(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	lines := string(b)
	return lines, nil
}
