package parser

import (
	"../ast"
	"../lexer"
	"../token"
	"fmt"
)

const (
	_ int = iota
	LOWEST
	SUM     // +
	PRODUCT // *
	PREFIX  // -X or !X
	CALL    // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.ASSIGN:   SUM,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

type (
	infixParseFn func(ast.Expression) ast.Expression
)

type Function struct {
	Name string
	Argc int
}

type Parser struct {
	l      *lexer.Lexer
	errors []string

	infixParseFns map[token.TokenType]infixParseFn

	variableTable  []string   // 宣言済みの変数名を登録する
	prototypeTable []Function // プロトタイプ宣言済みの関数
	functionTable  []Function // 定義済みの関数
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.ASSIGN, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	return p
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) Parse() *ast.TranslationUnit {
	program := &ast.TranslationUnit{}

	// printnum関数を事前定義
	printnum := ast.Prototype{
		Token: *token.New(token.INTTYPE, "int", 0),
	}
	name := &ast.Identifier{
		Token: *token.New(token.IDENT, "printnum", 0),
		Value: "printnum",
	}
	param := &ast.Identifier{
		Token: *token.New(token.IDENT, "i", 0),
		Value: "i",
	}
	printnum.Name = name
	printnum.Parameters = []*ast.Identifier{param}
	program.Prototypes = append(program.Prototypes, printnum)
	p.prototypeTable = append(p.prototypeTable, Function{"printnum", 1})

Loop:
	for {
		switch p.l.GetCurType() {
		case token.INTTYPE:
			// プロトタイプ宣言
			prototype := p.parsePrototype()

			// 正当性チェックに使うオブジェクトを作成
			name := prototype.Name.Name()
			argc := len(prototype.Parameters)
			fn := Function{name, argc}

			// 次が;ならプロトタイプ宣言 {なら関数定義
			switch p.l.GetCurType() {

			case token.SEMICOLON:
				// プロトタイプ宣言

				// 再定義チェック
				if ok := p.checkReDefinition(fn); !ok {
					panicMsg := name + " is already definitions"
					panic(panicMsg)
				}
				p.prototypeTable = append(p.prototypeTable, fn)

				program.Prototypes = append(program.Prototypes, *prototype)
				p.l.GetNextToken()

			case token.LBRACE:
				// 関数定義

				// プロトタイプ宣言が正当かチェック
				if ok := p.checkCorrectDefinition(fn); !ok {
					panicMsg := name + " is invalid definition"
					panic(panicMsg)
				}

				program.Functions = append(program.Functions, *p.parseFunctionLiteral(prototype))

			default:
				panic("invalid token")

			}
		case token.EOF:
			break Loop
		default:
			panic("not prototype or EOF")
		}
	}

	return program
}

func (p *Parser) parsePrototype() *ast.Prototype {
	paramList := []string{}

	prototype := &ast.Prototype{Token: p.l.GetToken()}

	p.l.GetNextToken() // int => identifier

	if p.l.GetCurType() != token.IDENT {
		panic("panic")
	}
	prototype.Name = p.parseIdentifier()
	p.l.GetNextToken() // identifier => (

	if p.l.GetCurType() != token.LPAREN {
		panic("panic")
	}
	p.l.GetNextToken() // ( => parameter

	// parameter (int a, int b, ...)
	for {
		if p.l.GetCurType() == token.RPAREN {
			p.l.GetNextToken()
			break
		} else if p.l.GetCurType() != token.INTTYPE {
			panic("panic")
		}

		p.l.GetNextToken()
		identifier := p.parseIdentifier()
		if contains(paramList, identifier.Token.Literal) {
			panic("already used")
		}
		prototype.Parameters = append(prototype.Parameters, identifier)
		paramList = append(paramList, identifier.Token.Literal)
		p.l.GetNextToken()

		if p.l.GetCurType() == token.RPAREN {
			p.l.GetNextToken()
			break
		} else if p.l.GetCurType() == token.COMMA {
			p.l.GetNextToken()
			continue
		} else {
			panic("panic")
		}
	}

	return prototype
}

func (p *Parser) parseFunctionLiteral(prototype *ast.Prototype) *ast.FunctionLiteral {
	functionLiteral := &ast.FunctionLiteral{
		Token:     p.l.GetToken(),
		Prototype: *prototype,
	}
	p.variableTable = []string{}
	functionLiteral.Body = *p.parseFunctionStatement(prototype)
	p.functionTable = append(p.functionTable, Function{prototype.GetName(), prototype.GetParamNum()})

	if p.l.GetCurType() == token.RBRACE {
		p.l.GetNextToken() // } => 次の関数
	}

	return functionLiteral
}

func (p *Parser) parseFunctionStatement(prototype *ast.Prototype) *ast.FunctionStatement {
	functionStmt := &ast.FunctionStatement{}

	p.l.GetNextToken() // { => ...

	// parse parameter
	argc := len(prototype.Parameters)
	for i := 0; i < argc; i++ {
		vdecl := &ast.DeclarationStatement{
			Token: p.l.GetToken(),
			Name:  *prototype.Parameters[i],
		}
		vdecl.SetDeclType(ast.Param)
		p.variableTable = append(p.variableTable, vdecl.Name.Name())
		functionStmt.Declarations = append(functionStmt.Declarations, *vdecl)
	}

	// parse DeclarationStatements
	for p.l.GetCurType() == token.INTTYPE {
		stmt := p.parseDeclarationStatement()
		if contains(p.variableTable, stmt.Name.Name()) {
			panic("already declared")
		}
		p.variableTable = append(p.variableTable, stmt.Name.Name())
		functionStmt.Declarations = append(functionStmt.Declarations, *stmt)
		p.l.GetNextToken()
	}

	// parse Statements
	var stmt ast.Statement
	for p.l.GetCurType() != token.RBRACE {
		stmt = p.parseStatement()
		functionStmt.Statements = append(functionStmt.Statements, stmt)
		p.l.GetNextToken()
	}

	// 最後のStatementがReturnか確認
	if stmt.TokenLiteral() != "return" {
		panic("not return")
	}

	return functionStmt
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	identifier := &ast.Identifier{
		Token: p.l.GetToken(),
		Value: p.l.GetCurString(),
	}
	return identifier
}

func (p *Parser) parseDeclarationStatement() *ast.DeclarationStatement {
	declarationStatement := &ast.DeclarationStatement{
		Token: p.l.GetToken(),
	}
	declarationStatement.SetDeclType(ast.Local)
	p.l.GetNextToken() // INTTYPE => identifer

	declarationStatement.Name = *p.parseIdentifier()

	for p.l.GetNextType() == token.SEMICOLON {
		p.l.GetNextToken() // identifer => semicolon
	}
	return declarationStatement
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.l.GetCurType() {
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{
		Token: p.l.GetToken(),
	}

	p.l.GetNextToken() // return => expression

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.l.GetNextType() == token.SEMICOLON {
		p.l.GetNextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token: p.l.GetToken(),
	}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.l.GetNextType() == token.SEMICOLON {
		p.l.GetNextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	var exp ast.Expression
	switch p.l.GetCurType() {
	case token.IDENT:
		exp = p.parseIdentifier()
	case token.DIGIT:
		exp = p.parseNumber()
	}

	for p.l.GetNextType() != token.SEMICOLON && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.l.GetNextType()]
		if infix == nil {
			return exp
		}

		p.l.GetNextToken()
		exp = infix(exp)
	}

	return exp
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.l.GetCurType()]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	pt := p.l.GetNextType()
	if p, ok := precedences[pt]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseNumber() *ast.Number {
	number := &ast.Number{
		Token: p.l.GetToken(),
		Value: p.l.GetCurNumVal(),
	}
	return number
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.l.GetToken(),
		Operator: p.l.GetCurString(),
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.l.GetNextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	call := &ast.CallExpression{
		Token:    p.l.GetToken(),
		Function: function,
	}

	// プロトタイプ宣言されているか確認し引数の数をテーブルから取得
	argc := -1
	for _, prototype := range p.prototypeTable {
		if function.TokenLiteral() == prototype.Name {
			argc = prototype.Argc
			break
		}
	}

	// 関数定義済みであるか確認し、引数の数をテーブルから取得
	if argc == -1 {
		for _, fn := range p.functionTable {
			if function.TokenLiteral() == fn.Name {
				argc = fn.Argc
				break
			}
		}
	}

	if argc == -1 {
		msg := fmt.Sprintf("%s is not defined", function.TokenLiteral())
		panic(msg)
	}

	call.Arguments = p.parseExpressionList(token.RPAREN)

	// 引数の数を確認
	if argc != len(call.Arguments) {
		fmt.Println(argc, len(call.Arguments))
		panic("argc is not correct")
	}

	return call
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.l.GetNextType() == end {
		p.l.GetNextToken()
		return list
	}

	p.l.GetNextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.l.GetNextType() == token.COMMA {
		p.l.GetNextToken() // previous => ,
		p.l.GetNextToken() // , => expression
		list = append(list, p.parseExpression(LOWEST))
	}

	if p.l.GetNextType() != end {
		return nil
	}
	p.l.GetNextToken() // ) => {

	return list
}

func contains(slice []string, target string) bool {
	for _, s := range slice {
		if target == s {
			return true
		}
	}
	return false
}
