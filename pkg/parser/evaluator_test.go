package parser

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type SuccessCases struct {
	template string
	expect   interface{}
}

type ErrorCases struct {
	template string
	msg      string
}

var cases = []SuccessCases{
	{"Just raw text", "Just raw text"},
	{"{{ 123 * (45.67) }}", float64(123 * 45.67)},
	{"{{-123 * (45.67) }} juxtaposed", "-5617.41 juxtaposed"},
	{
		"{{-123 * (45.67) }} ",
		"-5617.41 ",
	}, // converts to string if the template braces don't begin and end the string
	{
		`{{ 3 * 3 }} with text in-between {{ true ? "changed" : "not changed" }}`,
		"9 with text in-between changed",
	},
	{
		`{{'{{'}} 3 * 3 }} escaped template with template after {{ true ? "yes" : "no" }}`,
		"{{ 3 * 3 }} escaped template with template after yes",
	},
	{"{{ 4 * 4 }}", float64(16)},
	{"{{ 8 / 4 }}", float64(2)},
	{"{{ 5 > 4 }}", true},
	{"{{ 5 < 4 }}", false},
	{"{{ 5 == 5 }}", true},
	{"{{ 5 != 5 }}", false},
	{"{{ 5 >= 5 }}", true},
	{"{{ 5 <= 5 }}", true},
	{"{{ 5 >= 4 }}", true},
	{"{{ 5 <= 4 }}", false},
	{"{{ true && true}}", true},
	{"{{ true && false}}", false},
	{"{{ false && false}}", false},
	{"{{ true || true}}", true},
	{"{{ true || false}}", true},
	{"{{ false || false}}", false},
	{"{{ 4 > 5 && 5 == 5 }}", false},
	{"{{ 4 > 5 || 5 == 5 }}", true},
	{"{{ (4 > 5 || 5) }}", true},
	{"{{ (4 > 5 && 5) }}", false},
	{"{{ (4 > 5 || 5) == true}}", true},
	{"{{ (4 > 5 && 5) == true}}", false},
	{"{{ true && true && true}}", true},
	{"{{ true && false && true}}", false},
	{"{{ true && false && true}}", false},
	{`{{ "a string" == "a string"}}`, true},
	{`{{ "a string" != "a different string"}}`, true},
	{`{{ "a string" == "a different string"}}`, false},
	{`{{ "a string" != "a string"}}`, false},
	{`{{ "a string" != "a string"}}`, false},
	{`{{[1, 2, true, "a"]}}`, []interface{}{float64(1), float64(2), true, "a"}},
	{`{{[1, 2, true, "a"]}} `, "[1 2 true a] "},
	{`{{ "a string" + " " + "Joined" }}`, "a string Joined"},
	{`{{ concat("string", "joined by", "another") }}`, "stringjoined byanother"},
	{
		`{{ concat("string", " ", concat("with another", concat(" ", "recursive"))) }}`,
		"string with another recursive",
	},
	{
		"{{ getDeepObject().deep.object.with.values }}",
		[]interface{}{3, 2, 1},
	},
	{"{{ funcable()() }}", "a function"},
}

var errorCases = []ErrorCases{
	{"{{ 5 > }}", "parse error: Error at position 7. Expect expression. got }}"},
	{"{{ 5 ", "parse error: Error at position 5. Expect '}}' after expression. got unclosed action"},
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
	evaluator.SetFunctions(createTestTemplateFunctions())
	evaluator.SetTimeout(5 * time.Millisecond)
	for _, c := range cases {
		ast := NewParser(c.template).Parse()
		res, err := evaluator.Evaluate(ast)
		assert.Nil(t, err)
		assert.Equal(t, c.expect, res)
	}

}

func TestExampleParserErrors(t *testing.T) {
	evaluator := NewInterpreter()
	evaluator.SetTimeout(5 * time.Millisecond)
	evaluator.SetFunctions(createTestTemplateFunctions())
	for _, c := range errorCases {
		ast := NewParser(c.template).Parse()
		_, err := evaluator.Evaluate(ast)
		assert.NotNil(t, err)
		assert.Equal(t, c.msg, err.Error())
	}
}

func createTestTemplateFunctions() map[string]interface{} {

	return map[string]interface{}{
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
		"funcable": func() func() string {
			return func() string {
				return "a function"
			}
		},
	}
}
