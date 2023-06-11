package parser

import (
	"fmt"
	"strconv"
)

// Grammar:
// template          → ( valueTemplate | TEXT )* ;
// valueTemplate     → TEMPLATE_START expression TEMPLATE_END ;
// expression        → ternary ;
// ternary           → equality ( QMARK expression COLON expression )? ;
// equality          → comparison ( ( BANG_EQUAL | EQUAL_EQUAL ) comparison )* ;
// comparison        → term ( ( GREATER | GREATER_EQUAL | LESS | LESS_EQUAL ) term )* ;
// term              → factor ( ( MINUS | PLUS ) factor )* ;
// factor            → unary ( ( SLASH | STAR ) unary )* ;
// unary             → ( BANG | MINUS ) unary | call ;
// call              → primary ( LPAREN arguments? RPAREN | DOT identifier )* ;
// primary           → number | string | TRUE | FALSE | NIL | identifier | LPAREN expression RPAREN
// arguments         → expression ( COMMA expression )* ;
// identifier        → LETTER ( LETTER | DIGIT )* ;
// number            → DIGIT+ ( DOT DIGIT+ )? ;
// string            → (DQUOTE characters? DQUOTE) ;
// characters        → ( escape | char )* ;
// escape            → "\\" char ;
// LETTER            → [a-zA-Z] ;
// DIGIT             → [0-9] ;
// TEMPLATE_START    → "{{" ;
// TEMPLATE_END      → "}}" ;
// QMARK             → "?" ;
// COLON             → ":" ;
// COMMA             → "," ;
// DOT               → "." ;
// EQUAL             → "=" ;
// BANG              → "!" ;
// MINUS             → "-" ;
// PLUS              → "+" ;
// DQUOTE            → "\"" ;
// SQUOTE            → "\'" ;
// SLASH             → "/" ;
// STAR              → "*" ;
// LPAREN            → "(" ;
// RPAREN            → ")" ;
// EQUAL_EQUAL       → "==" ;
// BANG_EQUAL        → "!=" ;
// GREATER           → ">" ;
// GREATER_EQUAL     → ">=" ;
// LESS              → "<" ;
// LESS_EQUAL        → "<=" ;
// TRUE              → "true" ;
// FALSE             → "false" ;
// NIL               → "nil" ;
// TEXT              → [^\{\}]+ ;
// char              → [^\"] ;

type Expr interface {
	// Expr is an interface that all expressions implement.
	// It has an accept method that takes a visitor interface.
	// The visitor interface has a visit method for each of the expression classes.
	// The accept method calls the visit method for the expression’s class.
	// The visit method then calls the accept method on the expression’s children.
	// This is the essence of the Visitor pattern.
	accept(v Visitor) interface{}
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func NewBinary(left Expr, operator Token, right Expr) *Binary {
	return &Binary{left: left, operator: operator, right: right}
}

func (b *Binary) accept(v Visitor) interface{} {
	return v.visitBinaryExpr(b)
}

type Grouping struct {
	expression Expr
}

func NewGrouping(expression Expr) *Grouping {
	return &Grouping{expression: expression}
}

func (g *Grouping) accept(v Visitor) interface{} {
	return v.visitGroupingExpr(g)
}

type Literal struct {
	value interface{}
	raw   string
}

func NewLiteral(value interface{}, raw string) *Literal {
	return &Literal{value: value, raw: raw}
}

func (l *Literal) accept(v Visitor) interface{} {
	return v.visitLiteralExpr(l)
}

type Unary struct {
	operator Token
	right    Expr
}

func NewUnary(operator Token, right Expr) *Unary {
	return &Unary{operator: operator, right: right}
}

func (u *Unary) accept(v Visitor) interface{} {
	return v.visitUnaryExpr(u)
}

type Template struct {
	expressions []Expr
}

func NewTemplate(expressions []Expr) *Template {
	return &Template{expressions: expressions}
}

func (t *Template) accept(v Visitor) interface{} {
	return v.visitTemplateExpr(t)
}

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(source string) *Parser {
	lexer := NewLexer(source)
	tokens := lexer.scanTokens()
	return &Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse(i Interpreter) interface{} {
	return i.interpret(p.template())
}

func (p *Parser) template() Expr {
	var exprs []Expr
	for !p.isAtEnd() {
		if p.match(TEMPLATE_LEFT_BRACE) {
			exprs = append(exprs, p.valueTemplate())
		} else {
			exprs = append(exprs, p.text())
		}
	}
	return NewTemplate(exprs)
}

type Ternary struct {
	condition Expr
	trueExpr  Expr
	falseExpr Expr
}

func NewTernary(condition Expr, trueExpr Expr, falseExpr Expr) *Ternary {
	return &Ternary{condition: condition, trueExpr: trueExpr, falseExpr: falseExpr}
}

func (t *Ternary) accept(v Visitor) interface{} {
	return v.visitTernaryExpr(t)
}

func (p *Parser) valueTemplate() Expr {
	expr := p.expression()
	p.consume(TEMPLATE_RIGHT_BRACE, "Expect '}}' after expression.")
	return expr
}

func (p *Parser) text() Expr {
	token := p.advance()
	return NewLiteral(token.lexeme, token.lexeme)
}

func (p *Parser) expression() Expr {
	return p.ternary()
}

func (p *Parser) ternary() Expr {
	expr := p.equality()
	if p.match(QMARK) {
		trueExpr := p.expression()
		p.consume(COLON, "Expect ':' after true expression.")
		falseExpr := p.expression()
		expr = NewTernary(expr, trueExpr, falseExpr)
	}
	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL, BANG_EQUAL, EQUAL_EQUAL) {
		expr = NewBinary(expr, p.previous(), p.term())
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(MINUS, PLUS) {
		expr = NewBinary(expr, p.previous(), p.factor())
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(SLASH, STAR) {
		expr = NewBinary(expr, p.previous(), p.unary())
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		return NewUnary(p.previous(), p.unary())
	}
	return p.call()
}

type Get struct {
	object Expr
	name   Token
}

func NewGet(object Expr, name Token) *Get {
	return &Get{object: object, name: name}
}

func (g *Get) accept(v Visitor) interface{} {
	return v.visitGetExpr(g)
}

type Call struct {
	callee    Expr
	paren     Token
	arguments []Expr
}

func NewCall(callee Expr, paren Token, arguments []Expr) *Call {
	return &Call{callee: callee, paren: paren, arguments: arguments}
}

func (c *Call) accept(v Visitor) interface{} {
	return v.visitCallExpr(c)
}

func (p *Parser) call() Expr {
	expr := p.primary()
	for {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else if p.match(DOT) {
			getToken := p.consume(IDENTIFIER, "Expect property name after '.'.")
			expr = NewGet(expr, getToken)
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) finishCall(expr Expr) Expr {
	args := make([]Expr, 0)
	if !p.check(RIGHT_PAREN) {
		args = append(args, p.expression())
		for p.match(COMMA) {
			args = append(args, p.expression())
		}
	}
	paren := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")
	return NewCall(expr, paren, args)
}

func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return NewLiteral(false, "false")
	}
	if p.match(TRUE) {
		return NewLiteral(true, "true")
	}
	if p.match(NIL) {
		return NewLiteral(nil, "nil")
	}
	if p.match(NUMBER, STRING) {
		num, _ := strconv.ParseFloat(p.previous().lexeme, 64)

		return NewLiteral(num, p.previous().lexeme)
	}
	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return NewGrouping(expr)
	}
	panic(p.error(p.peek(), "Expect expression."))
}

func (p *Parser) consume(tokenType TokenType, errorMessage string) Token {
	if p.check(tokenType) {
		return p.advance()
	}
	panic(p.error(p.peek(), errorMessage).message)
}

func (p *Parser) error(token Token, errorMessage string) ParseError {
	return ParseError{token: token, message: errorMessage}
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().tokenType == tokenType
}

func (p *Parser) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

type ParseError struct {
	token   Token
	message string
}

func (pe *ParseError) Error() string {
	return fmt.Sprintf("Error at '%s': %s", pe.token.lexeme, pe.message)
}
