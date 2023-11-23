// xpress is a simple and fast template engine for Go.
//
// # Basic Usage
//
//	package main
//
//	import (
//		"fmt"
//		"github.com/nonsocode/xpress/pkg/parser"
//	)
//
//	func main() {
//		// create a parser
//		p := parser.NewParser("Hello, @{{ name }}!")
//
//		// parse the template
//		ast := p.Parse()
//
//		// create an evaluator
//		e := parser.NewInterpreter()
//
//		// set members
//		e.SetMembers(map[string]interface{}{
//			"name": "John",
//		})
//
//		// evaluate the template
//		res, err := e.Evaluate(ast)
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println(res)	// Hello, John!
//	}
//
// You can reuse the evaluator for multiple templates:
//
//	func main() {
//		// create an evaluator with members
//		e := parser.NewInterpreter()
//		e.SetMembers(map[string]interface{}{
//			"name": "John",
//			"concat": func(a, b string) string {
//				return a + b
//			},
//		})
//
//		// create an AST
//		ast := parser.NewParser("Hello, @{{ name }}!").Parse()
//
//		// evaluate the template
//		res, err := e.Evaluate(ast)
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println(res)	// Hello, John!
//
//		// create another AST
//		p = parser.NewParser("Hello, @{{ concat(name, ' Doe') }}!").Parse()
//
//		// evaluate the template
//		res, err = e.Evaluate(ast)
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println(res)	// Hello, John Doe!
//	}
//
// # Template Syntax
//
// The template syntax is loosely inspired by javascript syntax. When a template is evaluated,
// If the template begins and ends exactly with the template delimiters, the template will be
// return the exact result of the operation regardless of the type of the result. Otherwise,
// the result will be converted to a string.
//
//	// "@{{ 1 + 2 }}" will return number `3`
//	// "@{{ 1 + 2 }} "  will return string `3 `.
//
// # Members
//
// Members are variables and functions that can be accessed from the template. When defining
// functions, You can optionally define the first argument as context.Context. A context
// will be passed to the function when it is called. Functions must return either 1 or 2 values.
// If the function returns 1 value, it will be used as the result of the function call. If the
// function returns 2 values, the first value will be used as the result of the function call
// and the second value will be used as an error. If the error is not nil, the evaluation will
// be aborted and the error will be returned.
//
// # Benchmarks
//
//	goos: linux
//	goarch: amd64
//	pkg: github.com/nonsocode/xpress/pkg/parser
//	cpu: AMD EPYC 7763 64-Core Processor
//	BenchmarkComplexParser-2          126980              9935 ns/op            7648 B/op         70 allocs/op
//	BenchmarkSimpleParser-2           630464              2714 ns/op            1696 B/op         18 allocs/op
//	BenchmarkEvaluator-2               34920             36895 ns/op            4440 B/op         97 allocs/op
package xpress
