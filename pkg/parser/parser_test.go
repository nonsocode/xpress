package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestASTPrinter(t *testing.T) {
	printer := ASTPrinter{}
	assert.Equal(t, "(* 6 (group (+ 4 2)))", NewParser("6 * (4 + 2)").Parse(&printer))
}

func TestEvaluator(t *testing.T) {
	evaluator := Evaluator{}
	assert.Equal(t, float64(36), NewParser("6 * (4 + 2)").Parse(&evaluator))
}
