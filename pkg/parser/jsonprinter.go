package parser

import (
	"context"
	"encoding/json"
)

type (
	JSONPrinter struct { // implements Visitor
	}
	JNode struct {
		Type string `json:"type"`
	}

	JBinary struct {
		JNode    `json:",inline"`
		Left     interface{} `json:"left"`
		Operator string      `json:"operator"`
		Right    interface{} `json:"right"`
	}

	JGrouping struct {
		JNode      `json:",inline"`
		Expression interface{} `json:"expression"`
	}

	JLiteral struct {
		JNode `json:",inline"`
		Value interface{} `json:"value"`
		Raw   string      `json:"raw"`
	}

	JUnary struct {
		JNode    `json:",inline"`
		Operator string      `json:"operator"`
		Right    interface{} `json:"right"`
	}

	JTemplate struct {
		JNode `json:",inline"`
		Exprs []interface{} `json:"exprs"`
	}

	JTernary struct {
		JNode `json:",inline"`
		Cond  interface{} `json:"cond"`
		Then  interface{} `json:"then"`
		Else  interface{} `json:"else"`
	}

	JGet struct {
		JNode      `json:",inline"`
		Object     interface{} `json:"object"`
		Identifier string      `json:"identifier"`
	}

	JIndex struct {
		JNode  `json:",inline"`
		Object interface{} `json:"object"`
		Index  interface{} `json:"index"`
	}

	JArray struct {
		JNode  `json:",inline"`
		Values []interface{} `json:"values"`
	}

	JVariable struct {
		JNode `       json:",inline"`
		Name  string `json:"name"`
	}

	JParseError struct {
		JNode   `       json:",inline"`
		Message string `json:"message"`
	}

	JCall struct {
		JNode     `json:",inline"`
		Callee    interface{}   `json:"callee"`
		Arguments []interface{} `json:"arguments"`
	}
)

func NewJSONPrinter() *JSONPrinter {
	return &JSONPrinter{}
}

func (jp *JSONPrinter) visitBinaryExpr(ctx context.Context, expr *Binary) (interface{}, error) {
	left, _ := expr.left.Accept(ctx, jp)
	right, _ := expr.right.Accept(ctx, jp)
	return &JBinary{
		JNode:    JNode{Type: "Binary"},
		Left:     left,
		Operator: expr.operator.tokenType.String(),
		Right:    right,
	}, nil
}

func (jp *JSONPrinter) visitGroupingExpr(ctx context.Context, expr *Grouping) (interface{}, error) {
	expression, _ := expr.expression.Accept(ctx, jp)
	return &JGrouping{
		JNode:      JNode{Type: "Grouping"},
		Expression: expression,
	}, nil
}

func (jp *JSONPrinter) visitLiteralExpr(ctx context.Context, expr *Literal) (interface{}, error) {
	return &JLiteral{
		JNode: JNode{Type: "Literal"},
		Value: expr.value,
		Raw:   expr.raw,
	}, nil
}

func (jp *JSONPrinter) visitUnaryExpr(ctx context.Context, expr *Unary) (interface{}, error) {
	right, _ := expr.right.Accept(ctx, jp)
	return &JUnary{
		JNode:    JNode{Type: "Unary"},
		Operator: expr.operator.tokenType.String(),
		Right:    right,
	}, nil
}

func (jp *JSONPrinter) visitTemplateExpr(ctx context.Context, expr *Template) (interface{}, error) {
	exprs := make([]interface{}, len(expr.expressions))
	for i, e := range expr.expressions {
		exprs[i], _ = e.Accept(ctx, jp)
	}
	return &JTemplate{
		JNode: JNode{Type: "Template"},
		Exprs: exprs,
	}, nil
}

func (jp *JSONPrinter) visitTernaryExpr(ctx context.Context, expr *Ternary) (interface{}, error) {
	cond, _ := expr.condition.Accept(ctx, jp)
	then, _ := expr.trueExpr.Accept(ctx, jp)
	els, _ := expr.falseExpr.Accept(ctx, jp)
	return &JTernary{
		JNode: JNode{Type: "Ternary"},
		Cond:  cond,
		Then:  then,
		Else:  els,
	}, nil
}

func (jp *JSONPrinter) visitGetExpr(ctx context.Context, expr *Get) (interface{}, error) {
	obj, _ := expr.object.Accept(ctx, jp)

	return &JGet{
		JNode:      JNode{Type: "Get"},
		Object:     obj,
		Identifier: expr.name.lexeme,
	}, nil
}

func (jp *JSONPrinter) visitIndexExpr(ctx context.Context, expr *Index) (interface{}, error) {
	obj, _ := expr.object.Accept(ctx, jp)
	index, _ := expr.index.Accept(ctx, jp)
	return &JIndex{
		JNode:  JNode{Type: "Index"},
		Object: obj,
		Index:  index,
	}, nil
}

func (jp *JSONPrinter) visitVariableExpr(ctx context.Context, expr *Variable) (interface{}, error) {
	return &JVariable{
		JNode: JNode{Type: "Variable"},
		Name:  expr.name.String(),
	}, nil
}

func (jp *JSONPrinter) visitCallExpr(ctx context.Context, expr *Call) (interface{}, error) {
	callee, _ := expr.callee.Accept(ctx, jp)
	args := make([]interface{}, len(expr.arguments))
	for i, a := range expr.arguments {
		args[i], _ = a.Accept(ctx, jp)
	}
	return &JCall{
		JNode:     JNode{Type: "Call"},
		Callee:    callee,
		Arguments: args,
	}, nil
}

func (jp *JSONPrinter) visitArrayExpr(ctx context.Context, expr *Array) (interface{}, error) {
	values := make([]interface{}, len(expr.values))
	for i, v := range expr.values {
		values[i], _ = v.Accept(ctx, jp)
	}
	return &JArray{
		JNode:  JNode{Type: "Array"},
		Values: values,
	}, nil
}

func (jp *JSONPrinter) visitParseErrorExpr(
	ctx context.Context,
	err *ParseError,
) (interface{}, error) {
	return &JParseError{
		JNode:   JNode{Type: "ParseError"},
		Message: err.Error(),
	}, nil
}

func (jp *JSONPrinter) Print(expr Expr) (string, error) {
	ctx := context.Background()
	j, err := expr.Accept(ctx, jp)
	if err != nil {
		return "", err
	}
	b, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
