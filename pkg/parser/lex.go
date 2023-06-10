package parser

import "fmt"

type Token struct {
	lexeme    string
	tokenType TokenType
}
type TokenType int

type Lexer struct {
	source  string
	tokens  []Token
	start   int
	current int
}

func NewLexer(source string) *Lexer {
	return &Lexer{source: source}
}

func (l *Lexer) scanTokens() []Token {
	for !l.isAtEnd() {
		l.start = l.current
		l.scanToken()
	}
	l.tokens = append(l.tokens, NewToken("", EOF))
	return l.tokens
}

func (l *Lexer) scanToken() {
	c := l.advance()
	switch c {
	case '(':
		l.addToken(LEFT_PAREN)
	case ')':
		l.addToken(RIGHT_PAREN)
	case '{':
		l.addToken(LEFT_BRACE)
	case '}':
		l.addToken(RIGHT_BRACE)
	case ',':
		l.addToken(COMMA)
	case '.':
		l.addToken(DOT)
	case '-':
		l.addToken(MINUS)
	case '+':
		l.addToken(PLUS)
	case ';':
		l.addToken(SEMICOLON)
	case '*':
		l.addToken(STAR)
	case '?':
		l.addToken(QMARK)
	case '!':
		if l.match('=') {
			l.addToken(BANG_EQUAL)
		} else {
			l.addToken(BANG)
		}
	case '=':
		if l.match('=') {
			l.addToken(EQUAL_EQUAL)
		} else {
			l.addToken(EQUAL)
		}
	case '<':
		if l.match('=') {
			l.addToken(LESS_EQUAL)
		} else {
			l.addToken(LESS)
		}
	case '>':
		if l.match('=') {
			l.addToken(GREATER_EQUAL)
		} else {
			l.addToken(GREATER)
		}
	case '/':
		if l.match('/') {
			for l.peek() != '\n' && !l.isAtEnd() {
				l.advance()
			}
		} else {
			l.addToken(SLASH)
		}
	case ' ', '\r', '\t':
		// Ignore whitespace.
	case '\n':
		l.addToken(SEMICOLON)
	case '"':
		l.string()
	default:
		if isDigit(c) {
			l.number()
		}
	}
}

func (l *Lexer) string() {
	for l.peek() != '"' && !l.isAtEnd() {
		if l.peek() == '\n' {
			l.advance()
		}
	}
	if l.isAtEnd() {
		fmt.Println("Unterminated string.")
		return
	}
	l.advance()
	value := l.source[l.start+1 : l.current-1]
	l.addTokenLiteral(STRING, value)
}

func (l *Lexer) number() {
	for isDigit(l.peek()) {
		l.advance()
	}
	if l.peek() == '.' && isDigit(l.peekNext()) {
		l.advance()
		for isDigit(l.peek()) {
			l.advance()
		}
	}
	l.addTokenLiteral(NUMBER, l.source[l.start:l.current])
}

func (l *Lexer) peekNext() byte {
	if l.current+1 >= len(l.source) {
		return '\000'
	}
	return l.source[l.current+1]
}

func (l *Lexer) match(expected byte) bool {
	if l.isAtEnd() {
		return false
	}
	if l.source[l.current] != expected {
		return false
	}
	l.current++
	return true
}

func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return '\000'
	}
	return l.source[l.current]
}

func (l *Lexer) advance() byte {
	l.current++
	return l.source[l.current-1]
}

func (l *Lexer) addToken(tokenType TokenType) {
	l.addTokenLiteral(tokenType, nil)
}

func (l *Lexer) addTokenLiteral(tokenType TokenType, literal interface{}) {
	text := l.source[l.start:l.current]
	l.tokens = append(l.tokens, NewToken(text, tokenType))
	l.start = l.current
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func NewToken(lexeme string, tokenType TokenType) Token {
	return Token{lexeme: lexeme, tokenType: tokenType}
}
