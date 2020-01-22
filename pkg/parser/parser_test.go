package parser

import (
	"../ast"
	"../lexer"
	"testing"
)

func TestDeclarationStatement(t *testing.T) {
	input := `int main() {
		int i;
		return 0;
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

func TestExpression(t *testing.T) {
	input := `int test(int j) {
		int i;
		i = j*10;
		return 0;
	}`

	l := lexer.New(input)
	p := New(l)
	translationUnit := p.Parse()

	mainFunction := translationUnit.Functions[0]

	stmt := mainFunction.Body.Statements[0]
	if _, ok := stmt.(*ast.ExpressionStatement); !ok {
		t.Fatal("stmt is not ExpressionStatement")
	}
}

func TestCallStatement(t *testing.T) {
	input := `int test(int j);

	int main() {
		test(10);
		return 0;
	}`

	l := lexer.New(input)
	p := New(l)
	translationUnit := p.Parse()

	mainFunction := translationUnit.Functions[0]
	statements := mainFunction.Body.Statements

	if len(statements) != 2 {
		t.Fatalf("statements does not contain %d statements. got=%d\n", 1, len(statements))
	}

	stmt, ok := statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			statements[0])
	}

	_, ok = stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}
}
