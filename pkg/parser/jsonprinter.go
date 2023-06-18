package parser

import "encoding/json"

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
		JNode `json:",inline"`
		Name  string `json:"name"`
	}

	JParseError struct {
		JNode   `json:",inline"`
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

func (jp *JSONPrinter) visitBinaryExpr(expr *Binary) (interface{}, error) {
	left, _ := expr.Left().Accept(jp)
	right, _ := expr.Right().Accept(jp)
	return &JBinary{
		JNode:    JNode{Type: "Binary"},
		Left:     left,
		Operator: expr.Operator().Type().String(),
		Right:    right,
	}, nil
}

func (jp *JSONPrinter) visitGroupingExpr(expr *Grouping) (interface{}, error) {
	expression, _ := expr.Expression().Accept(jp)
	return &JGrouping{
		JNode:      JNode{Type: "Grouping"},
		Expression: expression,
	}, nil
}

func (jp *JSONPrinter) visitLiteralExpr(expr *Literal) (interface{}, error) {
	return &JLiteral{
		JNode: JNode{Type: "Literal"},
		Value: expr.Value(),
		Raw:   expr.Raw(),
	}, nil
}

func (jp *JSONPrinter) visitUnaryExpr(expr *Unary) (interface{}, error) {
	right, _ := expr.Right().Accept(jp)
	return &JUnary{
		JNode:    JNode{Type: "Unary"},
		Operator: expr.Operator().Type().String(),
		Right:    right,
	}, nil
}

func (jp *JSONPrinter) visitTemplateExpr(expr *Template) (interface{}, error) {
	exprs := make([]interface{}, len(expr.Expressions()))
	for i, e := range expr.Expressions() {
		exprs[i], _ = e.Accept(jp)
	}
	return &JTemplate{
		JNode: JNode{Type: "Template"},
		Exprs: exprs,
	}, nil
}

func (jp *JSONPrinter) visitTernaryExpr(expr *Ternary) (interface{}, error) {
	cond, _ := expr.Condition().Accept(jp)
	then, _ := expr.TrueExpr().Accept(jp)
	els, _ := expr.FalseExpr().Accept(jp)
	return &JTernary{
		JNode: JNode{Type: "Ternary"},
		Cond:  cond,
		Then:  then,
		Else:  els,
	}, nil
}

func (jp *JSONPrinter) visitGetExpr(expr *Get) (interface{}, error) {
	obj, _ := expr.Object().Accept(jp)

	return &JGet{
		JNode:      JNode{Type: "Get"},
		Object:     obj,
		Identifier: expr.Name().Lexeme(),
	}, nil
}

func (jp *JSONPrinter) visitIndexExpr(expr *Index) (interface{}, error) {
	obj, _ := expr.Object().Accept(jp)
	index, _ := expr.Index().Accept(jp)
	return &JIndex{
		JNode:  JNode{Type: "Index"},
		Object: obj,
		Index:  index,
	}, nil
}

func (jp *JSONPrinter) visitVariableExpr(expr *Variable) (interface{}, error) {
	return &JVariable{
		JNode: JNode{Type: "Variable"},
		Name:  expr.Name().String(),
	}, nil
}

func (jp *JSONPrinter) visitCallExpr(expr *Call) (interface{}, error) {
	callee, _ := expr.Callee().Accept(jp)
	args := make([]interface{}, len(expr.Arguments()))
	for i, a := range expr.Arguments() {
		args[i], _ = a.Accept(jp)
	}
	return &JCall{
		JNode:     JNode{Type: "Call"},
		Callee:    callee,
		Arguments: args,
	}, nil
}

func (jp *JSONPrinter) visitArrayExpr(expr *Array) (interface{}, error) {
	values := make([]interface{}, len(expr.Values()))
	for i, v := range expr.Values() {
		values[i], _ = v.Accept(jp)
	}
	return &JArray{
		JNode:  JNode{Type: "Array"},
		Values: values,
	}, nil
}

func (jp *JSONPrinter) visitParseErrorExpr(err *ParseError) (interface{}, error) {
	return &JParseError{
		JNode:   JNode{Type: "ParseError"},
		Message: err.Error(),
	}, nil
}

func (jp *JSONPrinter) Print(expr Expr) (string, error) {
	j, err := expr.Accept(jp)
	if err != nil {
		return "", err
	}
	b, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
