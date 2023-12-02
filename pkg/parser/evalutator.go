package parser

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
)

type (
	Evaluator struct {
		members map[string]interface{}
		timeout time.Duration
		lock    sync.RWMutex
	}

	EvaluationError struct {
		message string
	}
)

const (
	// DefaultTimeout is the default timeout for evaluating expressions.
	DefaultTimeout = 10 * time.Millisecond
)

var (
	EvaluationCancelledErrror = NewEvaluationError("evaluation cancelled")
)

func NewInterpreter() *Evaluator {
	return &Evaluator{
		members: make(map[string]interface{}),
		timeout: DefaultTimeout,
	}
}

func (i *Evaluator) SetTimeout(timeout time.Duration) {
	i.timeout = timeout
}

func NewEvaluationError(message string, args ...interface{}) *EvaluationError {
	return &EvaluationError{message: fmt.Sprintf(message, args...)}
}

func (e *EvaluationError) Error() string {
	return e.message
}

func (i *Evaluator) AddMember(name string, member interface{}) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.members[name] = member
	return nil
}

func (i *Evaluator) SetMembers(members map[string]interface{}) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.members = members
	return nil
}

func (i *Evaluator) visitBinaryExpr(ctx context.Context, expr *Binary) (interface{}, error) {
	if ctx.Err() != nil {
		return nil, EvaluationCancelledErrror
	}
	left, err := i.interpret(ctx, expr.left)
	if err != nil {
		return nil, err
	}
	right, err := i.interpret(ctx, expr.right)
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

	switch l := left.(type) {
	case string:
		if r, ok := right.(string); ok {
			return l + r, nil
		}
	case float64:
		if r, ok := right.(float64); ok {
			return l + r, nil
		}
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

func (i *Evaluator) visitParseErrorExpr(
	ctx context.Context,
	expr *ParseError,
) (interface{}, error) {
	if ctx.Err() != nil {
		return nil, EvaluationCancelledErrror
	}
	return nil, fmt.Errorf("parse error: %s", expr.Error())
}

func (i *Evaluator) visitGroupingExpr(ctx context.Context, expr *Grouping) (interface{}, error) {
	if ctx.Err() != nil {
		return nil, EvaluationCancelledErrror
	}
	return i.interpret(ctx, expr.expression)
}

func (i *Evaluator) visitLiteralExpr(ctx context.Context, expr *Literal) (interface{}, error) {
	if ctx.Err() != nil {
		return nil, EvaluationCancelledErrror
	}
	return expr.value, nil
}

func (i *Evaluator) visitUnaryExpr(ctx context.Context, expr *Unary) (interface{}, error) {
	if ctx.Err() != nil {
		return nil, EvaluationCancelledErrror
	}
	right, err := i.interpret(ctx, expr.right)
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

func (i *Evaluator) visitTemplateExpr(ctx context.Context, expr *Template) (interface{}, error) {
	if ctx.Err() != nil {
		return nil, EvaluationCancelledErrror
	}
	evaluations := make([]interface{}, 0)
	for _, e := range expr.expressions {
		ev, err := i.interpret(ctx, e)
		if err != nil {
			return nil, err
		}
		evaluations = append(evaluations, ev)
	}
	if len(evaluations) == 1 {
		return evaluations[0], nil
	}

	str := strings.Builder{}
	for _, e := range evaluations {
		str.WriteString(fmt.Sprintf("%v", e))
	}
	return str.String(), nil
}

func (i *Evaluator) visitTernaryExpr(ctx context.Context, expr *Ternary) (interface{}, error) {
	if ctx.Err() != nil {
		return nil, EvaluationCancelledErrror
	}
	condition, err := i.interpret(ctx, expr.condition)
	if err != nil {
		return nil, err
	}
	if i.isTruthy(condition) {
		return i.interpret(ctx, expr.trueExpr)
	}
	return i.interpret(ctx, expr.falseExpr)
}

func (i *Evaluator) visitVariableExpr(ctx context.Context, expr *Variable) (interface{}, error) {
	if ctx.Err() != nil {
		return nil, EvaluationCancelledErrror
	}
	if member, ok := i.members[expr.name.lexeme]; ok {
		return member, nil
	}
	return nil, nil
}

func (i *Evaluator) visitGetExpr(ctx context.Context, expr *Get) (interface{}, error) {
	if ctx.Err() != nil {
		return nil, EvaluationCancelledErrror
	}
	if expr.object == nil {
		return nil, fmt.Errorf("object is nil")
	}
	obj, err := i.interpret(ctx, expr.object)
	if err != nil {
		return nil, err
	}
	if obj == nil {
		return nil, fmt.Errorf("cannot get property '%s' of nil", expr.name.lexeme)
	}

	value := reflect.ValueOf(obj) // hack for pointer receivers

	switch value.Kind() {
	case reflect.Map:
		key := reflect.ValueOf(expr.name.lexeme)
		if value.MapIndex(key).IsValid() {
			return value.MapIndex(key).Interface(), nil
		}
		return nil, nil // TODO: return error?
	case reflect.Struct, reflect.Ptr:
		if field, ok := getFieldFromStructOrPointer(value, expr.name.lexeme); ok {
			return field, nil
		}
		if method, ok := getMethodFromStructOrPointer(value, expr.name.lexeme); ok {
			return method, nil
		}
		return nil, nil // TODO: return error?
	case reflect.Slice, reflect.Array, reflect.String:
		if expr.name.lexeme == "length" {
			return value.Len(), nil
		}
		return nil, fmt.Errorf("property '%s' does not exist", expr.name.lexeme)
	default:
		return nil, fmt.Errorf("cannot get property '%s' of type %T", expr.name.lexeme, obj)
	}
}

func getFieldFromStructOrPointer(value reflect.Value, name string) (interface{}, bool) {
	switch value.Kind() {
	case reflect.Ptr:
		value = value.Elem()
	case reflect.Struct:
	default:
		return nil, false
	}
	if field := value.FieldByName(name); field.IsValid() {
		return field.Interface(), true
	}
	return nil, false
}

func getMethodFromStructOrPointer(value reflect.Value, name string) (interface{}, bool) {
	if method := value.MethodByName(name); method.IsValid() {
		return method.Interface(), true
	}
	if value.CanAddr() {
		if method := value.Addr().MethodByName(name); method.IsValid() {
			return method.Interface(), true
		}
	}
	switch value.Kind() {
	case reflect.Ptr:
	case reflect.Struct:
		v := reflect.New(value.Type())
		v.Elem().Set(value)
		value = v
	default:
		return nil, false
	}
	if method := value.MethodByName(name); method.IsValid() {
		return method.Interface(), true
	}
	return nil, false
}

func (i *Evaluator) visitIndexExpr(ctx context.Context, expr *Index) (interface{}, error) {
	if ctx.Err() != nil {
		return nil, EvaluationCancelledErrror
	}
	obj, err := i.interpret(ctx, expr.object)
	if err != nil {
		return nil, err
	}
	indexValue, err := i.interpret(ctx, expr.index)
	if err != nil {
		return nil, err
	}

	if obj == nil {
		return nil, fmt.Errorf("cannot index into nil")
	}

	value := reflect.ValueOf(obj)

	switch value.Kind() {
	case reflect.Map:
		key := reflect.ValueOf(indexValue)
		if value.MapIndex(key).IsValid() {
			return value.MapIndex(key).Interface(), nil
		}
		return nil, nil // TODO: return error?
	case reflect.Struct, reflect.Ptr:
		key, ok := indexValue.(string)
		if !ok {
			return nil, fmt.Errorf("property '%s' does not exist", indexValue)
		}
		if field, ok := getFieldFromStructOrPointer(value, key); ok {
			return field, nil
		}
		if method, ok := getMethodFromStructOrPointer(value, key); ok {
			return method, nil
		}

		return nil, nil // TODO: return error?
	case reflect.Slice, reflect.Array, reflect.String:
		indexFloat, ok := indexValue.(float64)
		if !ok {
			return nil, fmt.Errorf("index '%v' is not an integer", indexValue)
		}
		index := int(indexFloat)
		if index < 0 || index >= value.Len() {
			return nil, fmt.Errorf("index '%v' is out of bounds", indexValue)
		}
		return value.Index(index).Interface(), nil
	default:
		return nil, fmt.Errorf("cannot index into type %T", obj)
	}
}

func (e *Evaluator) visitCallExpr(ctx context.Context, expr *Call) (interface{}, error) {
	if ctx.Err() != nil {
		return nil, EvaluationCancelledErrror
	}
	args := make([]interface{}, 0)
	for _, a := range expr.arguments {
		arg, err := e.interpret(ctx, a)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	callee, err := e.interpret(ctx, expr.callee)
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
	var argIndex int
	in := make([]reflect.Value, 0)
	if fn.Type().NumIn() > 0 && fn.Type().In(0) == reflect.TypeOf((*context.Context)(nil)).Elem() {
		in = append(in, reflect.ValueOf(ctx))
		argIndex = 1
	}
	if !isVariadic && fn.Type().NumIn() != (len(args)+argIndex) {
		return nil, NewEvaluationError(
			"function '%s' expects %d arguments, got %d",
			identifyCallee(expr),
			fn.Type().NumIn(),
			len(args),
		)
	}

	variadicIndex := fn.Type().NumIn() - 1
	for i, arg := range args {
		if isVariadic && i+argIndex >= variadicIndex {
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
		paramType := fn.Type().In(i + argIndex)

		if !argValue.Type().AssignableTo(paramType) {
			// attempt to convert arg to paramType
			if argValue.Type().ConvertibleTo(paramType) {
				argValue = argValue.Convert(paramType)
			} else {
				return nil, NewEvaluationError(
					"argument '%v' is not assignable to parameter '%s'",
					arg,
					paramType.String(),
				)
			}
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
	switch callee := expr.callee.(type) {
	case *Variable:
		return callee.name.lexeme
	case *Get:
		return callee.name.lexeme
	}
	return "unknown"
}

func (i *Evaluator) visitArrayExpr(ctx context.Context, expr *Array) (interface{}, error) {
	if ctx.Err() != nil {
		return nil, EvaluationCancelledErrror
	}
	values := make([]interface{}, len(expr.values))
	for index, v := range expr.values {
		value, err := i.interpret(ctx, v)
		if err != nil {
			return nil, err
		}
		values[index] = value
	}
	return values, nil
}

func (i *Evaluator) isTruthy(object interface{}) bool {
	switch value := object.(type) {
	case nil:
		return false
	case bool:
		return value
	default:
		return true
	}
}

func (i *Evaluator) Evaluate(ctx context.Context, expr Expr) (interface{}, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()
	ctx, cancel := context.WithTimeout(ctx, i.timeout)
	defer cancel()

	var result interface{}
	var err error
	var once sync.Once

	done := make(chan bool)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				once.Do(func() {
					err = NewEvaluationError("%v", r)
				})
			}
			close(done)
		}()
		obj, e := i.interpret(ctx, expr)
		once.Do(func() {
			result = obj
			err = e
		})
	}()

	select {
	case <-done:
		return result, err
	case <-ctx.Done():
		return nil, NewEvaluationError("evaluation timed out after %s", i.timeout.String())
	}
}

func (i *Evaluator) interpret(ctx context.Context, expr Expr) (interface{}, error) {
	if ctx.Err() != nil {
		return nil, EvaluationCancelledErrror
	}
	return expr.Accept(ctx, i)
}
