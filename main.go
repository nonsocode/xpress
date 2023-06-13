package main

import (
	"fmt"

	"github.com/nonsocode/artihm/pkg/parser"
)

func main() {

	evaluator := *parser.NewInterpreter()
	template1 := "hehe{{-123 * (45.67) }}"
	template2 := `{{ 3 * 3 }} withot spaces is {{ true ? "changed" : "not changed" }}`
	template3 := `This year's woman of the year is {{ hello().woman.of.the.year }}`
	evaluator.AddFunc("hello", func(args ...interface{}) (interface{}, error) {
		return map[string]interface{}{
			"woman": map[string]interface{}{
				"of": map[string]interface{}{
					"the": map[string]interface{}{
						"year": "2020",
					},
				},
			},
		}, nil
	})
	fmt.Println(parser.NewParser("-hello world").Parse(&evaluator))
	fmt.Println(parser.NewParser(template1).Parse(&evaluator))
	fmt.Println(parser.NewParser(template2).Parse(&evaluator))
	fmt.Println(parser.NewParser(template3).Parse(&evaluator))
	// ...
}
