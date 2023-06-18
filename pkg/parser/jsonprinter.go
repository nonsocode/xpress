package parser

type (
	JSONPrinter struct { // implements Visitor
		depth int
	}
)

func (jp *JSONPrinter) visitBinaryExpr(expr *Binary) (interface{}, error) {
	panic("not implemented") // TODO: Implement
}

func (jp *JSONPrinter) visitGroupingExpr(expr *Grouping) (interface{}, error) {
	return nil, nil
}

func (jp *JSONPrinter) visitLiteralExpr(expr *Literal) (interface{}, error) {
	return nil, nil
}

func (jp *JSONPrinter) visitUnaryExpr(expr *Unary) (interface{}, error) {
	return nil, nil
}
func (jp *JSONPrinter) visitTemplateExpr(expr *Template) (interface{}, error) {
	panic("not implemented") // TODO: Implement
}

func (jp *JSONPrinter) visitTernaryExpr(expr *Ternary) (interface{}, error) {
	panic("not implemented") // TODO: Implement
}

func (jp *JSONPrinter) visitGetExpr(expr *Get) (interface{}, error) {
	panic("not implemented") // TODO: Implement
}

func (jp *JSONPrinter) visitIndexExpr(expr *Index) (interface{}, error) {
	panic("not implemented") // TODO: Implement
}

func (jp *JSONPrinter) visitVariableExpr(expr *Variable) (interface{}, error) {
	panic("not implemented") // TODO: Implement
}

func (jp *JSONPrinter) visitCallExpr(expr *Call) (interface{}, error) {
	panic("not implemented") // TODO: Implement
}

func (jp *JSONPrinter) visitArrayExpr(expr *Array) (interface{}, error) {
	panic("not implemented") // TODO: Implement
}

func (jp *JSONPrinter) visitParseErrorExpr(err *ParseError) (interface{}, error) {
	panic("not implemented") // TODO: Implement
}
