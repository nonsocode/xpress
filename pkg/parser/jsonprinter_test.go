package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestASTPrinter(t *testing.T) {
	printer := NewJSONPrinter()
	ast := NewParser("6 * (4 + 2)").Parse()
	res, err := printer.Print(ast)
	assert.Nil(t, err)
	printed := `{
  "type": "Template",
  "exprs": [
    {
      "type": "Literal",
      "value": "6 * (4 + 2)",
      "raw": "6 * (4 + 2)"
    }
  ]
}`

	assert.Equal(t, printed, res)
}
