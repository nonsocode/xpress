package parser

import "fmt"

// ASTPrinter implements the Visitor interface
type ASTPrinter struct{}

func (ap *ASTPrinter) visitBinaryExpr(expr *Binary) (interface{}, error) {
	return ap.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (ap *ASTPrinter) visitGroupingExpr(expr *Grouping) (interface{}, error) {
	return ap.parenthesize("group", expr.expression)
}

func (ap *ASTPrinter) visitLiteralExpr(expr *Literal) (interface{}, error) {
	if expr.value == nil {
		return "nil", nil
	}
	return expr.value, nil
}

func (ap *ASTPrinter) visitUnaryExpr(expr *Unary) (interface{}, error) {
	return ap.parenthesize(expr.operator.lexeme, expr.right)
}

func (ap *ASTPrinter) visitTemplateExpr(expr *Template) (interface{}, error) {
	return ap.parenthesize("template", expr.expressions...)
}

func (ap *ASTPrinter) visitTernaryExpr(expr *Ternary) (interface{}, error) {
	return ap.parenthesize("ternary", expr.condition, expr.trueExpr, expr.falseExpr)
}

func (ap *ASTPrinter) visitGetExpr(expr *Get) (interface{}, error) {
	return ap.parenthesize("get", expr.object, NewLiteral("tail", "tail"))
}

func (ap *ASTPrinter) visitIndexExpr(expr *Index) (interface{}, error) {
	return ap.parenthesize("index", expr.object, expr.index)
}

func (ap *ASTPrinter) visitVariableExpr(expr *Variable) (interface{}, error) {
	return ap.parenthesize("var", expr)
}

func (ap *ASTPrinter) visitCallExpr(expr *Call) (interface{}, error) {
	args := make([]Expr, len(expr.arguments)+1)
	args[0] = expr.callee
	for i, arg := range expr.arguments {
		args[i+1] = arg
	}
	return ap.parenthesize("call", args...)
}

func (ap *ASTPrinter) visitArrayExpr(expr *Array) (interface{}, error) {
	return ap.parenthesize("array", expr.values...)
}

func (ap *ASTPrinter) visitParseErrorExpr(err *ParseError) (interface{}, error) {
	return nil, err
}

func (ap *ASTPrinter) parenthesize(name string, exprs ...Expr) (string, error) {
	str := "(" + name
	for _, expr := range exprs {
		str += " "
		val, err := expr.accept(ap)
		if err != nil {
			return "", err
		}
		str += fmt.Sprintf("%v", val)
	}
	str += ")"
	return str, nil
}

func (ap *ASTPrinter) Print(expr Expr) (interface{}, error) {
	return expr.accept(&ASTPrinter{})
}
