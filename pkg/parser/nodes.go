package parser

import "fmt"

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
