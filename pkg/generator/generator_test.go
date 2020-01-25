package generator

import (
	"../lexer"
	"../parser"
	"io/ioutil"
	"testing"
)

func TestGenerator(t *testing.T) {
	input, err := readFile("../../test/test.dc")
	if err != nil {
		t.Fatalf(err.Error())
	}

	l := lexer.New(input)
	p := parser.New(l)
	tu := p.Parse()
	g := New()
	g.Generate(tu, input)
}

func readFile(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	lines := string(b)
	return lines, nil
}
