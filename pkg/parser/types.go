package parser

import (
	"context"
)

type (
	Visitor interface {
		visitBinaryExpr(context.Context, *Binary) EvaluationResult
		visitGroupingExpr(context.Context, *Grouping) EvaluationResult
		visitLiteralExpr(context.Context, *Literal) EvaluationResult
		visitUnaryExpr(context.Context, *Unary) EvaluationResult

		visitTemplateExpr(context.Context, *Template) EvaluationResult
		visitTernaryExpr(context.Context, *Ternary) EvaluationResult
		visitGetExpr(context.Context, *Get) EvaluationResult
		visitIndexExpr(context.Context, *Index) EvaluationResult
		visitVariableExpr(context.Context, *Variable) EvaluationResult
		visitCallExpr(context.Context, *Call) EvaluationResult
		visitArrayExpr(context.Context, *Array) EvaluationResult
		visitMapExpr(context.Context, *Map) EvaluationResult
		visitMapEntryExpr(context.Context, *MapEntry) EvaluationResult

		visitParseErrorExpr(context.Context, *ParseError) EvaluationResult
	}
	Interpreter interface {
		Evaluate(expr Expr) EvaluationResult
	}

	EvaluationResult interface {
		Get() interface{}
		Error() error
	}

	optionalEvaluationResult struct {
		result
		absent bool
	}
	result struct {
		err   error
		value interface{}
	}
)

var (
	_ EvaluationResult = &optionalEvaluationResult{}
	_ EvaluationResult = &result{}
)

func (o *optionalEvaluationResult) IsAbsent() bool {
	return o.absent
}

func (c *result) Get() interface{} {
	return c.value
}

func (c *result) Error() error {
	return c.err
}
