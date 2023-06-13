package parser

import "fmt"

type (
	Evaluator struct {
		funcs map[string]func(...interface{}) (interface{}, error)
	}
)

func NewInterpreter() *Evaluator {
	return &Evaluator{
		funcs: make(map[string]func(...interface{}) (interface{}, error)),
	}
}

func (i *Evaluator) AddFunc(name string, fn func(...interface{}) (interface{}, error)) {
	i.funcs[name] = fn
}

func (i *Evaluator) visitBinaryExpr(expr *Binary) (interface{}, error) {
	left, err := i.interpret(expr.left)
	if err != nil {
		return nil, err
	}
	right, err := i.interpret(expr.right)
	if err != nil {
		return nil, err
	}
	switch expr.operator.tokenType {
	case MINUS:
		return left.(float64) - right.(float64), nil
	case SLASH:
		return left.(float64) / right.(float64), nil
	case STAR:
		return left.(float64) * right.(float64), nil
	case PLUS:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l + r, nil
			}
		}
		if l, ok := left.(string); ok {
			if r, ok := right.(string); ok {
				return l + r, nil
			}
		}
	case GREATER:
		return left.(float64) > right.(float64), nil
	case GREATER_EQUAL:
		return left.(float64) >= right.(float64), nil
	case LESS:
		return left.(float64) < right.(float64), nil
	case LESS_EQUAL:
		return left.(float64) <= right.(float64), nil
	}
	return nil, nil
}

func (i *Evaluator) visitGroupingExpr(expr *Grouping) (interface{}, error) {
	return i.interpret(expr.expression)
}

func (i *Evaluator) visitLiteralExpr(expr *Literal) (interface{}, error) {
	return expr.value, nil
}

func (i *Evaluator) visitUnaryExpr(expr *Unary) (interface{}, error) {
	right, err := i.interpret(expr.right)
	if err != nil {
		return nil, err
	}
	switch expr.operator.tokenType {
	case MINUS:
		return -(right.(float64)), nil
	case BANG:
		return !(i.isTruthy(right)), nil
	}
	return nil, nil
}

func (i *Evaluator) visitTemplateExpr(expr *Template) (interface{}, error) {
	evaluations := make([]interface{}, 0)
	for _, e := range expr.expressions {
		ev, err := i.interpret(e)
		if err != nil {
			return nil, err
		}
		evaluations = append(evaluations, ev)
	}
	if len(evaluations) == 1 {
		return evaluations[0], nil
	}
	// conert to string and concat
	ret := ""
	for _, e := range evaluations {
		ret += fmt.Sprintf("%v", e)
	}
	return ret, nil
}

func (i *Evaluator) visitTernaryExpr(expr *Ternary) (interface{}, error) {
	condition, err := i.interpret(expr.condition)
	if err != nil {
		return nil, err
	}
	if i.isTruthy(condition) {
		return i.interpret(expr.trueExpr)
	}
	return i.interpret(expr.falseExpr)
}

func (i *Evaluator) visitVariableExpr(expr *Variable) (interface{}, error) {
	if fn, ok := i.funcs[expr.name.lexeme]; ok {
		return fn, nil
	}
	return nil, nil
}

func (i *Evaluator) visitGetExpr(expr *Get) (interface{}, error) {
	obj, err := i.interpret(expr.object)
	if err != nil {
		return nil, err
	}

	// Assume that obj is a map from string to interface{} and get the field.
	// TODO: Check if obj is actually a map and handle errors.
	return obj.(map[string]interface{})[expr.name.lexeme], nil
}

func (i *Evaluator) visitIndexExpr(expr *Index) (interface{}, error) {
	// TODO: implement
	return expr.object, nil
}

func (e *Evaluator) visitCallExpr(expr *Call) (interface{}, error) {
	args := make([]interface{}, 0)
	for _, a := range expr.arguments {
		arg, err := e.interpret(a)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	callee, err := e.interpret(expr.callee)
	if err != nil {
		return nil, err
	}
	if fn, ok := callee.(func(...interface{}) (interface{}, error)); ok {
		return fn(args...)
	}
	return nil, nil
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

func (i *Evaluator) interpret(expr Expr) (interface{}, error) {
	return expr.accept(i)
}
