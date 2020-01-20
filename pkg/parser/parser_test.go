package parser

import (
	"dummyc/pkg/ast"
	"dummyc/pkg/lexer"
	"testing"
)

func TestDeclarationStatement(t *testing.T) {
	input := `int main() {
		int i;
	}
	`

	l := lexer.New(input)
	p := New(l)
	translationUnit := p.Parse()
	function := translationUnit.Functions[0]
	stmt := function.Body.Declarations[0]
	if stmt.Name.String() != "i" {
		t.Fatalf("%s != i", stmt.Name.String())
	}
}

func TestReturnStatement(t *testing.T) {
	input := `int main() {
		return 0;
	}
	`

	l := lexer.New(input)
	p := New(l)
	translationUnit := p.Parse()

	mainFunction := translationUnit.Functions[0]

	stmt := mainFunction.Body.Statements[0]
	if _, ok := stmt.(*ast.ReturnStatement); !ok {
		t.Fatal("stmt is not returnStatement")
	}
}
