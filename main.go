package main

import (
	"fmt"

	"github.com/nonsocode/artihm/pkg/parser"
)

func main() {

	regularPrinter := parser.ASTPrinter{}
	evaluator := parser.Evaluator{}
	fmt.Println(parser.NewParser("-123  (45.67)").Parse(&evaluator))
	fmt.Println(parser.NewParser("-123  (45.67)").Parse(&regularPrinter))
	fmt.Println(parser.NewParser("6 * (4 + 2)").Parse(&evaluator))
	// ...
}
