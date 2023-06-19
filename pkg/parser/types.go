package parser

import "context"

type (
	Visitor interface {
		visitBinaryExpr(context.Context, *Binary) (interface{}, error)
		visitGroupingExpr(context.Context, *Grouping) (interface{}, error)
		visitLiteralExpr(context.Context, *Literal) (interface{}, error)
		visitUnaryExpr(context.Context, *Unary) (interface{}, error)

		visitTemplateExpr(context.Context, *Template) (interface{}, error)
		visitTernaryExpr(context.Context, *Ternary) (interface{}, error)
		visitGetExpr(context.Context, *Get) (interface{}, error)
		visitIndexExpr(context.Context, *Index) (interface{}, error)
		visitVariableExpr(context.Context, *Variable) (interface{}, error)
		visitCallExpr(context.Context, *Call) (interface{}, error)
		visitArrayExpr(context.Context, *Array) (interface{}, error)

		visitParseErrorExpr(context.Context, *ParseError) (interface{}, error)
	}
	Interpreter interface {
		Evaluate(expr Expr) (interface{}, error)
	}
)
