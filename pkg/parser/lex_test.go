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
	lex := NewLexer(`@{{"Hello"}}`)
	lex.run()
	assert.Equal(
		t,
		[]Token{
			{lexeme: "@{{", tokenType: TEMPLATE_LEFT_BRACE, start: 0, line: 1},
			{lexeme: "\"Hello\"", tokenType: STRING, start: 3, line: 1},
			{lexeme: "}}", tokenType: TEMPLATE_RIGHT_BRACE, start: 10, line: 1},
			{lexeme: "", tokenType: EOF, start: 12, line: 1},
		},
		lex.tokens,
	)
}

func TestBasicArithmetic(t *testing.T) {
	lex := NewLexer("@{{ 3 * 3 }}")
	lex.run()
	assert.Equal(
		t,
		[]Token{
			{lexeme: "@{{", tokenType: TEMPLATE_LEFT_BRACE, start: 0, line: 1},
			{lexeme: "3", tokenType: NUMBER, start: 4, line: 1},
			{lexeme: "*", tokenType: STAR, start: 6, line: 1},
			{lexeme: "3", tokenType: NUMBER, start: 8, line: 1},
			{lexeme: "}}", tokenType: TEMPLATE_RIGHT_BRACE, start: 10, line: 1},
			{lexeme: "", tokenType: EOF, start: 12, line: 1},
		},
		lex.tokens,
	)
}

func TestArithmetic(t *testing.T) {
	lex := NewLexer(`@{{1+2-3*4/5}}`)
	lex.run()
	assert.Equal(
		t,
		[]Token{
			{lexeme: "@{{", tokenType: TEMPLATE_LEFT_BRACE, start: 0, line: 1},
			{lexeme: "1", tokenType: NUMBER, start: 3, line: 1},
			{lexeme: "+", tokenType: PLUS, start: 4, line: 1},
			{lexeme: "2", tokenType: NUMBER, start: 5, line: 1},
			{lexeme: "-", tokenType: MINUS, start: 6, line: 1},
			{lexeme: "3", tokenType: NUMBER, start: 7, line: 1},
			{lexeme: "*", tokenType: STAR, start: 8, line: 1},
			{lexeme: "4", tokenType: NUMBER, start: 9, line: 1},
			{lexeme: "/", tokenType: SLASH, start: 10, line: 1},
			{lexeme: "5", tokenType: NUMBER, start: 11, line: 1},
			{lexeme: "}}", tokenType: TEMPLATE_RIGHT_BRACE, start: 12, line: 1},
			{lexeme: "", tokenType: EOF, start: 14, line: 1},
		},
		lex.tokens,
	)
}

func TestIdentifiers(t *testing.T) {
	lex := NewLexer(`@{{foo.bar}}`)
	lex.run()
	assert.Equal(
		t,
		[]Token{
			{lexeme: "@{{", tokenType: TEMPLATE_LEFT_BRACE, start: 0, line: 1},
			{lexeme: "foo", tokenType: IDENTIFIER, start: 3, line: 1},
			{lexeme: ".", tokenType: 11, start: 6, line: 1},
			{lexeme: "bar", tokenType: IDENTIFIER, start: 7, line: 1},
			{lexeme: "}}", tokenType: TEMPLATE_RIGHT_BRACE, start: 10, line: 1},
			{lexeme: "", tokenType: EOF, start: 12, line: 1},
		},
		lex.tokens,
	)
}

func TestFunctionCalls(t *testing.T) {
	lex := NewLexer(`@{{foo(bar, "shift")}}`)
	lex.run()
	assert.Equal(
		t,
		[]Token{
			{lexeme: "@{{", tokenType: TEMPLATE_LEFT_BRACE, start: 0, line: 1},
			{lexeme: "foo", tokenType: IDENTIFIER, start: 3, line: 1},
			{lexeme: "(", tokenType: LEFT_PAREN, start: 6, line: 1},
			{lexeme: "bar", tokenType: IDENTIFIER, start: 7, line: 1},
			{lexeme: ",", tokenType: COMMA, start: 10, line: 1},
			{lexeme: "\"shift\"", tokenType: STRING, start: 12, line: 1},
			{lexeme: ")", tokenType: RIGHT_PAREN, start: 19, line: 1},
			{lexeme: "}}", tokenType: TEMPLATE_RIGHT_BRACE, start: 20, line: 1},
			{lexeme: "", tokenType: EOF, start: 22, line: 1},
		},
		lex.tokens,
	)
}

func TestNestedCallsAndChaining(t *testing.T) {
	lex := NewLexer(`@{{foo(bar(baz), "shift").qux + 1}}`)
	lex.run()
	assert.Equal(
		t,
		[]Token{
			{lexeme: "@{{", tokenType: TEMPLATE_LEFT_BRACE, start: 0, line: 1},
			{lexeme: "foo", tokenType: IDENTIFIER, start: 3, line: 1},
			{lexeme: "(", tokenType: LEFT_PAREN, start: 6, line: 1},
			{lexeme: "bar", tokenType: IDENTIFIER, start: 7, line: 1},
			{lexeme: "(", tokenType: LEFT_PAREN, start: 10, line: 1},
			{lexeme: "baz", tokenType: IDENTIFIER, start: 11, line: 1},
			{lexeme: ")", tokenType: RIGHT_PAREN, start: 14, line: 1},
			{lexeme: ",", tokenType: COMMA, start: 15, line: 1},
			{lexeme: "\"shift\"", tokenType: STRING, start: 17, line: 1},
			{lexeme: ")", tokenType: RIGHT_PAREN, start: 24, line: 1},
			{lexeme: ".", tokenType: DOT, start: 25, line: 1},
			{lexeme: "qux", tokenType: IDENTIFIER, start: 26, line: 1},
			{lexeme: "+", tokenType: PLUS, start: 30, line: 1},
			{lexeme: "1", tokenType: NUMBER, start: 32, line: 1},
			{lexeme: "}}", tokenType: TEMPLATE_RIGHT_BRACE, start: 33, line: 1},
			{lexeme: "", tokenType: EOF, start: 35, line: 1},
		},
		lex.tokens,
	)
}

func TestTernary(t *testing.T) {
	lex := NewLexer(`@{{foo ? bar : baz}}`)
	lex.run()
	assert.Equal(
		t,
		[]Token{
			{lexeme: "@{{", tokenType: TEMPLATE_LEFT_BRACE, start: 0, line: 1},
			{lexeme: "foo", tokenType: IDENTIFIER, start: 3, line: 1},
			{lexeme: "?", tokenType: QMARK, start: 7, line: 1},
			{lexeme: "bar", tokenType: IDENTIFIER, start: 9, line: 1},
			{lexeme: ":", tokenType: COLON, start: 13, line: 1},
			{lexeme: "baz", tokenType: IDENTIFIER, start: 15, line: 1},
			{lexeme: "}}", tokenType: TEMPLATE_RIGHT_BRACE, start: 18, line: 1},
			{lexeme: "", tokenType: EOF, start: 20, line: 1},
		},
		lex.tokens,
	)
}

func TestTernaryWithFunctionCalls(t *testing.T) {
	lex := NewLexer(`@{{foo(bar) ? baz(qux) : quux(true)}}`)
	lex.run()
	assert.Equal(
		t,
		[]Token{
			{lexeme: "@{{", tokenType: TEMPLATE_LEFT_BRACE, start: 0, line: 1},
			{lexeme: "foo", tokenType: IDENTIFIER, start: 3, line: 1},
			{lexeme: "(", tokenType: LEFT_PAREN, start: 6, line: 1},
			{lexeme: "bar", tokenType: IDENTIFIER, start: 7, line: 1},
			{lexeme: ")", tokenType: RIGHT_PAREN, start: 10, line: 1},
			{lexeme: "?", tokenType: QMARK, start: 12, line: 1},
			{lexeme: "baz", tokenType: IDENTIFIER, start: 14, line: 1},
			{lexeme: "(", tokenType: LEFT_PAREN, start: 17, line: 1},
			{lexeme: "qux", tokenType: IDENTIFIER, start: 18, line: 1},
			{lexeme: ")", tokenType: RIGHT_PAREN, start: 21, line: 1},
			{lexeme: ":", tokenType: COLON, start: 23, line: 1},
			{lexeme: "quux", tokenType: IDENTIFIER, start: 25, line: 1},
			{lexeme: "(", tokenType: LEFT_PAREN, start: 29, line: 1},
			{lexeme: "true", tokenType: TRUE, start: 30, line: 1},
			{lexeme: ")", tokenType: RIGHT_PAREN, start: 34, line: 1},
			{lexeme: "}}", tokenType: TEMPLATE_RIGHT_BRACE, start: 35, line: 1},
			{lexeme: "", tokenType: EOF, start: 37, line: 1},
		},
		lex.tokens,
	)
}
