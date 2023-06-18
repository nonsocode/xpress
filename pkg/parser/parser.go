package parser

import (
	"fmt"
	"strconv"
)

type (
	Expr interface {
		Accept(v Visitor) (interface{}, error)
	}

	Binary struct {
		left     Expr
		operator Token
		right    Expr
	}

	Grouping struct {
		expression Expr
	}

	Literal struct {
		value interface{}
		raw   string
	}

	Unary struct {
		operator Token
		right    Expr
	}

	Template struct {
		expressions []Expr
	}

	Parser struct {
		tokens  []Token
		current int
	}

	Ternary struct {
		condition Expr
		trueExpr  Expr
		falseExpr Expr
	}

	Get struct {
		object Expr
		name   Token
	}

	Call struct {
		callee    Expr
		paren     Token
		arguments []Expr
	}

	Index struct {
		object  Expr
		bracket Token
		index   Expr
	}

	Array struct {
		values []Expr
	}

	Variable struct {
		name Token
	}

	ParseError struct {
		token   Token
		message string
	}
)

func NewBinary(left Expr, operator Token, right Expr) *Binary {
	return &Binary{left: left, operator: operator, right: right}
}

func (b *Binary) Accept(v Visitor) (interface{}, error) {
	return v.visitBinaryExpr(b)
}

func NewGrouping(expression Expr) *Grouping {
	return &Grouping{expression: expression}
}

func (g *Grouping) Accept(v Visitor) (interface{}, error) {
	return v.visitGroupingExpr(g)
}

func NewLiteral(value interface{}, raw string) *Literal {
	return &Literal{value: value, raw: raw}
}

func (l *Literal) Accept(v Visitor) (interface{}, error) {
	return v.visitLiteralExpr(l)
}

func NewUnary(operator Token, right Expr) *Unary {
	return &Unary{operator: operator, right: right}
}

func (u *Unary) Accept(v Visitor) (interface{}, error) {
	return v.visitUnaryExpr(u)
}

func NewTemplate(expressions []Expr) *Template {
	return &Template{expressions: expressions}
}

func (t *Template) Accept(v Visitor) (interface{}, error) {
	return v.visitTemplateExpr(t)
}

func NewTernary(condition Expr, trueExpr Expr, falseExpr Expr) *Ternary {
	return &Ternary{condition: condition, trueExpr: trueExpr, falseExpr: falseExpr}
}

func (t *Ternary) Accept(v Visitor) (interface{}, error) {
	return v.visitTernaryExpr(t)
}

func NewGet(object Expr, name Token) *Get {
	return &Get{object: object, name: name}
}

func (g *Get) Accept(v Visitor) (interface{}, error) {
	return v.visitGetExpr(g)
}

func NewCall(callee Expr, paren Token, arguments []Expr) *Call {
	return &Call{callee: callee, paren: paren, arguments: arguments}
}

func (c *Call) Accept(v Visitor) (interface{}, error) {
	return v.visitCallExpr(c)
}
func NewIndex(object Expr, bracket Token, index Expr) *Index {
	return &Index{object: object, bracket: bracket, index: index}
}

func (i *Index) Accept(v Visitor) (interface{}, error) {
	return v.visitIndexExpr(i)
}
func NewArray(values []Expr) *Array {
	return &Array{values: values}
}

func (a *Array) Accept(visitor Visitor) (interface{}, error) {
	return visitor.visitArrayExpr(a)
}

func NewVariable(name Token) *Variable {
	return &Variable{name: name}
}

func (v *Variable) Accept(visitor Visitor) (interface{}, error) {
	return visitor.visitVariableExpr(v)
}

func (pe *ParseError) Error() string {
	return fmt.Sprintf("Error at '%s': %s", pe.token.lexeme, pe.message)
}

func (p *ParseError) Accept(v Visitor) (interface{}, error) {
	return v.visitParseErrorExpr(p)
}

func NewParser(source string) *Parser {
	lexer := NewLexer(source)
	tokens := lexer.scanTokens()
	return &Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse() Expr {
	return p.template()
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

func (p *Parser) valueTemplate() Expr {
	expr := p.expression()
	_, ok := p.consume(TEMPLATE_RIGHT_BRACE)
	if !ok {
		return p.error(fmt.Sprintf("Expect '}}' after expression. got %v", p.peek().lexeme))
	}
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
	expr := p.logicalOr()
	if p.match(QMARK) {
		trueExpr := p.expression()
		_, ok := p.consume(COLON)
		if !ok {
			return p.error(fmt.Sprintf("Expect ':' after true expression. got %v", p.peek().lexeme))
		}
		falseExpr := p.expression()
		expr = NewTernary(expr, trueExpr, falseExpr)
	}
	return expr
}

func (p *Parser) logicalOr() Expr {
	expr := p.logicalAnd()
	for p.match(OR) {
		operator := p.previous()
		right := p.logicalAnd()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) logicalAnd() Expr {
	expr := p.equality()
	for p.match(AND) {
		operator := p.previous()
		right := p.equality()
		expr = NewBinary(expr, operator, right)
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

func (p *Parser) call() Expr {
	expr := p.primary()
	for {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else if p.match(DOT) {
			name, ok := p.consume(IDENTIFIER)
			if !ok {
				return p.error(fmt.Sprintf("Expect property name after '.' at column %d. got %v", p.current, p.peek().lexeme))
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
	bracket, ok := p.consume(RIGHT_BRACKET)
	if !ok {
		return p.error(fmt.Sprintf("Expect ']' after index expression. got %v", p.peek().lexeme))
	}
	return NewIndex(expr, bracket, index)
}

func (p *Parser) finishCall(expr Expr) Expr {
	args := make([]Expr, 0)
	if !p.check(RIGHT_PAREN) {
		args = append(args, p.expression())
		for p.match(COMMA) {
			args = append(args, p.expression())
		}
	}
	paren, ok := p.consume(RIGHT_PAREN)
	if !ok {
		return p.error(fmt.Sprintf("Expect ')' after arguments. got %v", p.peek().lexeme))
	}
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
			return p.error(fmt.Sprintf("Expect ')' after expression. got %v", p.peek().lexeme))
		}
		return NewGrouping(expr)
	}
	if p.match(LEFT_BRACKET) {
		return p.array()
	}

	return p.error(fmt.Sprintf("Expect expression. got %v", p.peek().lexeme))
}

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
		return p.error(fmt.Sprintf("Expect ']' after array expression. got %v", p.peek().lexeme))
	}

	return NewArray(values)
}

func (p *Parser) consume(tokenType TokenType) (Token, bool) {
	if p.check(tokenType) {
		return p.advance(), true
	}
	return Token{}, false
}

func (p *Parser) error(errorMessage string) *ParseError {
	return &ParseError{message: errorMessage}
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
