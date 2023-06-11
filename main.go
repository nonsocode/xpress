package main

import (
	"fmt"

	"github.com/nonsocode/artihm/pkg/parser"
)

func main() {

	evaluator := parser.Evaluator{}
	template1 := "hehe{{-123 * (45.67) }}"
	template2 := `{{ 3 * 3 }} withot spaces is {{ true ? "changed" : "not changed" }}`
	fmt.Println(parser.NewParser("-hello world").Parse(&evaluator))
	fmt.Println(parser.NewParser(template1).Parse(&evaluator))
	fmt.Println(parser.NewParser(template2).Parse(&evaluator))
	// ...
}
