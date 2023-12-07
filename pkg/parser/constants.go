package parser

const (
	// Single-character tokens.
	EOF TokenType = iota
	ERROR
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	LEFT_BRACKET
	RIGHT_BRACKET
	PERCENT
	COLON
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR
	QMARK
	// One or two character tokens.
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL
	TEMPLATE_LEFT_BRACE
	TEMPLATE_RIGHT_BRACE
	OPTIONALCHAIN
	AND
	OR
	// Literals.
	IDENTIFIER
	STRING
	NUMBER
	TEXT
	// Keywords.
	_keywordStart
	FALSE
	TRUE
	NIL
)
