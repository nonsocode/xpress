package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexerCanLexText(t *testing.T) {
	lex := NewLexer("Hello World")
	lex.run()
	assert.Equal(t, []Token{
		{"Hello World", TEXT, 0, 1},
		{"", EOF, 11, 1},
	}, lex.tokens)
}

func TestLexerCanLexTextWithNewlines(t *testing.T) {
	lex := NewLexer("Hello\nWorld")
	lex.run()
	assert.Equal(t, []Token{
		{"Hello\nWorld", TEXT, 0, 1},
		{"", EOF, 11, 2},
	}, lex.tokens)
}

func TestLexerCanLexAction(t *testing.T) {
	lex := NewLexer(`{{"Hello"}}`)
	lex.run()
	assert.Equal(
		t,
		[]Token{
			{lexeme: "{{", tokenType: TEMPLATE_LEFT_BRACE, start: 0, line: 1},
			{lexeme: "\"Hello\"", tokenType: STRING, start: 2, line: 1},
			{lexeme: "}}", tokenType: TEMPLATE_RIGHT_BRACE, start: 9, line: 1},
			{lexeme: "", tokenType: EOF, start: 11, line: 1},
		},
		lex.tokens,
	)
}

func TestBasicArithmetic(t *testing.T) {
	lex := NewLexer("{{ 3 * 3 }}")
	lex.run()
	assert.Equal(
		t,
		[]Token{
			{lexeme: "{{", tokenType: TEMPLATE_LEFT_BRACE, start: 0, line: 1},
			{lexeme: "3", tokenType: NUMBER, start: 3, line: 1},
			{lexeme: "*", tokenType: STAR, start: 5, line: 1},
			{lexeme: "3", tokenType: NUMBER, start: 7, line: 1},
			{lexeme: "}}", tokenType: TEMPLATE_RIGHT_BRACE, start: 9, line: 1},
			{lexeme: "", tokenType: EOF, start: 11, line: 1},
		},
		lex.tokens,
	)
}

func TestArithmetic(t *testing.T) {
	lex := NewLexer(`{{1+2-3*4/5}}`)
	lex.run()
	assert.Equal(
		t,
		[]Token{
			{lexeme: "{{", tokenType: TEMPLATE_LEFT_BRACE, start: 0, line: 1},
			{lexeme: "1", tokenType: NUMBER, start: 2, line: 1},
			{lexeme: "+", tokenType: PLUS, start: 3, line: 1},
			{lexeme: "2", tokenType: NUMBER, start: 4, line: 1},
			{lexeme: "-", tokenType: MINUS, start: 5, line: 1},
			{lexeme: "3", tokenType: NUMBER, start: 6, line: 1},
			{lexeme: "*", tokenType: STAR, start: 7, line: 1},
			{lexeme: "4", tokenType: NUMBER, start: 8, line: 1},
			{lexeme: "/", tokenType: SLASH, start: 9, line: 1},
			{lexeme: "5", tokenType: NUMBER, start: 10, line: 1},
			{lexeme: "}}", tokenType: TEMPLATE_RIGHT_BRACE, start: 11, line: 1},
			{lexeme: "", tokenType: EOF, start: 13, line: 1},
		},
		lex.tokens,
	)
}

func TestIdentifiers(t *testing.T) {
	lex := NewLexer(`{{foo.bar}}`)
	lex.run()
	assert.Equal(
		t,
		[]Token{
			{lexeme: "{{", tokenType: TEMPLATE_LEFT_BRACE, start: 0, line: 1},
			{lexeme: "foo", tokenType: 30, start: 2, line: 1},
			{lexeme: ".", tokenType: 11, start: 5, line: 1},
			{lexeme: "bar", tokenType: 30, start: 6, line: 1},
			{lexeme: "}}", tokenType: TEMPLATE_RIGHT_BRACE, start: 9, line: 1},
			{lexeme: "", tokenType: EOF, start: 11, line: 1},
		},
		lex.tokens,
	)
}

func TestFunctionCalls(t *testing.T) {
	lex := NewLexer(`{{foo(bar, "shift")}}`)
	lex.run()
	assert.Equal(
		t,
		[]Token{
			{lexeme: "{{", tokenType: TEMPLATE_LEFT_BRACE, start: 0, line: 1},
			{lexeme: "foo", tokenType: 30, start: 2, line: 1},
			{lexeme: "(", tokenType: 2, start: 5, line: 1},
			{lexeme: "bar", tokenType: 30, start: 6, line: 1},
			{lexeme: ",", tokenType: 10, start: 9, line: 1},
			{lexeme: "\"shift\"", tokenType: 31, start: 11, line: 1},
			{lexeme: ")", tokenType: 3, start: 18, line: 1},
			{lexeme: "}}", tokenType: TEMPLATE_RIGHT_BRACE, start: 19, line: 1},
			{lexeme: "", tokenType: EOF, start: 21, line: 1},
		},
		lex.tokens,
	)
}

func TestNestedCallsAndChaining(t *testing.T) {
	lex := NewLexer(`{{foo(bar(baz), "shift").qux + 1}}`)
	lex.run()
	assert.Equal(
		t,
		[]Token{
			{lexeme: "{{", tokenType: TEMPLATE_LEFT_BRACE, start: 0, line: 1},
			{lexeme: "foo", tokenType: 30, start: 2, line: 1},
			{lexeme: "(", tokenType: 2, start: 5, line: 1},
			{lexeme: "bar", tokenType: 30, start: 6, line: 1},
			{lexeme: "(", tokenType: 2, start: 9, line: 1},
			{lexeme: "baz", tokenType: 30, start: 10, line: 1},
			{lexeme: ")", tokenType: 3, start: 13, line: 1},
			{lexeme: ",", tokenType: 10, start: 14, line: 1},
			{lexeme: "\"shift\"", tokenType: 31, start: 16, line: 1},
			{lexeme: ")", tokenType: 3, start: 23, line: 1},
			{lexeme: ".", tokenType: 11, start: 24, line: 1},
			{lexeme: "qux", tokenType: 30, start: 25, line: 1},
			{lexeme: "+", tokenType: PLUS, start: 29, line: 1},
			{lexeme: "1", tokenType: NUMBER, start: 31, line: 1},
			{lexeme: "}}", tokenType: TEMPLATE_RIGHT_BRACE, start: 32, line: 1},
			{lexeme: "", tokenType: EOF, start: 34, line: 1},
		},
		lex.tokens,
	)
}

func TestTernary(t *testing.T) {
	lex := NewLexer(`{{foo ? bar : baz}}`)
	lex.run()
	assert.Equal(
		t,
		[]Token{
			{lexeme: "{{", tokenType: TEMPLATE_LEFT_BRACE, start: 0, line: 1},
			{lexeme: "foo", tokenType: IDENTIFIER, start: 2, line: 1},
			{lexeme: "?", tokenType: QMARK, start: 6, line: 1},
			{lexeme: "bar", tokenType: IDENTIFIER, start: 8, line: 1},
			{lexeme: ":", tokenType: COLON, start: 12, line: 1},
			{lexeme: "baz", tokenType: IDENTIFIER, start: 14, line: 1},
			{lexeme: "}}", tokenType: TEMPLATE_RIGHT_BRACE, start: 17, line: 1},
			{lexeme: "", tokenType: EOF, start: 19, line: 1},
		},
		lex.tokens,
	)
}

func TestTernaryWithFunctionCalls(t *testing.T) {
	lex := NewLexer(`{{foo(bar) ? baz(qux) : quux(true)}}`)
	lex.run()
	assert.Equal(
		t,
		[]Token{
			{lexeme: "{{", tokenType: TEMPLATE_LEFT_BRACE, start: 0, line: 1},
			{lexeme: "foo", tokenType: IDENTIFIER, start: 2, line: 1},
			{lexeme: "(", tokenType: LEFT_PAREN, start: 5, line: 1},
			{lexeme: "bar", tokenType: IDENTIFIER, start: 6, line: 1},
			{lexeme: ")", tokenType: RIGHT_PAREN, start: 9, line: 1},
			{lexeme: "?", tokenType: QMARK, start: 11, line: 1},
			{lexeme: "baz", tokenType: IDENTIFIER, start: 13, line: 1},
			{lexeme: "(", tokenType: LEFT_PAREN, start: 16, line: 1},
			{lexeme: "qux", tokenType: IDENTIFIER, start: 17, line: 1},
			{lexeme: ")", tokenType: RIGHT_PAREN, start: 20, line: 1},
			{lexeme: ":", tokenType: COLON, start: 22, line: 1},
			{lexeme: "quux", tokenType: IDENTIFIER, start: 24, line: 1},
			{lexeme: "(", tokenType: LEFT_PAREN, start: 28, line: 1},
			{lexeme: "true", tokenType: TRUE, start: 29, line: 1},
			{lexeme: ")", tokenType: RIGHT_PAREN, start: 33, line: 1},
			{lexeme: "}}", tokenType: TEMPLATE_RIGHT_BRACE, start: 34, line: 1},
			{lexeme: "", tokenType: EOF, start: 36, line: 1},
		},
		lex.tokens,
	)
}
