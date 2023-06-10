package parser

type (
	Visitor interface {
		visitBinaryExpr(expr *Binary) interface{}
		visitGroupingExpr(expr *Grouping) interface{}
		visitLiteralExpr(expr *Literal) interface{}
		visitUnaryExpr(expr *Unary) interface{}
	}
	Interpreter interface {
		interpret(expr Expr) interface{}
	}
)
