package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestASTPrinter(t *testing.T) {
	printer := ASTPrinter{}
	res, err := NewParser("6 * (4 + 2)").Parse(&printer)
	assert.Nil(t, err)
	assert.Equal(t, "(template 6 * (4 + 2))", res)
}

func TestEvaluator(t *testing.T) {
	evaluator := Evaluator{}
	res, err := NewParser("{{6 * (4 + 2)}}").Parse(&evaluator)
	assert.Nil(t, err)
	assert.Equal(t, float64(36), res)
}
