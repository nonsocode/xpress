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
		visitIndexExpr(expr *Index) (interface{}, error)
		visitVariableExpr(expr *Variable) (interface{}, error)
		visitCallExpr(expr *Call) (interface{}, error)
		visitArrayExpr(expr *Array) (interface{}, error)

		visitParseErrorExpr(err *ParseError) (interface{}, error)
	}
	Interpreter interface {
		Evaluate(expr Expr) (interface{}, error)
	}
)
