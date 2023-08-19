package main

import (
	"fmt"

	"github.com/nonsocode/artihm/pkg/parser"
)

func main() {
	evaluator := *parser.NewInterpreter()
	// jsonPrinter := parser.NewJSONPrinter()
	template1 := parser.NewParser("{{ 3 * 6 + (5 + 6) / 5 }}").Parse()
	template2 := parser.NewParser(`{{ 3 * 3 }} withot spaces is {{ true ? "changed" : "not changed" }}`).
		Parse()
	template3 := parser.NewParser(`This year's woman of the year is {{ hello("df").woman.funcs().justice }}`).
		Parse()
	evaluator.AddMember("hello", func(someting string) (interface{}, error) {
		// panic("not implemented")
		return map[string]interface{}{
			"woman": map[string]interface{}{
				"of": []interface{}{
					2020, 2021, 2022,
				},
				"funcs": func() (interface{}, error) {
					return map[string]interface{}{
						"justice": "Ruth Bader Ginsburg",
					}, nil
				},
			},
		}, nil
	})
	fmt.Println(evaluator.Evaluate(template1))
	fmt.Println(evaluator.Evaluate(template2))
	fmt.Println(evaluator.Evaluate(template3))
	// fmt.Println(jsonPrinter.Print(template1))
	// fmt.Println(jsonPrinter.Print(template2))
	// fmt.Println(jsonPrinter.Print(template3))
}
