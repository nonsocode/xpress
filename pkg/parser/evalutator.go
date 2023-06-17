package parser

import (
	"fmt"
)

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
		return i.sub(left, right)
	case SLASH:
		return i.div(left, right)
	case STAR:
		return i.mul(left, right)
	case PLUS:
		return i.add(left, right)
	case GREATER:
		return i.greater(left, right)
	case GREATER_EQUAL:
		return i.greaterEqual(left, right)
	case LESS:
		return i.less(left, right)
	case LESS_EQUAL:
		return i.lessEqual(left, right)
	case BANG_EQUAL:
		return !i.isEqual(left, right), nil
	case EQUAL_EQUAL:
		return i.isEqual(left, right), nil
	case AND:
		return i.isTruthy(left) && i.isTruthy(right), nil
	case OR:
		return i.isTruthy(left) || i.isTruthy(right), nil
	}
	return nil, nil
}

func (e *Evaluator) add(left, right interface{}) (interface{}, error) {
	if left == nil || right == nil {
		return nil, fmt.Errorf("cannot add nil values: adding %v and %v", left, right)
	}

	_, oks1 := left.(string)
	_, oks2 := right.(string)
	if oks1 || oks2 {
		return fmt.Sprintf("%v%v", left, right), nil
	}
	return left.(float64) + right.(float64), nil
}

func (e *Evaluator) sub(left, right interface{}) (interface{}, error) {
	if left == nil || right == nil {
		return nil, fmt.Errorf("cannot subtract nil values: adding %v and %v", left, right)
	}

	leftNum, ok1 := left.(float64)
	rightNum, ok2 := right.(float64)

	if !ok1 || !ok2 {
		return nil, fmt.Errorf("cannot subtract non-numbers or strings: %v - %v", left, right)
	}

	return leftNum - rightNum, nil
}

func (e *Evaluator) mul(left, right interface{}) (interface{}, error) {
	if left, ok := left.(float64); ok {
		if right, ok := right.(float64); ok {
			return left * right, nil
		}
	}
	return nil, fmt.Errorf("cannot multiply non-numbers: %v * %v", left, right)
}

func (e *Evaluator) div(left, right interface{}) (interface{}, error) {
	if left, ok := left.(float64); ok {
		if right, ok := right.(float64); ok {
			if right == 0 {
				return nil, fmt.Errorf("cannot divide by zero: %v / %v", left, right)
			}
			return left / right, nil
		}
	}
	return nil, fmt.Errorf("cannot divide non-numbers: %v / %v", left, right)
}
func (e *Evaluator) greater(left, right interface{}) (interface{}, error) {
	switch l := left.(type) {
	case float64:
		if r, ok := right.(float64); ok {
			return l > r, nil
		}
	case string:
		if r, ok := right.(string); ok {
			return l > r, nil
		}
	}
	return nil, fmt.Errorf("cannot compare %T with %T", left, right)
}

func (e *Evaluator) greaterEqual(left, right interface{}) (interface{}, error) {
	switch l := left.(type) {
	case float64:
		if r, ok := right.(float64); ok {
			return l >= r, nil
		}
	case string:
		if r, ok := right.(string); ok {
			return l >= r, nil
		}
	}
	return nil, fmt.Errorf("cannot compare %T with %T", left, right)
}

func (e *Evaluator) less(left, right interface{}) (interface{}, error) {
	switch l := left.(type) {
	case float64:
		if r, ok := right.(float64); ok {
			return l < r, nil
		}
	case string:
		if r, ok := right.(string); ok {
			return l < r, nil
		}
	}
	return nil, fmt.Errorf("cannot compare %T with %T", left, right)
}

func (e *Evaluator) lessEqual(left, right interface{}) (interface{}, error) {
	switch l := left.(type) {
	case float64:
		if r, ok := right.(float64); ok {
			return l <= r, nil
		}
	case string:
		if r, ok := right.(string); ok {
			return l <= r, nil
		}
	}
	return nil, fmt.Errorf("cannot compare %T with %T", left, right)
}

func (i *Evaluator) isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	if b == nil {
		return false
	}
	return a == b
}

// func (i *Evaluator) compare(a, b interface{}, op TokenType) bool, error {
// 	// both types have to be the same otherwise we can't compare them
// }

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
	obj, err := i.interpret(expr.object)
	if err != nil {
		return nil, err
	}
	indexValue, err := i.interpret(expr.index)
	if err != nil {
		return nil, err
	}
	// Assume that obj is a map from string to interface{} and get the field.
	switch obj.(type) {
	case []interface{}:
		indexNum, ok := indexValue.(float64)
		if !ok {
			return nil, fmt.Errorf("cannot index into array with non-number key %v", indexValue)
		}

		return obj.([]interface{})[int(indexNum)], nil
	case map[string]interface{}:
		indexValueStr, ok := indexValue.(string)
		if !ok {
			return nil, fmt.Errorf("cannot index into map with non-string key %v", indexValue)
		}
		return obj.(map[string]interface{})[indexValueStr], nil
	default:
		return nil, fmt.Errorf("cannot index into type %T", obj)
	}
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

func (i *Evaluator) visitArrayExpr(expr *Array) (interface{}, error) {
	values := make([]interface{}, 0)
	for _, v := range expr.values {
		value, err := i.interpret(v)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, nil
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
