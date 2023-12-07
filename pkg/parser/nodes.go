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
		object  Expr
		getType Token
		name    Token
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

	Map struct {
		entries []*MapEntry
	}

	MapEntry struct {
		key   Expr
		value Expr
	}
)

func NewBinary(left Expr, operator Token, right Expr) *Binary {
	return &Binary{left: left, operator: operator, right: right}
}

func (b *Binary) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitBinaryExpr(ctx, b)
}

func NewGrouping(expression Expr) *Grouping {
	return &Grouping{expression: expression}
}

func (g *Grouping) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitGroupingExpr(ctx, g)
}

func NewLiteral(value interface{}, raw string) *Literal {
	return &Literal{value: value, raw: raw}
}

func (l *Literal) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitLiteralExpr(ctx, l)
}

func NewUnary(operator Token, right Expr) *Unary {
	return &Unary{operator: operator, right: right}
}

func (u *Unary) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitUnaryExpr(ctx, u)
}

func NewTemplate(expressions []Expr) *Template {
	return &Template{expressions: expressions}
}

func (t *Template) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitTemplateExpr(ctx, t)
}

func NewTernary(condition Expr, trueExpr Expr, falseExpr Expr) *Ternary {
	return &Ternary{condition: condition, trueExpr: trueExpr, falseExpr: falseExpr}
}

func (t *Ternary) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitTernaryExpr(ctx, t)
}

func NewGet(object Expr, getType, name Token) *Get {
	return &Get{object: object, getType: getType, name: name}
}

func (g *Get) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitGetExpr(ctx, g)
}

func NewCall(callee Expr, arguments []Expr) *Call {
	return &Call{callee: callee, arguments: arguments}
}

func (c *Call) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitCallExpr(ctx, c)
}

func NewIndex(object Expr, index Expr) *Index {
	return &Index{object: object, index: index}
}

func (i *Index) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitIndexExpr(ctx, i)
}

func NewArray(values []Expr) *Array {
	return &Array{values: values}
}

func (a *Array) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitArrayExpr(ctx, a)
}

func NewVariable(name Token) *Variable {
	return &Variable{name: name}
}

func (v *Variable) Accept(ctx context.Context, vis Visitor) (interface{}, error) {
	return vis.visitVariableExpr(ctx, v)
}

func (pe *ParseError) Error() string {
	return fmt.Sprintf("Error at position %d. %s", pe.token.start, pe.message)
}

func (p *ParseError) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitParseErrorExpr(ctx, p)
}

func NewMap(entries []*MapEntry) *Map {
	return &Map{entries: entries}
}

func (m *Map) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitMapExpr(ctx, m)
}

func NewMapEntry(key Expr, value Expr) *MapEntry {
	return &MapEntry{key: key, value: value}
}

func (me *MapEntry) Accept(ctx context.Context, v Visitor) (interface{}, error) {
	return v.visitMapEntryExpr(ctx, me)
}
