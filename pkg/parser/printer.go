package parser

import "fmt"

// ASTPrinter implements the Visitor interface
type ASTPrinter struct{}

func (ap *ASTPrinter) visitBinaryExpr(expr *Binary) interface{} {
	return ap.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (ap *ASTPrinter) visitGroupingExpr(expr *Grouping) interface{} {
	return ap.parenthesize("group", expr.expression)
}

func (ap *ASTPrinter) visitLiteralExpr(expr *Literal) interface{} {
	if expr.value == nil {
		return "nil"
	}
	return expr.value
}

func (ap *ASTPrinter) visitUnaryExpr(expr *Unary) interface{} {
	return ap.parenthesize(expr.operator.lexeme, expr.right)
}

func (ap *ASTPrinter) parenthesize(name string, exprs ...Expr) string {
	str := "(" + name
	for _, expr := range exprs {
		str += " "
		str += fmt.Sprintf("%v", expr.accept(ap))
	}
	str += ")"
	return str
}

func (ap *ASTPrinter) interpret(expr Expr) interface{} {
	return expr.accept(&ASTPrinter{})
}

// RPNPrinter implements the Visitor interface
type RPNPrinter struct{}

func (rp *RPNPrinter) visitBinaryExpr(expr *Binary) interface{} {
	return fmt.Sprintf("%s %s %s", expr.left.accept(rp), expr.right.accept(rp), expr.operator.lexeme)
}

func (rp *RPNPrinter) visitGroupingExpr(expr *Grouping) interface{} {
	return fmt.Sprintf("%s", expr.expression.accept(rp))
}

func (rp *RPNPrinter) visitLiteralExpr(expr *Literal) interface{} {
	if expr.value == nil {
		return "nil"
	}
	return expr.value
}

func (rp *RPNPrinter) visitUnaryExpr(expr *Unary) interface{} {
	return fmt.Sprintf("%s %s", expr.right.accept(rp), expr.operator.lexeme)
}

// PrintRPN prints the AST in RPN
func (rp *RPNPrinter) interpret(expr Expr) interface{} {
	return expr.accept(&RPNPrinter{})
}
