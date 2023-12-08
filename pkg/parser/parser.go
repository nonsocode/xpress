package parser

import (
	"context"
	"fmt"
	"strconv"
)

type (
	Expr interface {
		Accept(context.Context, Visitor) EvaluationResult
	}
)

func NewParser(source string) *Parser {
	lexer := NewLexer(source)
	tokens := lexer.scanTokens()
	return &Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse() (exp Expr) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(*ParseError); ok {
				exp = err
			} else {
				exp = &ParseError{
					message: fmt.Sprintf("%v", r),
					token:   p.peek(),
				}
			}
		}
	}()
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

	if ok := p.consume(TEMPLATE_RIGHT_BRACE); !ok {
		p.error(
			fmt.Sprintf("Expect '}}' after expression. got %v", p.peek().lexeme),
			p.peek(),
		)
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
		if ok := p.consume(COLON); !ok {
			p.error(
				fmt.Sprintf("Expect ':' after true expression. got %v", p.peek().lexeme),
				p.peek(),
			)
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
		} else if p.matchAny(DOT, OPTIONALCHAIN) {
			expr = p.get(expr, p.previous())
		} else if p.match(LEFT_BRACKET) {
			expr = p.index(expr)
		} else {
			break
		}
	}
	return expr
}

// Grammar:
// get → (QMARK)? DOT IDENTIFIER ;
func (p *Parser) get(expr Expr, token Token) Expr {
	if ok := p.consume(IDENTIFIER); !ok {
		p.error(fmt.Sprintf("Expect property name after '%s'. got %v", p.previous().lexeme, p.peek().lexeme), p.peek())
	}
	return NewGet(expr, token, p.previous())
}

// Grammar:
// index → LBRACKET expression RBRACKET ;
func (p *Parser) index(expr Expr) Expr {
	index := p.expression()
	if ok := p.consume(RIGHT_BRACKET); !ok {
		p.error(
			fmt.Sprintf("Expect ']' after index expression. got %v", p.peek().lexeme),
			p.peek(),
		)
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
	if ok := p.consume(RIGHT_PAREN); !ok {
		p.error(fmt.Sprintf("Expect ')' after arguments. got %v", p.peek().lexeme), p.peek())
	}
	return NewCall(expr, args)
}

// Grammar:
// primary  → "true" | "false" | "nil" | NUMBER | STRING | IDENTIFIER | LPAREN expression RPAREN | array | map;
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
		if ok := p.consume(RIGHT_PAREN); !ok {
			p.error(
				fmt.Sprintf("Expect ')' after expression. got %v", p.peek().lexeme),
				p.peek(),
			)
		}
		return NewGrouping(expr)
	}
	if p.match(LEFT_BRACKET) {
		return p.array()
	}
	if p.match(LEFT_BRACE) {
		return p.mapExpr()
	}

	p.error(fmt.Sprintf("Expect expression. got %v", p.peek().lexeme), p.peek())
	return nil
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
	if ok := p.consume(RIGHT_BRACKET); !ok {
		p.error(
			fmt.Sprintf("Expect ']' after array expression. got %v", p.peek().lexeme),
			p.peek(),
		)
	}

	return NewArray(values)
}

// Grammar:
// map  → LBRACE ( mapEntry ( COMMA mapEntry )* )? RBRACE ;
func (p *Parser) mapExpr() Expr {
	entries := make([]*MapEntry, 0)
	if !p.check(RIGHT_BRACE) {
		entries = append(entries, p.mapEntry())
		for p.match(COMMA) {
			entries = append(entries, p.mapEntry())
		}
	}
	if ok := p.consume(RIGHT_BRACE); !ok {
		p.error(
			fmt.Sprintf("Expect '}' after map expression. got %v", p.peek().lexeme),
			p.peek(),
		)
	}

	return NewMap(entries)
}

// Grammar:
// mapEntry  → ( identifier | string | LBRACKET expression RBRACKET ) COLON expression ;
func (p *Parser) mapEntry() *MapEntry {
	var key Expr
	if p.match(IDENTIFIER) {
		key = NewLiteral(p.previous().lexeme, p.previous().lexeme)
	}
	if p.match(STRING) {
		str := p.previous().lexeme
		key = NewLiteral(str[1:len(str)-1], str)
	}
	if p.match(LEFT_BRACKET) {
		key = p.expression()
		if ok := p.consume(RIGHT_BRACKET); !ok {
			p.error(
				fmt.Sprintf("Expect ']' after index expression. got %v", p.peek().lexeme),
				p.peek(),
			)
		}
	}
	if key == nil {
		p.error(
			fmt.Sprintf("Expect map key. got %v", p.peek().lexeme),
			p.peek(),
		)
	}
	if ok := p.consume(COLON); !ok {
		p.error(
			fmt.Sprintf("Expect ':' after map key. got %v", p.peek().lexeme),
			p.peek(),
		)
	}
	value := p.expression()
	return NewMapEntry(key, value)
}

func (p *Parser) consume(tokenType TokenType) bool {
	if p.check(tokenType) {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) error(errorMessage string, token Token) {
	panic(&ParseError{message: errorMessage, token: token})
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

func (p *Parser) matchAny(types ...TokenType) bool {
	for _, t := range types {
		if p.match(t) {
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
