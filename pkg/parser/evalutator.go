package parser

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

type (
	Evaluator struct {
		funcs   map[string]interface{}
		timeout time.Duration
	}

	EvaluationError struct {
		message string
	}
)

const (
	// DefaultTimeout is the default timeout for evaluating expressions.
	DefaultTimeout = 5 * time.Hour
)

func NewInterpreter() *Evaluator {
	return &Evaluator{
		funcs:   make(map[string]interface{}),
		timeout: DefaultTimeout,
	}
}

func NewEvaluationError(message string, args ...interface{}) *EvaluationError {
	return &EvaluationError{message: fmt.Sprintf(message, args...)}
}

func (e *EvaluationError) Error() string {
	return e.message
}

func (i *Evaluator) AddFunc(name string, fn interface{}) error {
	i.funcs[name] = fn
	return nil
}

func (i *Evaluator) SetFunctions(funcs map[string]interface{}) error {
	for name, fn := range funcs {
		i.AddFunc(name, fn)
	}
	return nil
}

func (i *Evaluator) visitBinaryExpr(ctx context.Context, expr *Binary) (interface{}, error) {
	left, err := i.interpret(ctx, expr.Left())
	if err != nil {
		return nil, err
	}
	right, err := i.interpret(ctx, expr.Right())
	if err != nil {
		return nil, err
	}
	switch expr.Operator().Type() {
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

	l, okf1 := left.(float64)
	r, okf2 := right.(float64)
	if okf1 || okf2 {
		return l + r, nil
	}
	return nil, fmt.Errorf("cannot add non-numbers or strings: %v + %v", left, right)
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

func (i *Evaluator) visitParseErrorExpr(ctx context.Context, expr *ParseError) (interface{}, error) {
	return nil, fmt.Errorf("parse error: %s", expr.Error())
}

func (i *Evaluator) visitGroupingExpr(ctx context.Context, expr *Grouping) (interface{}, error) {
	return i.interpret(ctx, expr.Expression())
}

func (i *Evaluator) visitLiteralExpr(ctx context.Context, expr *Literal) (interface{}, error) {
	return expr.Value(), nil
}

func (i *Evaluator) visitUnaryExpr(ctx context.Context, expr *Unary) (interface{}, error) {
	right, err := i.interpret(ctx, expr.Right())
	if err != nil {
		return nil, err
	}
	switch expr.Operator().Type() {
	case MINUS:
		return -(right.(float64)), nil
	case BANG:
		return !(i.isTruthy(right)), nil
	}
	return nil, nil
}

func (i *Evaluator) visitTemplateExpr(ctx context.Context, expr *Template) (interface{}, error) {
	evaluations := make([]interface{}, 0)
	for _, e := range expr.Expressions() {
		ev, err := i.interpret(ctx, e)
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

func (i *Evaluator) visitTernaryExpr(ctx context.Context, expr *Ternary) (interface{}, error) {
	condition, err := i.interpret(ctx, expr.Condition())
	if err != nil {
		return nil, err
	}
	if i.isTruthy(condition) {
		return i.interpret(ctx, expr.TrueExpr())
	}
	return i.interpret(ctx, expr.FalseExpr())
}

func (i *Evaluator) visitVariableExpr(ctx context.Context, expr *Variable) (interface{}, error) {
	if fn, ok := i.funcs[expr.Name().Lexeme()]; ok {
		return fn, nil
	}
	return nil, nil
}

func (i *Evaluator) visitGetExpr(ctx context.Context, expr *Get) (interface{}, error) {
	obj, err := i.interpret(ctx, expr.Object())
	if err != nil {
		return nil, err
	}

	// Assume that obj is a map from string to interface{} and get the field.
	// TODO: Check if obj is actually a map and handle errors.
	return obj.(map[string]interface{})[expr.Name().Lexeme()], nil
}

func (i *Evaluator) visitIndexExpr(ctx context.Context, expr *Index) (interface{}, error) {
	obj, err := i.interpret(ctx, expr.Object())
	if err != nil {
		return nil, err
	}
	indexValue, err := i.interpret(ctx, expr.Index())
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

func (e *Evaluator) visitCallExpr(ctx context.Context, expr *Call) (interface{}, error) {
	args := make([]interface{}, 0)
	for _, a := range expr.Arguments() {
		arg, err := e.interpret(ctx, a)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	callee, err := e.interpret(ctx, expr.Callee())
	if err != nil {
		return nil, err
	}
	fn := reflect.ValueOf(callee)
	if fn.Kind() != reflect.Func {
		return nil, NewEvaluationError(
			"cannot call non-function '%s' of type %T",
			identifyCallee(expr),
			callee,
		)
	}
	if fn.Type().NumOut() > 2 {
		return nil, NewEvaluationError(
			"function '%s' returns more than 2 values",
			identifyCallee(expr),
		)
	}
	if fn.Type().NumOut() == 2 {
		if fn.Type().Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
			return nil, NewEvaluationError(
				"function '%s' second return value must be of type error",
				identifyCallee(expr),
			)
		}
	}

	isVariadic := fn.Type().IsVariadic()
	if !isVariadic && fn.Type().NumIn() != len(args) {
		return nil, NewEvaluationError(
			"function '%s' expects %d arguments, got %d",
			identifyCallee(expr),
			fn.Type().NumIn(),
			len(args),
		)
	}

	in := make([]reflect.Value, 0)
	variadicIndex := fn.Type().NumIn() - 1
	for i, arg := range args {
		if isVariadic && i >= variadicIndex {
			// Variadic argument
			varsType := fn.Type().In(variadicIndex)
			paramType := varsType.Elem()
			for _, a := range args[i:] {
				if !reflect.TypeOf(a).AssignableTo(paramType) {
					return nil, NewEvaluationError(
						"variadic argument '%v' is not assignable to type '%s'",
						arg,
						paramType.String(),
					)
				}
				in = append(in, reflect.ValueOf(a))
			}
			break
		}
		argValue := reflect.ValueOf(arg)
		paramType := fn.Type().In(i)

		if !argValue.Type().AssignableTo(paramType) {
			return nil, NewEvaluationError(
				"argument '%v' is not assignable to parameter '%s'",
				arg,
				paramType.String(),
			)
		}

		in = append(in, argValue)
	}

	out := fn.Call(in)
	if len(out) == 2 {
		if out[1].Interface() != nil {
			return nil, out[1].Interface().(error)
		}
	}
	return out[0].Interface(), nil
}

func identifyCallee(expr *Call) string {
	switch callee := expr.Callee().(type) {
	case *Variable:
		return callee.Name().Lexeme()
	case *Get:
		return callee.Name().Lexeme()
	}
	return "unknown"
}

func (i *Evaluator) visitArrayExpr(ctx context.Context, expr *Array) (interface{}, error) {
	values := make([]interface{}, 0)
	for _, v := range expr.Values() {
		value, err := i.interpret(ctx, v)
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

func (i *Evaluator) Evaluate(expr Expr) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), i.timeout)
	defer cancel()

	result := make(chan interface{})
	errChan := make(chan error)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				errChan <- NewEvaluationError("%v", r)
			}
			close(result)
			close(errChan)
		}()
		obj, err := i.interpret(ctx, expr)
		if err != nil {
			errChan <- err
			return
		}
		result <- obj
	}()

	select {
	case obj := <-result:
		return obj, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, NewEvaluationError("evaluation timed out after %s", i.timeout.String())
	}
}

func (i *Evaluator) interpret(ctx context.Context, expr Expr) (interface{}, error) {
	return expr.Accept(ctx, i)
}
