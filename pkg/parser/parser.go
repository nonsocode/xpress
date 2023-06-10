package parser

import (
	"fmt"
	"strconv"
)

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

// Visitor interface

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
	return i.interpret(p.expression())
}

func (p *Parser) expression() Expr {
	return p.equality()
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
	return p.primary()
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

func (p *Parser) consume(tokenType TokenType, errorMessage string) {
	if p.check(tokenType) {
		p.advance()
		return
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
