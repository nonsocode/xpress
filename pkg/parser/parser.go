package parser

import (
	"context"
	"fmt"
	"strconv"
)

type (
	Expr interface {
		Accept(context.Context, Visitor) (interface{}, error)
	}
)

func NewParser(source string) *Parser {
	lexer := NewLexer(source)
	tokens := lexer.scanTokens()
	return &Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse() Expr {
	return p.template()
}

// Grammar:
// template  → ( valueTemplate | TEXT )* ;
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

// Grammar:
// valueTemplate → TEMPLATE_START expression TEMPLATE_END ;
func (p *Parser) valueTemplate() Expr {
	expr := p.expression()
	_, ok := p.consume(TEMPLATE_RIGHT_BRACE)
	if !ok {
		return p.error(fmt.Sprintf("Expect '}}' after expression. got %v", p.peek().lexeme), p.peek())
	}
	return expr
}

// Grammar:
// TEXT → [^\{\}]+ ;
func (p *Parser) text() Expr {
	token := p.advance()
	return NewLiteral(token.lexeme, token.lexeme)
}

// Grammar:
// expression  → ternary ;
func (p *Parser) expression() Expr {
	return p.ternary()
}

// Grammar:
// ternary → logicalOr ( QMARK expression COLON expression )? ;
func (p *Parser) ternary() Expr {
	expr := p.logicalOr()
	if p.match(QMARK) {
		trueExpr := p.expression()
		_, ok := p.consume(COLON)
		if !ok {
			return p.error(fmt.Sprintf("Expect ':' after true expression. got %v", p.peek().lexeme), p.peek())
		}
		falseExpr := p.expression()
		expr = NewTernary(expr, trueExpr, falseExpr)
	}
	return expr
}

// Grammar:
// logicalOr  → logicalAnd ( OR logicalAnd )* ;
func (p *Parser) logicalOr() Expr {
	expr := p.logicalAnd()
	for p.match(OR) {
		operator := p.previous()
		right := p.logicalAnd()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

// Grammar:
// logicalAnd  → equality ( AND equality )* ;
func (p *Parser) logicalAnd() Expr {
	expr := p.equality()
	for p.match(AND) {
		operator := p.previous()
		right := p.equality()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

// Grammar:
// equality  → comparison ( ( BANG_EQUAL | EQUAL_EQUAL ) comparison )* ;
func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

// Grammar:
// comparison  → term ( ( GREATER | GREATER_EQUAL | LESS | LESS_EQUAL ) term )* ;
func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL, BANG_EQUAL, EQUAL_EQUAL) {
		expr = NewBinary(expr, p.previous(), p.term())
	}
	return expr
}

// Grammar:
// term  → factor ( ( MINUS | PLUS ) factor )* ;
func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(MINUS, PLUS) {
		expr = NewBinary(expr, p.previous(), p.factor())
	}
	return expr
}

// Grammar:
// factor  → unary ( ( SLASH | STAR ) unary )* ;
func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(SLASH, STAR) {
		expr = NewBinary(expr, p.previous(), p.unary())
	}
	return expr
}

// Grammar:
// unary  → ( BANG | MINUS ) unary | call ;
func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		return NewUnary(p.previous(), p.unary())
	}
	return p.call()
}

// Grammar:
// call → primary ( LPAREN arguments? RPAREN )* ( get )* ;
func (p *Parser) call() Expr {
	expr := p.primary()
	for {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else if p.match(DOT) {
			name, ok := p.consume(IDENTIFIER)
			if !ok {
				return p.error(fmt.Sprintf("Expect property name after '.'. got %v", p.peek().lexeme), p.peek())
			}
			expr = NewGet(expr, name)
		} else if p.match(LEFT_BRACKET) {
			expr = p.finishIndex(expr)
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) finishIndex(expr Expr) Expr {
	index := p.expression()
	_, ok := p.consume(RIGHT_BRACKET)
	if !ok {
		return p.error(fmt.Sprintf("Expect ']' after index expression. got %v", p.peek().lexeme), p.peek())
	}
	return NewIndex(expr, index)
}

func (p *Parser) finishCall(expr Expr) Expr {
	args := make([]Expr, 0)
	if !p.check(RIGHT_PAREN) {
		args = append(args, p.expression())
		for p.match(COMMA) {
			args = append(args, p.expression())
		}
	}
	_, ok := p.consume(RIGHT_PAREN)
	if !ok {
		return p.error(fmt.Sprintf("Expect ')' after arguments. got %v", p.peek().lexeme), p.peek())
	}
	return NewCall(expr, args)
}

// Grammar:
// primary  → "true" | "false" | "nil" | NUMBER | STRING | IDENTIFIER | LPAREN expression RPAREN | array ;
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
	if p.match(NUMBER) {
		num, _ := strconv.ParseFloat(p.previous().lexeme, 64)

		return NewLiteral(num, p.previous().lexeme)
	}
	if p.match(IDENTIFIER) {
		return NewVariable(p.previous())
	}
	if p.match(STRING) {
		str := p.previous().lexeme
		return NewLiteral(str[1:len(str)-1], p.previous().lexeme)
	}
	if p.match(LEFT_PAREN) {
		expr := p.expression()
		_, ok := p.consume(RIGHT_PAREN)
		if !ok {
			return p.error(fmt.Sprintf("Expect ')' after expression. got %v", p.peek().lexeme), p.peek())
		}
		return NewGrouping(expr)
	}
	if p.match(LEFT_BRACKET) {
		return p.array()
	}

	return p.error(fmt.Sprintf("Expect expression. got %v", p.peek().lexeme), p.peek())
}

// Grammar:
// array  → LBRACKET ( expression ( COMMA expression )* )? RBRACKET ;
func (p *Parser) array() Expr {
	values := make([]Expr, 0)
	if !p.check(RIGHT_BRACKET) {
		values = append(values, p.expression())
		for p.match(COMMA) {
			values = append(values, p.expression())
		}
	}
	_, ok := p.consume(RIGHT_BRACKET)
	if !ok {
		return p.error(fmt.Sprintf("Expect ']' after array expression. got %v", p.peek().lexeme), p.peek())
	}

	return NewArray(values)
}

func (p *Parser) consume(tokenType TokenType) (Token, bool) {
	if p.check(tokenType) {
		return p.advance(), true
	}
	return Token{}, false
}

func (p *Parser) error(errorMessage string, token Token) *ParseError {
	return &ParseError{message: errorMessage, token: token}
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
