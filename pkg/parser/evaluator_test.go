package parser

import (
	"context"
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/gval"
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
}

var cases = []SuccessCases{
	{template: "Just raw text", expect: "Just raw text"},
	{template: "{{ 123 * (45.67) }}", expect: float64(123 * 45.67)},
	{template: "{{-123 * (45.67) }} juxtaposed", expect: "-5617.41 juxtaposed"},
	{
		template: "{{-123 * (45.67) }} ", // converts to string if the template braces don't begin and end the string
		expect:   "-5617.41 ",
	},
	{
		template: `{{ 3 * 3 }} with text in-between {{ true ? "changed" : "not changed" }}`,
		expect:   "9 with text in-between changed",
	},
	{
		template: `{{'{{'}} 3 * 3 }} escaped template with template after {{ true ? "yes" : "no" }}`,
		expect:   "{{ 3 * 3 }} escaped template with template after yes",
	},
	{template: "{{ true ? 1 : 2 }}", expect: float64(1)},
	{template: "{{ false ? 1 : 2 }}", expect: float64(2)},
	{template: "{{ 4 * 4 }}", expect: float64(16)},
	{template: "{{ 4 + 4 }}", expect: float64(8)},
	{template: "{{ 4 + -4 }}", expect: float64(0)},
	{template: "{{ 10 - 4 }}", expect: float64(6)},
	{template: "{{ 8 / 4 }}", expect: float64(2)},
	{template: "{{ 5 > 4 }}", expect: true},
	{template: "{{ 3 > 4 }}", expect: false},
	{template: "{{ 5 < 4 }}", expect: false},
	{template: "{{ 4 < 5 }}", expect: true},
	{template: "{{ 5 == 5 }}", expect: true},
	{template: "{{ 5 == 4 }}", expect: false},
	{template: "{{ 5 != 5 }}", expect: false},
	{template: "{{ 5 != 6 }}", expect: true},
	{template: "{{ 5 >= 4 }}", expect: true},
	{template: "{{ 5 >= 5 }}", expect: true},
	{template: "{{ 5 >= 6 }}", expect: false},
	{template: "{{ 5 <= 4 }}", expect: false},
	{template: "{{ 5 <= 5 }}", expect: true},
	{template: "{{ math.abs(-5) }}", expect: float64(5)},
	{template: "{{ 'a' <= 'b' }}", expect: true},
	{template: "{{ 'b' <= 'b' }}", expect: true},
	{template: "{{ 'c' <= 'a' }}", expect: false},
	{template: "{{ 'a' >= 'b' }}", expect: false},
	{template: "{{ 'b' >= 'b' }}", expect: true},
	{template: "{{ 'c' >= 'a' }}", expect: true},
	{template: "{{ 'a' < 'b' }}", expect: true},
	{template: "{{ 'b' < 'b' }}", expect: false},
	{template: "{{ 'c' < 'b' }}", expect: false},
	{template: "{{ 'a' > 'b' }}", expect: false},
	{template: "{{ 'b' > 'b' }}", expect: false},
	{template: "{{ 'c' > 'b' }}", expect: true},
	{template: "{{ true }}", expect: true},
	{template: "{{ false }}", expect: false},
	{template: "{{ !true }}", expect: false},
	{template: "{{ !false }}", expect: true},
	{template: "{{ true && true}}", expect: true},
	{template: "{{ true && false}}", expect: false},
	{template: "{{ false && false}}", expect: false},
	{template: "{{ true || true}}", expect: true},
	{template: "{{ true || false}}", expect: true},
	{template: "{{ false || false}}", expect: false},
	{template: "{{ 4 > 5 && 5 == 5 }}", expect: false},
	{template: "{{ 4 > 5 || 5 == 5 }}", expect: true},
	{template: "{{ (4 > 5 || 5) }}", expect: true},
	{template: "{{ (4 > 5 && 5) }}", expect: false},
	{template: "{{ (4 > 5 || 5) == true}}", expect: true},
	{template: "{{ (4 > 5 && 5) == true}}", expect: false},
	{template: "{{ true && true && true}}", expect: true},
	{template: "{{ true && false && true}}", expect: false},
	{template: "{{ true && false && true}}", expect: false},
	{template: `{{ "a string" == "a string"}}`, expect: true},
	{template: `{{ "a string" != "a different string"}}`, expect: true},
	{template: `{{ "a string" == "a different string"}}`, expect: false},
	{template: `{{ "a string" != "a string"}}`, expect: false},
	{template: `{{ "a string" != "a string"}}`, expect: false},
	{template: `{{[1, 2, true, "a"]}}`, expect: []interface{}{float64(1), float64(2), true, "a"}},
	{template: `{{[1, 2, true, "a"]}} `, expect: "[1 2 true a] "},
	{template: `{{ "a string" + " " + "Joined" }}`, expect: "a string Joined"},
	{template: `{{ concat("string", "joined by", "another") }}`, expect: "stringjoined byanother"},
	{template: `{{ concat("string", " ", concat("with another", concat(" ", "recursive"))) }}`, expect: "string with another recursive"},
	{template: "{{ getDeepObject().deep.object.with.values }}", expect: []interface{}{3, 2, 1}},
	{template: "{{ funcable()('host') }}", expect: "a function with host"},
	{template: "{{ someObject.key }}", expect: "value"},
	{template: "{{ someObject['key'] }}", expect: "value"},
	{template: "{{ someObject.nested.key1 }}", expect: "value2"},
	{template: "{{ someObject['nested'].key1 }}", expect: "value2"},
	{template: "{{ someObject.nested['key1'] }}", expect: "value2"},
	{template: "{{ someObject['nested']['key1'] }}", expect: "value2"},
	{template: "{{ getDeepObject().deep.object.with.values[0] }}", expect: int(3)},
}

var errorCases = []ErrorCases{
	{"{{ 5 > }}", "parse error: Error at position 7. Expect expression. got }}"},
	{
		"{{ 5 ",
		"parse error: Error at position 5. Expect '}}' after expression. got unclosed action",
	},
	{"{{ 5 6 }}", "parse error: Error at position 5. Expect '}}' after expression. got 6"},
	{
		"{{ nonexistentFunction() }}",
		"cannot call non-function 'nonexistentFunction' of type <nil>",
	},
	{
		"{{ waitMs(10) }}",
		"evaluation timed out after 5ms",
	},
	{
		"{{ longLoopWithContext() }}",
		"evaluation timed out after 5ms",
	},
}

func TestExampleParser(t *testing.T) {
	evaluator := NewInterpreter()
	evaluator.AddMembers(createTestTemplateFunctions())
	evaluator.SetTimeout(5 * time.Millisecond)
	test := func(cas *SuccessCases) {
		ast := NewParser(cas.template).Parse()
		res, err := evaluator.Evaluate(ast)
		assert.Nil(t, err)
		assert.Equal(t, cas.expect, res)
	}
	for _, c := range cases {
		if c.only {
			test(&c)
			return
		}
	}

	for _, c := range cases {
		test(&c)
	}

}

func TestExampleParserErrors(t *testing.T) {
	evaluator := NewInterpreter()
	evaluator.SetTimeout(5 * time.Millisecond)
	evaluator.AddMembers(createTestTemplateFunctions())
	for _, c := range errorCases {
		ast := NewParser(c.template).Parse()
		_, err := evaluator.Evaluate(ast)
		assert.NotNil(t, err)
		assert.Equal(t, c.msg, err.Error())
	}
}

func BenchmarkComplexParser(b *testing.B) {
	// create a parser with complex expression
	for n := 0; n < b.N; n++ {
		NewParser("{{ concat('string', ' ', concat('with another', concat(' ', 'recursive'))) }} and some advanced math {{ math.pow(2, 3) + 56 / 4 * 6 }} and some object access {{ someObject.nested.key1 }} with other function calls {{ getDeepObject().deep.object.with.values }}").Parse()
	}
}

func BenchmarkSimpleParser(b *testing.B) {
	// create a parser with simple expression
	for n := 0; n < b.N; n++ {
		NewParser("{{ 54 * (6 / 2) }}").Parse()
	}
}
func BenchmarkEvaluator(b *testing.B) {
	evaluator := NewInterpreter()
	evaluator.AddMembers(createTestTemplateFunctions())
	evaluator.SetTimeout(5 * time.Millisecond)
	// create a parser with complex expression
	ast := NewParser("{{ concat('string', ' ', concat('with another', concat(' ', 'recursive'))) }} and some advanced math {{ math.pow(2, 3) + 56 / 4 * 6 }} and some object access {{ someObject.nested.key1 }} with other function calls {{getDeepObject().deep.object.with.values}}").Parse()
	for n := 0; n < b.N; n++ {
		evaluator.Evaluate(ast)
	}
}

func BenchmarkGvalComplex(b *testing.B) {
	ogFunctions := createTestTemplateFunctions()
	vars := map[string]interface{}{
		"concat": func(args ...string) string {
			builder := strings.Builder{}
			for _, arg := range args {
				builder.WriteString(arg)
			}
			return builder.String()
		},
		"math": ogFunctions["math"],
		"deepObject": func() map[string]interface{} {
			return map[string]interface{}{
				"deep": map[string]interface{}{
					"object": map[string]interface{}{
						"with": map[string]interface{}{
							"values": []interface{}{3, 2, 1},
						},
					},
				},
			}
		},
		"someObject": ogFunctions["someObject"],
	}

	for n := 0; n < b.N; n++ {
		gval.Evaluate(`concat("string", " ", concat("with another", concat(" ", "recursive"))) + " and some advanced math " + math.pow(2, 3) + 56 / 4 * 6 + " and some object access " + someObject.nested.key1 + " with other function calls " + deepObject().deep.object.with.values`, vars)
	}

}
func TestGvalComplex(t *testing.T) {
	ogFunctions := createTestTemplateFunctions()
	vars := map[string]interface{}{
		"concat": func(args ...string) string {
			builder := strings.Builder{}
			for _, arg := range args {
				builder.WriteString(arg)
			}
			return builder.String()
		},
		"math": ogFunctions["math"],
		"getDeepObject": map[string]interface{}{
			"deep": map[string]interface{}{
				"object": map[string]interface{}{
					"with": map[string]interface{}{
						"values": []interface{}{3, 2, 1},
					},
				},
			},
		},
		"someObject": ogFunctions["someObject"],
	}

	_, err := gval.Evaluate(`concat("string", " ", concat("with another", concat(" ", "recursive"))) + " and some advanced math " + math.pow(2, 3) + 56 / 4 * 6 + " and some object access " + someObject.nested.key1 + " with other function calls " + getDeepObject.deep.object.with.values`, vars)
	assert.Nil(t, err)
}

func createTestTemplateFunctions() map[string]interface{} {

	return map[string]interface{}{
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
				fmt.Println(count, "longLoopWithContext")
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
		"funcable": func() any {
			return func(ctx context.Context, stuff string) string {
				return "a function with " + stuff
			}
		},
		"someObject": map[string]interface{}{
			"key": "value",
			"nested": map[string]interface{}{
				"key1": "value2",
			},
		},
	}
}
