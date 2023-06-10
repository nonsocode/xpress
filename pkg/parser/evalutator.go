package parser

type (
	Evaluator struct{}
)

func NewInterpreter() *Evaluator {
	return &Evaluator{}
}

func (i *Evaluator) visitBinaryExpr(expr *Binary) interface{} {
	left := i.interpret(expr.left)
	right := i.interpret(expr.right)
	switch expr.operator.tokenType {
	case MINUS:
		return left.(float64) - right.(float64)
	case SLASH:
		return left.(float64) / right.(float64)
	case STAR:
		return left.(float64) * right.(float64)
	case PLUS:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l + r
			}
		}
		if l, ok := left.(string); ok {
			if r, ok := right.(string); ok {
				return l + r
			}
		}
	case GREATER:
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case LESS:
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		return left.(float64) <= right.(float64)
	}
	return nil
}

func (i *Evaluator) visitGroupingExpr(expr *Grouping) interface{} {
	return i.interpret(expr.expression)
}

func (i *Evaluator) visitLiteralExpr(expr *Literal) interface{} {
	if expr.value == nil {
		return "nil"
	}
	return expr.value
}

func (i *Evaluator) visitUnaryExpr(expr *Unary) interface{} {
	right := i.interpret(expr.right)
	switch expr.operator.tokenType {
	case MINUS:
		return -(right.(float64))
	case BANG:
		return !(i.isTruthy(right))
	}
	return nil
}

func (i *Evaluator) isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	if b, ok := object.(bool); ok {
		return b
	}
	return true
}

func (i *Evaluator) interpret(expr Expr) interface{} {
	return expr.accept(i)
}
