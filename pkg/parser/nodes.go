package parser

import (
	"context"
	"fmt"
)

type (
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
		arguments []Expr
	}

	Index struct {
		object Expr
		index  Expr
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

func (b *Binary) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitBinaryExpr(ctx, b)
}

func (b *Binary) Left() Expr {
	return b.left
}

func (b *Binary) Operator() Token {
	return b.operator
}

func (b *Binary) Right() Expr {
	return b.right
}

func NewGrouping(expression Expr) *Grouping {
	return &Grouping{expression: expression}
}

func (g *Grouping) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitGroupingExpr(ctx, g)
}

func (g *Grouping) Expression() Expr {
	return g.expression
}

func NewLiteral(value interface{}, raw string) *Literal {
	return &Literal{value: value, raw: raw}
}

func (l *Literal) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitLiteralExpr(ctx, l)
}

func (l *Literal) Value() interface{} {
	return l.value
}

func (l *Literal) Raw() string {
	return l.raw
}

func NewUnary(operator Token, right Expr) *Unary {
	return &Unary{operator: operator, right: right}
}

func (u *Unary) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitUnaryExpr(ctx, u)
}

func (u *Unary) Operator() Token {
	return u.operator
}

func (u *Unary) Right() Expr {
	return u.right
}

func NewTemplate(expressions []Expr) *Template {
	return &Template{expressions: expressions}
}

func (t *Template) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitTemplateExpr(ctx, t)
}

func (t *Template) Expressions() []Expr {
	return t.expressions
}

func NewTernary(condition Expr, trueExpr Expr, falseExpr Expr) *Ternary {
	return &Ternary{condition: condition, trueExpr: trueExpr, falseExpr: falseExpr}
}

func (t *Ternary) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitTernaryExpr(ctx, t)
}

func (t *Ternary) Condition() Expr {
	return t.condition
}

func (t *Ternary) TrueExpr() Expr {
	return t.trueExpr
}

func (t *Ternary) FalseExpr() Expr {
	return t.falseExpr
}

func NewGet(object Expr, name Token) *Get {
	return &Get{object: object, name: name}
}

func (g *Get) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitGetExpr(ctx, g)
}

func (g *Get) Object() Expr {
	return g.object
}

func (g *Get) Name() Token {
	return g.name
}

func NewCall(callee Expr, arguments []Expr) *Call {
	return &Call{callee: callee, arguments: arguments}
}

func (c *Call) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitCallExpr(ctx, c)
}

func (c *Call) Callee() Expr {
	return c.callee
}

func (c *Call) Arguments() []Expr {
	return c.arguments
}

func NewIndex(object Expr, index Expr) *Index {
	return &Index{object: object, index: index}
}

func (i *Index) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitIndexExpr(ctx, i)
}

func (i *Index) Object() Expr {
	return i.object
}

func (i *Index) Index() Expr {
	return i.index
}

func NewArray(values []Expr) *Array {
	return &Array{values: values}
}

func (a *Array) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitArrayExpr(ctx, a)
}

func (a *Array) Values() []Expr {
	return a.values
}

func NewVariable(name Token) *Variable {
	return &Variable{name: name}
}

func (v *Variable) Accept(ctx context.Context, vis Visitor) (interface{}, error) {
	return vis.visitVariableExpr(ctx, v)
}

func (v *Variable) Name() Token {
	return v.name
}

func (pe *ParseError) Error() string {
	return fmt.Sprintf("Error at position %d. %s", pe.token.start, pe.message)
}

func (p *ParseError) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitParseErrorExpr(ctx, p)
}
