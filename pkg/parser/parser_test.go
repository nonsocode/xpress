package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestASTPrinter(t *testing.T) {
	printer := ASTPrinter{}
	ast := NewParser("6 * (4 + 2)").Parse()
	res, err := printer.Print(ast)
	assert.Nil(t, err)
	assert.Equal(t, "(template 6 * (4 + 2))", res)
}

func TestEvaluator(t *testing.T) {
	evaluator := Evaluator{}
	ast := NewParser("{{6 * (4 + 2)}}").Parse()
	res, err := evaluator.Evaluate(ast)
	assert.Nil(t, err)
	assert.Equal(t, float64(36), res)
}
