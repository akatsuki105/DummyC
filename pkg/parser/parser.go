package parser

import (
	"dummyc/pkg/ast"
	"dummyc/pkg/lexer"
	"dummyc/pkg/token"
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
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

type (
	infixParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	infixParseFns map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
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

Loop:
	for {
		switch p.l.GetCurType() {
		case token.INTTYPE:
			// プロトタイプ宣言
			prototype := p.parsePrototype()
			// 次が;ならプロトタイプ宣言 {なら関数定義
			switch p.l.GetCurType() {
			case token.SEMICOLON:
				program.Prototypes = append(program.Prototypes, *prototype)
				p.l.GetNextToken()
			case token.LBRACE:
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
	prototype := &ast.Prototype{Token: p.l.GetToken()}

	p.l.GetNextToken() // int => identifier

	if p.l.GetCurType() != token.IDENT {
		panic("panic")
	}
	prototype.Name = p.parseIdentifier()

	if p.l.GetCurType() != token.LPAREN {
		panic("panic")
	}
	p.l.GetNextToken() // ( => parameter

	for {
		if p.l.GetCurType() == token.RPAREN {
			p.l.GetNextToken()
			break
		} else if p.l.GetCurType() != token.INTTYPE {
			panic("panic")
		}

		p.l.GetNextToken()
		prototype.Parameters = append(prototype.Parameters, p.parseIdentifier())

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
	functionLiteral.Body = *p.parseFunctionStatement()

	if p.l.GetCurType() == token.RBRACE {
		p.l.GetNextToken() // } => 次の関数
	}

	return functionLiteral
}

func (p *Parser) parseFunctionStatement() *ast.FunctionStatement {
	functionStmt := &ast.FunctionStatement{}

	p.l.GetNextToken() // { => ...

	// parse DeclarationStatements
	for p.l.GetCurType() == token.INTTYPE {
		functionStmt.Declarations = append(functionStmt.Declarations, *p.parseDeclarationStatement())
	}

	// parse Statements
	for p.l.GetCurType() != token.RBRACE {
		functionStmt.Statements = append(functionStmt.Statements, p.parseStatement())
		p.l.GetNextToken()
	}

	return functionStmt
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	identifier := &ast.Identifier{
		Token: p.l.GetToken(),
		Value: p.l.GetCurString(),
	}
	p.l.GetNextToken()
	return identifier
}

func (p *Parser) parseDeclarationStatement() *ast.DeclarationStatement {
	declarationStatement := &ast.DeclarationStatement{
		Token: p.l.GetToken(),
	}
	p.l.GetNextToken()

	declarationStatement.Name = *p.parseIdentifier()

	for p.l.GetCurType() != token.SEMICOLON {
		p.l.GetNextToken()
	}
	p.l.GetNextToken() // ; => 次
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

	p.l.GetNextToken()
	call.Arguments = p.parseExpressionList(token.RPAREN)

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
