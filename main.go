package main

import (
	"fmt"

	"github.com/nonsocode/artihm/pkg/parser"
)

func main() {
	var zero interface{} = 0
	fmt.Println(nil == zero)
	fmt.Println("hello")

	evaluator := *parser.NewInterpreter()
	template1 := "hehe{{-123 * (45.67) }}"
	template2 := `{{ 3 * 3 }} withot spaces is {{ true ? "changed" : "not changed" }}`
	template3 := `This year's woman of the year is {{ hello().woman.funcs().justice }}`
	evaluator.AddFunc("hello", func(args ...interface{}) (interface{}, error) {
		return map[string]interface{}{
			"woman": map[string]interface{}{
				"of": []interface{}{
					2020, 2021, 2022,
				},
				"funcs": func(args ...interface{}) (interface{}, error) {
					type HMap struct {
						justice string
					}
					return HMap{
						justice: "Ruth Bader Ginsburg",
					}, nil
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
