package parser

type (
	Visitor interface {
		visitBinaryExpr(expr *Binary) interface{}
		visitGroupingExpr(expr *Grouping) interface{}
		visitLiteralExpr(expr *Literal) interface{}
		visitUnaryExpr(expr *Unary) interface{}

		visitTemplateExpr(expr *Template) interface{}
		visitTernaryExpr(expr *Ternary) interface{}
		visitGetExpr(expr *Get) interface{}
		visitCallExpr(expr *Call) interface{}
	}
	Interpreter interface {
		interpret(expr Expr) interface{}
	}
)
