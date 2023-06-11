package parser

type (
	Visitor interface {
		visitBinaryExpr(expr *Binary) (interface{}, error)
		visitGroupingExpr(expr *Grouping) (interface{}, error)
		visitLiteralExpr(expr *Literal) (interface{}, error)
		visitUnaryExpr(expr *Unary) (interface{}, error)

		visitTemplateExpr(expr *Template) (interface{}, error)
		visitTernaryExpr(expr *Ternary) (interface{}, error)
		visitGetExpr(expr *Get) (interface{}, error)
		visitCallExpr(expr *Call) (interface{}, error)
	}
	Interpreter interface {
		interpret(expr Expr) (interface{}, error)
	}
)
