package parser

import (
	"testing"

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
	{"{{-123 * (45.67) }} ", "-5617.41 "}, // converts to string if the template braces don't begin and end the string
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
}

var errorCases = []ErrorCases{
	{"{{ 5 > }}", "unexpected end of expression"},
}

func TestExampleParser(t *testing.T) {
	evaluator := NewInterpreter()
	for _, c := range cases {
		res, err := NewParser(c.template).Parse(evaluator)
		assert.Nil(t, err)
		assert.Equal(t, c.expect, res)
	}
}

func TestExampleParserErrors(t *testing.T) {
	evaluator := NewInterpreter()
	for _, c := range errorCases {
		_, err := NewParser(c.template).Parse(evaluator)
		assert.NotNil(t, err)
		assert.Equal(t, c.msg, err.Error())
	}
}
