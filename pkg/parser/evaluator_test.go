package parser

import (
	"context"
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type SuccessCases struct {
	template string
	expect   interface{}
	only     bool
}

type ErrorCases struct {
	template string
	msg      string
	only     bool
}

type Dummy struct {
	chassy  string
	Exposed string
}

var cases = []SuccessCases{
	{template: "Just raw text", expect: "Just raw text"},
	{template: "@{{ 123 * (45.67) }}", expect: float64(123 * 45.67)},
	{template: "@{{ 1 + 2 * 3 }}", expect: float64(7)},
	{template: "@{{ 2 * 3 + 1 }}", expect: float64(7)},
	{template: "@{{-123 * (45.67) }} juxtaposed", expect: "-5617.41 juxtaposed"},
	{template: "@{{-123 * (45.67) }} ", expect: "-5617.41 "}, // converts to string if the template braces don't begin and end the string
	{template: `@{{ 3 * 3 }} with text in-between @{{ true ? "changed" : "not changed" }}`, expect: "9 with text in-between changed"},
	{template: `@{{ '@{{' }} 3 * 3 }} escaped template with template after @{{ true ? "yes" : "no" }}`, expect: "@{{ 3 * 3 }} escaped template with template after yes"},
	{template: "@{{ true ? 1 : 2 }}", expect: float64(1)},
	{template: "@{{ false ? 1 : 2 }}", expect: float64(2)},
	{template: "@{{ 4 * 4 }}", expect: float64(16)},
	{template: "@{{ 4 + 4 }}", expect: float64(8)},
	{template: "@{{ 4 + -4 }}", expect: float64(0)},
	{template: "@{{ 10 - 4 }}", expect: float64(6)},
	{template: "@{{ 8 / 4 }}", expect: float64(2)},
	{template: "@{{ 5 > 4 }}", expect: true},
	{template: "@{{ 3 > 4 }}", expect: false},
	{template: "@{{ 5 < 4 }}", expect: false},
	{template: "@{{ 4 < 5 }}", expect: true},
	{template: "@{{ 5 == 5 }}", expect: true},
	{template: "@{{ 5 == 4 }}", expect: false},
	{template: "@{{ 5 != 5 }}", expect: false},
	{template: "@{{ 5 != 6 }}", expect: true},
	{template: "@{{ 5 >= 4 }}", expect: true},
	{template: "@{{ 5 >= 5 }}", expect: true},
	{template: "@{{ 5 >= 6 }}", expect: false},
	{template: "@{{ 5 <= 4 }}", expect: false},
	{template: "@{{ 5 <= 5 }}", expect: true},
	{template: "@{{ math.abs(-5) }}", expect: float64(5)},
	{template: "@{{ 'a' <= 'b' }}", expect: true},
	{template: "@{{ 'b' <= 'b' }}", expect: true},
	{template: "@{{ 'c' <= 'a' }}", expect: false},
	{template: "@{{ 'a' >= 'b' }}", expect: false},
	{template: "@{{ 'b' >= 'b' }}", expect: true},
	{template: "@{{ 'c' >= 'a' }}", expect: true},
	{template: "@{{ 'a' < 'b' }}", expect: true},
	{template: "@{{ 'b' < 'b' }}", expect: false},
	{template: "@{{ 'c' < 'b' }}", expect: false},
	{template: "@{{ 'a' > 'b' }}", expect: false},
	{template: "@{{ 'b' > 'b' }}", expect: false},
	{template: "@{{ 'c' > 'b' }}", expect: true},
	{template: "@{{ true }}", expect: true},
	{template: "@{{ false }}", expect: false},
	{template: "@{{ !true }}", expect: false},
	{template: "@{{ !false }}", expect: true},
	{template: "@{{ true && true}}", expect: true},
	{template: "@{{ true && false}}", expect: false},
	{template: "@{{ false && false}}", expect: false},
	{template: "@{{ true || true}}", expect: true},
	{template: "@{{ true || false}}", expect: true},
	{template: "@{{ false || false}}", expect: false},
	{template: "@{{ 4 > 5 && 5 == 5 }}", expect: false},
	{template: "@{{ 4 > 5 || 5 == 5 }}", expect: true},
	{template: "@{{ (4 > 5 || 5) }}", expect: true},
	{template: "@{{ (4 > 5 && 5) }}", expect: false},
	{template: "@{{ (4 > 5 || 5) == true}}", expect: true},
	{template: "@{{ (4 > 5 && 5) == true}}", expect: false},
	{template: "@{{ true && true && true}}", expect: true},
	{template: "@{{ true && false && true}}", expect: false},
	{template: "@{{ true && false && true}}", expect: false},
	{template: `@{{ "a string" == "a string"}}`, expect: true},
	{template: `@{{ "a string" != "a different string"}}`, expect: true},
	{template: `@{{ "a string" == "a different string"}}`, expect: false},
	{template: `@{{ "a string" != "a string"}}`, expect: false},
	{template: `@{{ "a string" != "a string"}}`, expect: false},
	{template: `@{{[1, 2, true, "a"]}}`, expect: []interface{}{float64(1), float64(2), true, "a"}},
	{template: `@{{[1, 2, true, "a"]}} `, expect: "[1 2 true a] "},
	{template: `@{{ "a string" + " " + "Joined" }}`, expect: "a string Joined"},
	{template: `@{{ concat("string", "joined by", "another") }}`, expect: "stringjoined byanother"},
	{template: `@{{ concat("string", " ", concat("with another", concat(" ", "nested"))) }}`, expect: "string with another nested"},
	{template: "@{{ getDeepObject().deep.object.with.values }}", expect: []interface{}{3, 2, 1}},
	{template: "@{{ someFunc()('host') }}", expect: "a function with host"},
	{template: "@{{ someObject.key }}", expect: "value"},
	{template: `@{{ "somestring".length }}`, expect: float64(10)},
	{template: `@{{ "somestring".length * 2 }}`, expect: float64(20)},
	{template: "@{{ someObject['key'] }}", expect: "value"},
	{template: "@{{ someObject.nested.key1 }}", expect: "value2"},
	{template: "@{{ someObject.nested.struct.Key }}", expect: "StructValue"},
	{template: "@{{ pointerDummy.PointerReceiverMethod() }}", expect: "pointer value"},
	{template: "@{{ pointerDummy.StructReceiverMethod() }}", expect: "struct value"},
	{template: "@{{ dummy.PointerReceiverMethod() }}", expect: "pointer value"},
	{template: "@{{ dummy.StructReceiverMethod() }}", expect: "struct value"},
	{template: `@{{ pointerDummy["PointerReceiverMethod"]() }}`, expect: "pointer value"},
	{template: `@{{ pointerDummy["StructReceiverMethod"]() }}`, expect: "struct value"},
	{template: `@{{ dummy["PointerReceiverMethod"]() }}`, expect: "pointer value"},
	{template: `@{{ dummy["StructReceiverMethod"]() }}`, expect: "struct value"},
	{template: `@{{ someObject["nested"].key1 }}`, expect: "value2"},
	{template: "@{{ someObject.nested['key1'] }}", expect: "value2"},
	{template: "@{{ someObject['nested']['key1'] }}", expect: "value2"},
	{template: "@{{ getDeepObject().deep.object.with.values[0] }}", expect: int(3)},
	{template: "@{{ getDeepObject().nonexistent }}", expect: nil},
	{template: "@{{ getDeepObject()['nonexistent'] }}", expect: nil},
	{template: "@{{ someObject.nonexistent?.deply['nonexistent'].path }}", expect: nil},
	{template: "@{{ getDeepObject().nonexistent?.deply.nonexistent.path }}", expect: nil},
	{template: "@{{ getDeepObject().nonexistentFunc?.() }}", expect: nil},
	{template: "@{{ getDeepObject?.()?.deep?.object?.with?.values?.[0] }}", expect: int(3)},
	{template: "@{{ [1,2,3][1] }}", expect: float64(2)},
	{template: `@{{ [
		1,
		2,
		3
	][1] }}`, expect: float64(2)}, // multiline
	{template: "@{{ [[0,0], [0,1], [1,2]] }}", expect: []interface{}{[]interface{}{float64(0), float64(0)}, []interface{}{float64(0), float64(1)}, []interface{}{float64(1), float64(2)}}},
	{template: `@{{ 
		{
			a: concat("Hello ", "World"), 
		  "nested": {
				b: {c: {d: "this is d"}},
				someArray: [1, 2, "hello".length + 1],
				"stringKey": 87,
				[someObject.nested.key1]: "this is a value",
				[getDeepObject().deep.object.with.values[0]]: "this is another value"
		  }
		} 
  }}`, expect: map[string]interface{}{
		"a": "Hello World",
		"nested": map[string]interface{}{
			"b": map[string]interface{}{
				"c": map[string]interface{}{
					"d": "this is d",
				},
			},
			"someArray": []interface{}{float64(1), float64(2), float64(6)},
			"stringKey": float64(87),
			"value2":    "this is a value",
			"3":         "this is another value",
		},
	}},
}

var errorCases = []ErrorCases{
	{template: "@{{ 5 > }}", msg: "parse error: Error at position 8. Expect expression. got }}"},
	{template: "@{{ 5 ", msg: "parse error: Error at position 6. Expect '}}' after expression. got unclosed action"},
	{template: "@{{ 5 6 }}", msg: "parse error: Error at position 6. Expect '}}' after expression. got 6"},
	{template: "@{{ nonexistentFunction() }}", msg: "cannot call non-function 'nonexistentFunction' of type <nil>"},
	{template: "@{{ waitMs(10) }}", msg: "evaluation timed out after 5ms"},
	{template: "@{{ longLoopWithContext() }}", msg: "evaluation timed out after 5ms"},
	{template: "@{{ getDeepObject().deep.object.with.values[3] }}", msg: "index '3' is out of bounds"},
	{template: "@{{ getDeepObject().deep.object.with.values[-1] }}", msg: "index '-1' is out of bounds"},
	{template: "@{{ getDeepObject().nonexistent.key }}", msg: "cannot get property 'key' of nil"},
	{template: "@{{ [1,2,3,4 }}", msg: "parse error: Error at position 13. Expect ']' after array expression. got }"},
}

func (d *Dummy) PointerReceiverMethod() string {
	return "pointer value"
}

func (d Dummy) StructReceiverMethod() string {
	return "struct value"
}

func TestExampleParser(t *testing.T) {
	evaluator := NewInterpreter()
	evaluator.SetMembers(createTestTemplateFunctions())
	evaluator.SetTimeout(5 * time.Hour)
	test := func(t *testing.T, cas *SuccessCases) {
		ast := NewParser(cas.template).Parse()
		res, err := evaluator.Evaluate(context.TODO(), ast)
		assert.Nil(t, err)
		assert.Equal(t, cas.expect, res)
	}
	for _, c := range cases {
		if c.only {
			t.Run(c.template, func(t *testing.T) {
				test(t, &c)
			})
			return
		}
	}

	for _, c := range cases {
		t.Run(c.template, func(t *testing.T) {
			test(t, &c)
		})
	}

}

func TestExampleParserErrors(t *testing.T) {
	evaluator := NewInterpreter()
	evaluator.SetTimeout(5 * time.Millisecond)
	evaluator.SetMembers(createTestTemplateFunctions())
	test := func(t *testing.T, cas *ErrorCases) {
		ast := NewParser(cas.template).Parse()
		val, err := evaluator.Evaluate(context.TODO(), ast)
		assert.Nil(t, val)
		assert.NotNil(t, err)
		assert.Equal(t, cas.msg, err.Error())
	}
	for _, c := range errorCases {
		if c.only {
			t.Run(c.template, func(t *testing.T) {
				test(t, &c)
			})
			return
		}
	}
	for _, c := range errorCases {
		t.Run(c.template, func(t *testing.T) {
			test(t, &c)
		})
	}
}

func BenchmarkComplexParser(b *testing.B) {
	// create a parser with complex expression
	for n := 0; n < b.N; n++ {
		NewParser("@{{ concat('string', ' ', concat('with another', concat(' ', 'recursive'))) }} and some advanced math @{{ math.pow(2, 3) + 56 / 4 * 6 }} and some object access @{{ someObject.nested.key1 }} with other function calls @{{ getDeepObject().deep.object.with.values }}").Parse()
	}
}

func BenchmarkSimpleParser(b *testing.B) {
	// create a parser with simple expression
	for n := 0; n < b.N; n++ {
		NewParser("@{{ 54 * (6 / 2) }}").Parse()
	}
}
func BenchmarkEvaluator(b *testing.B) {
	template := `@{{ concat('string', ' ', concat('with another', concat(' ', 'recursive'))) }} and some advanced math @{{ math.pow(2, 3) + 56 / 4 * 6 }} and some object access @{{ someObject.nested.key1 }} with other function calls @{{getDeepObject().deep.object.with.values}}`
	evaluator := NewInterpreter()
	evaluator.SetMembers(createTestTemplateFunctions())
	evaluator.SetTimeout(5 * time.Millisecond)
	// create a parser with complex expression
	ast := NewParser(template).Parse()
	for n := 0; n < b.N; n++ {
		evaluator.Evaluate(context.TODO(), ast)
	}
}

func createTestTemplateFunctions() map[string]interface{} {
	return map[string]interface{}{
		"pointerDummy": &Dummy{},
		"dummy":        Dummy{},
		"math": map[string]interface{}{
			"abs":   math.Abs,
			"acos":  math.Acos,
			"asin":  math.Asin,
			"atan":  math.Atan,
			"atan2": math.Atan2,
			"ceil":  math.Ceil,
			"cos":   math.Cos,
			"exp":   math.Exp,
			"floor": math.Floor,
			"log":   math.Log,
			"log10": math.Log10,
			"max":   math.Max,
			"min":   math.Min,
			"pow":   math.Pow,
			"sin":   math.Sin,
			"sqrt":  math.Sqrt,
			"tan":   math.Tan,
		},
		"concat": func(ctx context.Context, args ...string) (interface{}, error) {
			builder := strings.Builder{}
			for _, arg := range args {
				builder.WriteString(arg)
			}
			return builder.String(), nil
		},
		"getDeepObject": func(ctx context.Context, args ...interface{}) (interface{}, error) {
			return map[string]interface{}{
				"deep": map[string]interface{}{
					"object": map[string]interface{}{
						"with": map[string]interface{}{
							"values": []interface{}{3, 2, 1},
						},
					},
				},
			}, nil
		},
		"deepObject": map[string]interface{}{
			"deep": map[string]interface{}{
				"object": map[string]interface{}{
					"with": map[string]interface{}{
						"values": []interface{}{3, 2, 1},
					},
				},
			},
		},
		"waitMs": func(msec int) bool {
			time.Sleep(time.Duration(msec) * time.Millisecond)
			return true
		},
		"longLoop": func() (bool, error) {
			var count int
			for {
				count++
				fmt.Println(count, "longLoop")
				time.Sleep(1 * time.Second)
				if count > 10 {
					break
				}
			}

			return true, fmt.Errorf("should not get here")
		},
		"longLoopWithContext": func(ctx context.Context) (bool, error) {
			var count int
			for {
				count++
				time.Sleep(1 * time.Second)
				if count > 2 {
					break
				}
				select {
				case <-ctx.Done():
					return false, ctx.Err()
				default:
				}
			}
			return true, nil
		},
		"someFunc": func() any {
			return func(ctx context.Context, stuff string) string {
				return "a function with " + stuff
			}
		},
		"someObject": map[string]interface{}{
			"key": "value",
			"nested": map[string]interface{}{
				"key1": "value2",
				"struct": struct{ Key string }{
					Key: "StructValue",
				},
			},
		},
	}
}
