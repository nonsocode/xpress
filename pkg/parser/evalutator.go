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

func (i *Evaluator) visitBinaryExpr(ctx context.Context, expr *Binary) EvaluationResult {
	if ctx.Err() != nil {
		return &result{err: EvaluationCancelledErrror}
	}

	res := i.interpret(ctx, expr.left)
	if res.Error() != nil {
		return res
	}
	left := res.Get()

	switch expr.operator.tokenType {
	case AND:
		if !i.isTruthy(left) {
			return &result{value: false}
		}
	case OR:
		if i.isTruthy(left) {
			return &result{value: true}
		}
	case NULLCOALESCING:
		if i.isTruthy(left) {
			return &result{value: left}
		}
	}

	res = i.interpret(ctx, expr.right)
	if res.Error() != nil {
		return res
	}
	right := res.Get()
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
		return &result{value: !i.isEqual(left, right)}
	case EQUAL_EQUAL:
		return &result{value: i.isEqual(left, right)}
	case AND:
		return &result{value: i.isTruthy(left) && i.isTruthy(right)}
	case OR:
		return &result{value: i.isTruthy(left) || i.isTruthy(right)}
	case NULLCOALESCING:
		if i.isTruthy(left) {
			return &result{value: left}
		}
		return &result{value: right}
	}
	return &result{}
}

func (e *Evaluator) add(left, right interface{}) EvaluationResult {
	if left == nil || right == nil {
		return &result{err: fmt.Errorf("cannot add nil values: adding %v and %v", left, right)}
	}

	switch l := left.(type) {
	case string:
		if r, ok := right.(string); ok {
			return &result{value: l + r}
		}
	}
	if areNumbers(left, right) {
		leftNum, _ := toFloat64(left)
		rightNum, _ := toFloat64(right)
		return &result{value: leftNum + rightNum}
	}

	return &result{err: fmt.Errorf("cannot add non-numbers or strings: %v + %v", left, right)}
}

func (e *Evaluator) sub(left, right interface{}) EvaluationResult {
	if left == nil || right == nil {
		return &result{err: fmt.Errorf("cannot subtract nil values: adding %v and %v", left, right)}
	}

	if areNumbers(left, right) {
		leftNum, _ := toFloat64(left)
		rightNum, _ := toFloat64(right)
		return &result{value: leftNum - rightNum}
	}
	return &result{err: fmt.Errorf("cannot subtract non-numbers or strings: %v - %v", left, right)}
}

func (e *Evaluator) mul(left, right interface{}) EvaluationResult {
	if areNumbers(left, right) {
		leftNum, _ := toFloat64(left)
		rightNum, _ := toFloat64(right)
		return &result{value: leftNum * rightNum}
	}
	return &result{err: fmt.Errorf("cannot multiply non-numbers: %v * %v", left, right)}
}

func (e *Evaluator) div(left, right interface{}) EvaluationResult {
	if areNumbers(left, right) {
		leftNum, _ := toFloat64(left)
		rightNum, _ := toFloat64(right)
		if rightNum == 0 {
			return &result{err: fmt.Errorf("cannot divide by zero: %f / %f", leftNum, rightNum)}
		}
		return &result{value: leftNum / rightNum}
	}
	return &result{err: fmt.Errorf("cannot divide non-numbers: %v / %v", left, right)}
}
func (e *Evaluator) greater(left, right interface{}) EvaluationResult {
	if areNumbers(left, right) {
		leftNum, _ := toFloat64(left)
		rightNum, _ := toFloat64(right)
		return &result{value: leftNum > rightNum}
	}
	switch l := left.(type) {
	case string:
		if r, ok := right.(string); ok {
			return &result{value: l > r}
		}
	}
	return &result{err: fmt.Errorf("cannot compare %T with %T", left, right)}
}

func (e *Evaluator) greaterEqual(left, right interface{}) EvaluationResult {
	if areNumbers(left, right) {
		leftNum, _ := toFloat64(left)
		rightNum, _ := toFloat64(right)
		return &result{value: leftNum >= rightNum}
	}
	switch l := left.(type) {
	case string:
		if r, ok := right.(string); ok {
			return &result{value: l >= r}
		}
	}
	return &result{err: fmt.Errorf("cannot compare %T with %T", left, right)}
}

func (e *Evaluator) less(left, right interface{}) EvaluationResult {
	if areNumbers(left, right) {
		leftNum, _ := toFloat64(left)
		rightNum, _ := toFloat64(right)
		return &result{value: leftNum < rightNum}
	}
	switch l := left.(type) {
	case string:
		if r, ok := right.(string); ok {
			return &result{value: l < r}
		}
	}
	return &result{err: fmt.Errorf("cannot compare %T with %T", left, right)}
}

func (e *Evaluator) lessEqual(left, right interface{}) EvaluationResult {
	if areNumbers(left, right) {
		leftNum, _ := toFloat64(left)
		rightNum, _ := toFloat64(right)
		return &result{value: leftNum <= rightNum}
	}
	switch l := left.(type) {
	case string:
		if r, ok := right.(string); ok {
			return &result{value: l <= r}
		}
	}
	return &result{err: fmt.Errorf("cannot compare %T with %T", left, right)}
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
) EvaluationResult {
	if ctx.Err() != nil {
		return &result{err: EvaluationCancelledErrror}
	}
	return &result{err: fmt.Errorf("parse error: %w", expr)}
}

func (i *Evaluator) visitGroupingExpr(ctx context.Context, expr *Grouping) EvaluationResult {
	if ctx.Err() != nil {
		return &result{err: EvaluationCancelledErrror}
	}
	return i.interpret(ctx, expr.expression)
}

func (i *Evaluator) visitLiteralExpr(ctx context.Context, expr *Literal) EvaluationResult {
	if ctx.Err() != nil {
		return &result{err: EvaluationCancelledErrror}
	}
	return &result{value: expr.value}
}

func (i *Evaluator) visitUnaryExpr(ctx context.Context, expr *Unary) EvaluationResult {
	if ctx.Err() != nil {
		return &result{err: EvaluationCancelledErrror}
	}
	res := i.interpret(ctx, expr.right)
	if res.Error() != nil {
		return res
	}
	right := res.Get()
	switch expr.operator.tokenType {
	case MINUS:
		return &result{value: -(right.(float64))} // TODO: handle other numeric types
	case BANG:
		return &result{value: !(i.isTruthy(right))}
	}
	return &result{}
}

func (i *Evaluator) visitTemplateExpr(ctx context.Context, expr *Template) EvaluationResult {
	if ctx.Err() != nil {
		return &result{err: EvaluationCancelledErrror}
	}
	evaluations := make([]interface{}, 0)
	for _, e := range expr.expressions {
		res := i.interpret(ctx, e)
		if res.Error() != nil {
			return res
		}
		evaluations = append(evaluations, res.Get())
	}
	if len(evaluations) == 1 {
		return &result{value: evaluations[0]}
	}

	str := strings.Builder{}
	for _, e := range evaluations {
		str.WriteString(fmt.Sprintf("%v", e))
	}
	return &result{value: str.String()}
}

func (i *Evaluator) visitTernaryExpr(ctx context.Context, expr *Ternary) EvaluationResult {
	if ctx.Err() != nil {
		return &result{err: EvaluationCancelledErrror}
	}
	res := i.interpret(ctx, expr.condition)
	if res.Error() != nil {
		return res
	}
	if i.isTruthy(res.Get()) {
		return i.interpret(ctx, expr.trueExpr)
	}
	return i.interpret(ctx, expr.falseExpr)
}

func (i *Evaluator) visitVariableExpr(ctx context.Context, expr *Variable) EvaluationResult {
	if ctx.Err() != nil {
		return &result{err: EvaluationCancelledErrror}
	}
	if member, ok := i.members[expr.name.lexeme]; ok {
		return &result{value: member}
	}
	return &result{}
}

func (i *Evaluator) visitOptionalExpr(ctx context.Context, expr *Optional) EvaluationResult {
	if ctx.Err() != nil {
		return &result{err: EvaluationCancelledErrror}
	}
	res := i.interpret(ctx, expr.left)
	if res.Error() != nil {
		return res
	}
	if res, ok := res.(*optionalEvaluationResult); ok && res.IsAbsent() {
		return res
	}
	if res.Get() == nil {
		return &optionalEvaluationResult{absent: true}
	}
	return res
}

func (i *Evaluator) visitGetExpr(ctx context.Context, expr *Get) EvaluationResult {
	if ctx.Err() != nil {
		return &result{err: EvaluationCancelledErrror}
	}
	res := i.interpret(ctx, expr.object)
	if res.Error() != nil {
		return res
	}
	if res, ok := res.(*optionalEvaluationResult); ok && res.IsAbsent() {
		return res
	}
	obj := res.Get()
	if obj == nil {
		return &result{err: fmt.Errorf("cannot get property '%s' of nil", expr.name.lexeme)}
	}

	value := reflect.ValueOf(obj) // hack for pointer receivers

	switch value.Kind() {
	case reflect.Map:
		key := reflect.ValueOf(expr.name.lexeme)
		if value.MapIndex(key).IsValid() {
			return &result{value: value.MapIndex(key).Interface()}
		}
		return &result{} // TODO: return error?
	case reflect.Struct, reflect.Ptr:
		if field, ok := getFieldFromStructOrPointer(value, expr.name.lexeme); ok {
			return &result{value: field}
		}
		if method, ok := getMethodFromStructOrPointer(value, expr.name.lexeme); ok {
			return &result{value: method}
		}
		return &result{} // TODO: return error?
	case reflect.Slice, reflect.Array, reflect.String:
		if expr.name.lexeme == "length" {
			return &result{value: float64(value.Len())}
		}
		return &result{err: fmt.Errorf("property '%s' does not exist", expr.name.lexeme)}
	default:
		return &result{err: fmt.Errorf("cannot get property '%s' of type %T", expr.name.lexeme, obj)}
	}
}

func (i *Evaluator) visitIndexExpr(ctx context.Context, expr *Index) EvaluationResult {
	if ctx.Err() != nil {
		return &result{err: EvaluationCancelledErrror}
	}
	res := i.interpret(ctx, expr.object)
	if res.Error() != nil {
		return res
	}
	if res, ok := res.(*optionalEvaluationResult); ok && res.IsAbsent() {
		return res
	}
	obj := res.Get()
	res = i.interpret(ctx, expr.index)
	if res.Error() != nil {
		return res
	}
	indexValue := res.Get()

	if obj == nil {
		return &result{err: fmt.Errorf("cannot index into nil")}
	}

	value := reflect.ValueOf(obj)

	switch value.Kind() {
	case reflect.Map:
		key := reflect.ValueOf(indexValue)
		if value.MapIndex(key).IsValid() {
			return &result{value: value.MapIndex(key).Interface()}
		}
		return &result{} // TODO: return error?
	case reflect.Struct, reflect.Ptr:
		key, ok := indexValue.(string)
		if !ok {
			return &result{err: fmt.Errorf("property '%s' does not exist", indexValue)}
		}
		if field, ok := getFieldFromStructOrPointer(value, key); ok {
			return &result{value: field}
		}
		if method, ok := getMethodFromStructOrPointer(value, key); ok {
			return &result{value: method}
		}

		return &result{} // TODO: return error?
	case reflect.Slice, reflect.Array, reflect.String:
		indexFloat, ok := indexValue.(float64)
		if !ok {
			return &result{err: fmt.Errorf("index '%v' is not an integer", indexValue)}
		}
		index := int(indexFloat)
		if index < 0 || index >= value.Len() {
			return &result{err: fmt.Errorf("index '%v' is out of bounds", indexValue)}
		}
		return &result{value: value.Index(index).Interface()}
	default:
		return &result{err: fmt.Errorf("cannot index into type %T", obj)}
	}
}

func (e *Evaluator) visitCallExpr(ctx context.Context, expr *Call) EvaluationResult {
	if ctx.Err() != nil {
		return &result{err: EvaluationCancelledErrror}
	}

	calleeRes := e.interpret(ctx, expr.callee)
	if calleeRes.Error() != nil {
		return calleeRes
	}
	if calleeRes, ok := calleeRes.(*optionalEvaluationResult); ok && calleeRes.IsAbsent() {
		return calleeRes
	}
	callee := calleeRes.Get()
	fn := reflect.ValueOf(callee)
	if fn.Kind() != reflect.Func {
		return &result{err: NewEvaluationError(
			"cannot call non-function '%s' of type %T",
			identifyCallee(expr),
			callee,
		)}
	}
	if fn.Type().NumOut() > 2 {
		return &result{err: NewEvaluationError(
			"function '%s' returns more than 2 values",
			identifyCallee(expr),
		)}
	}
	if fn.Type().NumOut() == 2 {
		if fn.Type().Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
			return &result{err: NewEvaluationError(
				"function '%s' second return value must be of type error",
				identifyCallee(expr),
			)}
		}
	}

	args := make([]interface{}, 0)
	for _, a := range expr.arguments {
		res := e.interpret(ctx, a)
		if res.Error() != nil {
			return res
		}
		args = append(args, res.Get())
	}

	isVariadic := fn.Type().IsVariadic()
	var argIndex int
	in := make([]reflect.Value, 0)
	if fn.Type().NumIn() > 0 && fn.Type().In(0) == reflect.TypeOf((*context.Context)(nil)).Elem() {
		in = append(in, reflect.ValueOf(ctx))
		argIndex = 1
	}
	if !isVariadic && fn.Type().NumIn() != (len(args)+argIndex) {
		return &result{err: NewEvaluationError(
			"function '%s' expects %d arguments, got %d",
			identifyCallee(expr),
			fn.Type().NumIn()-argIndex,
			len(args),
		)}
	}

	variadicIndex := fn.Type().NumIn() - 1
	for i, arg := range args {
		if isVariadic && i+argIndex >= variadicIndex {
			// Variadic argument
			varsType := fn.Type().In(variadicIndex)
			paramType := varsType.Elem()
			for _, a := range args[i:] {
				if !reflect.TypeOf(a).AssignableTo(paramType) {
					return &result{err: NewEvaluationError(
						"variadic argument '%v' is not assignable to type '%s'",
						arg,
						paramType.String(),
					)}
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
				return &result{err: NewEvaluationError(
					"argument '%v' is not assignable to parameter '%s'",
					arg,
					paramType.String(),
				)}
			}
		}

		in = append(in, argValue)
	}

	out := fn.Call(in)
	if len(out) == 2 {
		if out[1].Interface() != nil {
			return &result{err: out[1].Interface().(error)}
		}
	}
	return &result{value: out[0].Interface()}
}
func (i *Evaluator) visitArrayExpr(ctx context.Context, expr *Array) EvaluationResult {
	if ctx.Err() != nil {
		return &result{err: EvaluationCancelledErrror}
	}
	values := make([]interface{}, len(expr.values))
	for index, v := range expr.values {
		res := i.interpret(ctx, v)
		if res.Error() != nil {
			return res
		}
		values[index] = res.Get()
	}
	return &result{value: values}
}

func (i *Evaluator) visitMapExpr(ctx context.Context, expr *Map) EvaluationResult {
	if ctx.Err() != nil {
		return &result{err: EvaluationCancelledErrror}
	}
	m := make(map[string]interface{})
	for _, e := range expr.entries {
		var entry [2]interface{}
		res := i.interpret(ctx, e)
		if res.Error() != nil {
			return res
		}
		entry = res.Get().([2]interface{})
		key := entry[0]
		value := entry[1]
		m[fmt.Sprintf("%v", key)] = value
	}
	return &result{value: m}
}

func (i *Evaluator) visitMapEntryExpr(ctx context.Context, expr *MapEntry) EvaluationResult {
	if ctx.Err() != nil {
		return &result{err: EvaluationCancelledErrror}
	}
	keyRes := i.interpret(ctx, expr.key)
	if keyRes.Error() != nil {
		return keyRes
	}
	valueRes := i.interpret(ctx, expr.value)
	if valueRes.Error() != nil {
		return valueRes
	}
	return &result{value: [2]interface{}{keyRes.Get(), valueRes.Get()}}
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

func toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	default:
		return 0, fmt.Errorf("not a number")
	}
}

func isNumber(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	default:
		return false
	}
}

func areNumbers(values ...interface{}) bool {
	for _, v := range values {
		if !isNumber(v) {
			return false
		}
	}
	return true
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
		r := i.interpret(ctx, expr)
		once.Do(func() {
			result = r.Get()
			err = r.Error()
		})
	}()

	select {
	case <-done:
		return result, err
	case <-ctx.Done():
		return nil, NewEvaluationError("evaluation canceled: %s", ctx.Err().Error())
	}
}

func (i *Evaluator) interpret(ctx context.Context, expr Expr) EvaluationResult {
	if ctx.Err() != nil {
		return &result{err: EvaluationCancelledErrror}
	}
	return expr.Accept(ctx, i)
}
