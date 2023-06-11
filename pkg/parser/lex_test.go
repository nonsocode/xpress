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
