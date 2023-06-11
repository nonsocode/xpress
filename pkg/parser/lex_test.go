package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexerCanLexText(t *testing.T) {
	lex := NewLexer("Hello World")
	lex.run()
	assert.Equal(t, []Token{
		{"Hello World", TEXT},
		{"", EOF},
	}, lex.tokens)
}

func TestLexerCanLexTextWithNewlines(t *testing.T) {
	lex := NewLexer("Hello\nWorld")
	lex.run()
	assert.Equal(t, []Token{
		{"Hello\nWorld", TEXT},
		{"", EOF},
	}, lex.tokens)
}

func TestLexerCanLexAction(t *testing.T) {
	lex := NewLexer(`{{"Hello"}}`)
	lex.run()
	assert.Equal(t, []Token{
		{lexeme: "{{", tokenType: TEMPLATE_LEFT_BRACE},
		{lexeme: "\"Hello\"", tokenType: STRING},
		{lexeme: "}}", tokenType: TEMPLATE_RIGHT_BRACE},
		{lexeme: "", tokenType: EOF},
	}, lex.tokens)
}

func TestArithmetic(t *testing.T) {
	lex := NewLexer(`{{1+2-3*4/5}}`)
	lex.run()
	assert.Equal(t, []Token{
		{lexeme: "{{", tokenType: TEMPLATE_LEFT_BRACE},
		{lexeme: "1", tokenType: NUMBER},
		{lexeme: "+", tokenType: PLUS},
		{lexeme: "2", tokenType: NUMBER},
		{lexeme: "-", tokenType: MINUS},
		{lexeme: "3", tokenType: NUMBER},
		{lexeme: "*", tokenType: STAR},
		{lexeme: "4", tokenType: NUMBER},
		{lexeme: "/", tokenType: SLASH},
		{lexeme: "5", tokenType: NUMBER},
		{lexeme: "}}", tokenType: TEMPLATE_RIGHT_BRACE},
		{lexeme: "", tokenType: EOF},
	}, lex.tokens)
}

func TestIdentifiers(t *testing.T) {
	lex := NewLexer(`{{foo.bar}}`)
	lex.run()
	assert.Equal(t, []Token{
		{lexeme: "{{", tokenType: TEMPLATE_LEFT_BRACE},
		{lexeme: "foo", tokenType: IDENTIFIER},
		{lexeme: ".", tokenType: DOT},
		{lexeme: "bar", tokenType: IDENTIFIER},
		{lexeme: "}}", tokenType: TEMPLATE_RIGHT_BRACE},
		{lexeme: "", tokenType: EOF},
	}, lex.tokens)
}

func TestFunctionCalls(t *testing.T) {
	lex := NewLexer(`{{foo(bar, "shift")}}`)
	lex.run()
	assert.Equal(t, []Token{
		{lexeme: "{{", tokenType: TEMPLATE_LEFT_BRACE},
		{lexeme: "foo", tokenType: IDENTIFIER},
		{lexeme: "(", tokenType: LEFT_PAREN},
		{lexeme: "bar", tokenType: IDENTIFIER},
		{lexeme: ",", tokenType: COMMA},
		{lexeme: "\"shift\"", tokenType: STRING},
		{lexeme: ")", tokenType: RIGHT_PAREN},
		{lexeme: "}}", tokenType: TEMPLATE_RIGHT_BRACE},
		{lexeme: "", tokenType: EOF},
	}, lex.tokens)
}

func TestNestedCallsAndChaining(t *testing.T) {
	lex := NewLexer(`{{foo(bar(baz), "shift").qux + 1}}`)
	lex.run()
	assert.Equal(t, []Token{
		{lexeme: "{{", tokenType: TEMPLATE_LEFT_BRACE},
		{lexeme: "foo", tokenType: IDENTIFIER},
		{lexeme: "(", tokenType: LEFT_PAREN},
		{lexeme: "bar", tokenType: IDENTIFIER},
		{lexeme: "(", tokenType: LEFT_PAREN},
		{lexeme: "baz", tokenType: IDENTIFIER},
		{lexeme: ")", tokenType: RIGHT_PAREN},
		{lexeme: ",", tokenType: COMMA},
		{lexeme: "\"shift\"", tokenType: STRING},
		{lexeme: ")", tokenType: RIGHT_PAREN},
		{lexeme: ".", tokenType: DOT},
		{lexeme: "qux", tokenType: IDENTIFIER},
		{lexeme: "+", tokenType: PLUS},
		{lexeme: "1", tokenType: NUMBER},
		{lexeme: "}}", tokenType: TEMPLATE_RIGHT_BRACE},
		{lexeme: "", tokenType: EOF},
	}, lex.tokens)
}

func TestTernary(t *testing.T) {
	lex := NewLexer(`{{foo ? bar : baz}}`)
	lex.run()
	assert.Equal(t, []Token{
		{lexeme: "{{", tokenType: TEMPLATE_LEFT_BRACE},
		{lexeme: "foo", tokenType: IDENTIFIER},
		{lexeme: "?", tokenType: QMARK},
		{lexeme: "bar", tokenType: IDENTIFIER},
		{lexeme: ":", tokenType: COLON},
		{lexeme: "baz", tokenType: IDENTIFIER},
		{lexeme: "}}", tokenType: TEMPLATE_RIGHT_BRACE},
		{lexeme: "", tokenType: EOF},
	}, lex.tokens)
}

func TestTernaryWithFunctionCalls(t *testing.T) {
	lex := NewLexer(`{{foo(bar) ? baz(qux) : quux(true)}}`)
	lex.run()
	assert.Equal(t, []Token{
		{lexeme: "{{", tokenType: TEMPLATE_LEFT_BRACE},
		{lexeme: "foo", tokenType: IDENTIFIER},
		{lexeme: "(", tokenType: LEFT_PAREN},
		{lexeme: "bar", tokenType: IDENTIFIER},
		{lexeme: ")", tokenType: RIGHT_PAREN},
		{lexeme: "?", tokenType: QMARK},
		{lexeme: "baz", tokenType: IDENTIFIER},
		{lexeme: "(", tokenType: LEFT_PAREN},
		{lexeme: "qux", tokenType: IDENTIFIER},
		{lexeme: ")", tokenType: RIGHT_PAREN},
		{lexeme: ":", tokenType: COLON},
		{lexeme: "quux", tokenType: IDENTIFIER},
		{lexeme: "(", tokenType: LEFT_PAREN},
		{lexeme: "true", tokenType: TRUE},
		{lexeme: ")", tokenType: RIGHT_PAREN},
		{lexeme: "}}", tokenType: TEMPLATE_RIGHT_BRACE},
		{lexeme: "", tokenType: EOF},
	}, lex.tokens)
}
