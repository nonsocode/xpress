package parser

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	leftDelim    = "@{{"
	rightDelim   = "}}"
	eof          = -1
	alpha        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digit        = "0123456789"
	alphaNumeric = alpha + digit
)

var (
	keywords = map[string]TokenType{
		"and":   AND,
		"or":    OR,
		"false": FALSE,
		"true":  TRUE,
		"nil":   NIL,
	}

	tokenMap = map[TokenType]string{
		EOF:                  "EOF",
		ERROR:                "ERROR",
		LEFT_PAREN:           "LEFT_PAREN",
		RIGHT_PAREN:          "RIGHT_PAREN",
		LEFT_BRACE:           "LEFT_BRACE",
		RIGHT_BRACE:          "RIGHT_BRACE",
		LEFT_BRACKET:         "LEFT_BRACKET",
		RIGHT_BRACKET:        "RIGHT_BRACKET",
		PERCENT:              "PERCENT",
		COLON:                "COLON",
		COMMA:                "COMMA",
		DOT:                  "DOT",
		MINUS:                "MINUS",
		PLUS:                 "PLUS",
		SEMICOLON:            "SEMICOLON",
		SLASH:                "SLASH",
		STAR:                 "STAR",
		QMARK:                "QMARK",
		BANG:                 "BANG",
		BANG_EQUAL:           "BANG_EQUAL",
		EQUAL:                "EQUAL",
		EQUAL_EQUAL:          "EQUAL_EQUAL",
		GREATER:              "GREATER",
		GREATER_EQUAL:        "GREATER_EQUAL",
		LESS:                 "LESS",
		LESS_EQUAL:           "LESS_EQUAL",
		TEMPLATE_LEFT_BRACE:  "TEMPLATE_LEFT_BRACE",
		TEMPLATE_RIGHT_BRACE: "TEMPLATE_RIGHT_BRACE",
		OPTIONALCHAIN:        "OPTIONALCHAIN",
		AND:                  "AND",
		OR:                   "OR",
		IDENTIFIER:           "IDENTIFIER",
		STRING:               "STRING",
		NUMBER:               "NUMBER",
		TEXT:                 "TEXT",
		FALSE:                "FALSE",
		TRUE:                 "TRUE",
		NIL:                  "NIL",
	}
)

type (
	Token struct {
		lexeme    string
		tokenType TokenType
		start     int
		line      int
	}
	TokenType int

	Lexer struct {
		source    string
		tokens    []Token
		start     int
		lineStart int
		current   int
		line      int
		nesting   int
	}
	stateFn func(*Lexer) stateFn
)

func NewLexer(source string) *Lexer {
	return &Lexer{
		source:    source,
		tokens:    make([]Token, 0),
		start:     0,
		current:   0,
		lineStart: 1,
		line:      1,
	}
}

func (l *Lexer) scanTokens() []Token {
	l.run()
	return l.tokens
}
func (l *Lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	l.tokens = append(l.tokens, Token{"", EOF, l.current, l.lineStart})
}

func (l *Lexer) addToken(tokenType TokenType) {
	text := l.source[l.start:l.current]
	l.tokens = append(l.tokens, Token{text, tokenType, l.start, l.lineStart})
	l.start = l.current
	l.lineStart = l.line
}

func (l *Lexer) errorf(err string, args ...any) stateFn {
	l.tokens = append(l.tokens, Token{fmt.Sprintf(err, args...), ERROR, l.start, l.line})
	return nil
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) ignore() {
	l.start = l.current
	l.lineStart = l.line
}

func (l *Lexer) next() rune {
	if int(l.current) >= len(l.source) {
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.source[l.current:])
	l.current += w
	if r == '\n' {
		l.line++
	}
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune.
func (l *Lexer) backup() {
	if !l.isAtEnd() && l.current > 0 {
		r, w := utf8.DecodeLastRuneInString(l.source[:l.current])
		l.current -= (w)
		// Correct newline count.
		if r == '\n' {
			l.line--
		}
	}
}

// accept consumes the next rune if it's from the valid set.
func (l *Lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *Lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func (l *Lexer) atTerminator() bool {
	r := l.peek()
	if isSpace(r) {
		return true
	}
	switch r {
	case eof,
		'.',
		'?',
		',',
		'|',
		':',
		')',
		'(',
		'[',
		']',
		'+',
		'-',
		'*',
		'/',
		'%',
		'^',
		'=',
		'!',
		'<',
		'>',
		'&',
		';',
		'{',
		'}':
		return true
	}
	return strings.HasPrefix(l.source[l.current:], rightDelim)
}

func lexText(l *Lexer) stateFn {
	if x := strings.Index(l.source[l.start:], leftDelim); x >= 0 {
		if x > 0 {
			l.current += x
			l.line += strings.Count(l.source[l.start:l.current], "\n")
			l.addToken(TEXT)
		}
		return lexLeftDelim
	}
	l.current += len(l.source[l.start:])
	l.line += strings.Count(l.source[l.start:l.current], "\n")
	if l.start >= l.current {
		return nil
	}

	l.addToken(TEXT)
	return nil
}

func lexLeftDelim(l *Lexer) stateFn {
	l.current += len(leftDelim)
	l.addToken(TEMPLATE_LEFT_BRACE)
	l.nesting++
	return lexInsideAction
}

func lexRightDelim(l *Lexer) stateFn {
	l.current += len(rightDelim)
	l.addToken(TEMPLATE_RIGHT_BRACE)
	l.nesting--
	return lexText
}

func lexInsideAction(l *Lexer) stateFn {
	for {
		if strings.HasPrefix(l.source[l.current:], rightDelim) && l.nesting == 1 {
			return lexRightDelim
		}
		if l.isAtEnd() {
			return l.errorf("unclosed action")
		}
		switch c := l.next(); c {
		case '"':
			return lexDquote
		case '\'':
			return lexSquote
		case '(':
			l.addToken(LEFT_PAREN)
			l.nesting++
		case ')':
			l.addToken(RIGHT_PAREN)
			l.nesting--
		case ',':
			l.addToken(COMMA)
		case '.':
			l.addToken(DOT)
			return lexIdent
		case '-':
			l.addToken(MINUS)
		case '+':
			l.addToken(PLUS)
		case '*':
			l.addToken(STAR)
		case '/':
			l.addToken(SLASH)
		case '%':
			l.addToken(PERCENT)
		case '?':
			if l.accept(".") {
				l.addToken(OPTIONALCHAIN)
			} else if l.accept("?") {
				l.addToken(NULLCOALESCING)
			} else {
				l.addToken(QMARK)
			}
		case '[':
			l.addToken(LEFT_BRACKET)
			l.nesting++
		case ']':
			l.addToken(RIGHT_BRACKET)
			l.nesting--
		case ':':
			l.addToken(COLON)
			// numbers
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			l.backup()
			return lexNumber
		case '|':
			if l.accept("|") {
				l.addToken(OR)
			} else {
				return l.errorf(`illegal character "%s"`, string(c))
			}
		case '&':
			if l.accept("&") {
				l.addToken(AND)
			} else {
				return l.errorf(`illegal character "%s"`, string(c))
			}
		case '!':
			if l.accept("=") {
				l.addToken(BANG_EQUAL)
			} else {
				l.addToken(BANG)
			}
		case '=':
			if l.accept("=") {
				l.addToken(EQUAL_EQUAL)
			} else {
				l.addToken(EQUAL)
			}
		case '<':
			if l.accept("=") {
				l.addToken(LESS_EQUAL)
			} else {
				l.addToken(LESS)
			}
		case '>':
			if l.accept("=") {
				l.addToken(GREATER_EQUAL)
			} else {
				l.addToken(GREATER)
			}
		case ' ', '\t', '\r', '\n':
			l.ignore()
		case '{':
			l.addToken(LEFT_BRACE)
			l.nesting++
		case '}':
			l.addToken(RIGHT_BRACE)
			l.nesting--
		default:
			if isAlphaNumeric(c) {
				return lexIdent
			}
			return l.errorf(`illegal character "%s"`, string(c))
		}
	}
}

func lexIdent(l *Lexer) stateFn {
	var r rune
	for {
		r = l.next()
		if !isAlphaNumeric(r) {
			l.backup()
			break
		}
	}
	if !l.atTerminator() {
		return l.errorf("bad character %#U", r)
	}
	word := l.source[l.start:l.current]
	if len(word) > 0 {
		if keywords[word] > _keywordStart {
			l.addToken(keywords[word])
		} else {
			l.addToken(IDENTIFIER)
		}
	}
	return lexInsideAction
}

// lexDquote scans a quoted string.
func lexDquote(l *Lexer) stateFn {
	err := l.scanRawString('"')
	if err != nil {
		return l.errorf(err.Error())
	}
	l.addToken(STRING)
	return lexInsideAction
}

func lexSquote(l *Lexer) stateFn {
	err := l.scanRawString('\'')
	if err != nil {
		return l.errorf(err.Error())
	}
	l.addToken(STRING)
	return lexInsideAction
}

func (l *Lexer) scanRawString(delim rune) error {
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof, '\n':
			return errors.New("unterminated quoted string")
		case delim:
			return nil
		}
	}
}

// lexNumber scans a number: decimal, octal, hex and float. This
// isn't a perfect number scanner - for instance it accepts "." and "0x0.2"
// and "089" - but when it's wrong the input is invalid and the parser (via
// strconv) will notice.
func lexNumber(l *Lexer) stateFn {
	if !l.scanNumber() {
		return l.errorf("bad number syntax: %q", l.source[l.start:l.current])
	}

	l.addToken(NUMBER)
	return lexInsideAction
}

func (l *Lexer) scanNumber() bool {
	if l.accept("_") {
		return false
	}
	digits := "0123456789_"
	if l.accept("0") {
		// Note: Leading 0 does not mean octal in floats.
		if l.accept("xX") {
			digits = "0123456789abcdefABCDEF_"
		} else if l.accept("oO") {
			digits = "01234567_"
		} else if l.accept("bB") {
			digits = "01_"
		}
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}

	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(l.peek()) {
		l.next()
		return false
	}
	return true
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
