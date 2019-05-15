package main

import (
	"fmt"
	"strconv"
)

type scanner struct {
	source  []rune
	start   int
	current int
	line    int
	tokens  []token
}

func newScanner(code string) scanner {
	return scanner{
		source:  []rune(code),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (sc *scanner) scanTokens() ([]token, error) {
	for !sc.isAtEnd() {
		sc.start = sc.current
		if err := sc.scanToken(); err != nil {
			return nil, err
		}
	}

	sc.tokens = append(sc.tokens, token{Type: tokenEndOfFile, Line: sc.line})
	return sc.tokens, nil
}

func (sc *scanner) isAtEnd() bool {
	return sc.current >= len(sc.source)
}

func (sc *scanner) scanToken() error {
	c := sc.advance()
	switch c {
	case '(':
		sc.addEmptyToken(tokenLeftParenthesis)
		break
	case ')':
		sc.addEmptyToken(tokenRightParenthesis)
		break
	case '{':
		sc.addEmptyToken(tokenLeftBrace)
		break
	case '}':
		sc.addEmptyToken(tokenRightBrace)
		break
	case '[':
		sc.addEmptyToken(tokenLeftBracket)
		break
	case ']':
		sc.addEmptyToken(tokenRightBracket)
		break
	case '+':
		sc.addEmptyToken(tokenPlus)
		break
	case '-':
		sc.addEmptyToken(tokenMinus)
		break
	case '*':
		sc.addEmptyToken(tokenStar)
		break
	case '.':
		sc.addEmptyToken(tokenDot)
		break
	case ',':
		sc.addEmptyToken(tokenComma)
		break
	case ';':
		sc.addEmptyToken(tokenSemiColon)
		break
	case '!':
		if sc.match('=') {
			sc.addEmptyToken(tokenBangEqual)
		} else {
			sc.addEmptyToken(tokenBang)
		}
		break
	case '/':
		if sc.match('/') {
			// A comment goes until the end of the line.
			for sc.peek() != '\n' && !sc.isAtEnd() {
				sc.advance()
			}
			break
		}
		sc.addEmptyToken(tokenSlash)
		break
	case '=':
		if sc.match('=') {
			sc.addEmptyToken(tokenEqualEqual)
		} else {
			sc.addEmptyToken(tokenEqual)
		}
		break
	case '>':
		if sc.match('=') {
			sc.addEmptyToken(tokenGreaterEqual)
		} else {
			sc.addEmptyToken(tokenGreater)
		}
		break
	case '<':
		if sc.match('=') {
			sc.addEmptyToken(tokenLessEqual)
		} else {
			sc.addEmptyToken(tokenLess)
		}
		break
	case ' ':
	case '\r':
	case '\t':
		// Ignore whitespace.
		break
	case '\n':
		sc.line++
		break
	case '"':
		sc.string()
		break
	default:
		if sc.isDigit(c) {
			sc.number()
		} else if sc.isAlpha(c) {
			sc.identifier()
		} else {
			return fmt.Errorf("Error: Line: %d, Unexpected character: %c", sc.line, c)
		}
	}
	return nil
}

func (sc *scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (sc *scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func (sc *scanner) isAlphaNumeric(c rune) bool {
	return sc.isAlpha(c) || sc.isDigit(c)
}

func (sc *scanner) identifier() {
	for sc.isAlphaNumeric(sc.peek()) {
		sc.advance()
	}

	text := sc.source[sc.start:sc.current]
	if tokenType, exist := keywords[string(text)]; exist {
		sc.addEmptyToken(tokenType)
	} else {
		sc.addEmptyToken(tokenIdentifier)
	}
}

func (sc *scanner) number() {
	for sc.isDigit(sc.peek()) {
		sc.advance()
	}

	// Look for a fractional part.
	if sc.peek() == '.' && sc.isDigit(sc.peekNext()) {
		// Consume the "."
		sc.advance()

		for sc.isDigit(sc.peek()) {
			sc.advance()
		}
	}

	float, err := strconv.ParseFloat(string(sc.source[sc.start:sc.current]), 64)
	if err == nil {
		sc.addToken(tokenNumber, float)
	}
}

func (sc *scanner) peek() rune {
	if sc.isAtEnd() {
		return rune(-1)
	}
	return sc.source[sc.current]
}

func (sc *scanner) peekNext() rune {
	if sc.current+1 >= len(sc.source) {
		return rune(-1)
	}
	return sc.source[sc.current+1]
}

func (sc *scanner) match(c rune) bool {
	if sc.isAtEnd() {
		return false
	}
	if sc.source[sc.current] != c {
		return false
	}
	sc.current++
	return true
}

func (sc *scanner) addToken(t tokenType, lit interface{}) {
	text := sc.source[sc.start:sc.current]
	sc.tokens = append(sc.tokens, token{
		Type:    t,
		Lexeme:  string(text),
		Literal: lit,
		Line:    sc.line,
	})
}

func (sc *scanner) addEmptyToken(t tokenType) {
	sc.addToken(t, nil)
}

func (sc *scanner) advance() rune {
	c := sc.source[sc.current]
	sc.current++
	return c
}

func (sc *scanner) string() {

	for sc.peek() != '"' && !sc.isAtEnd() {
		if sc.peek() == '\n' {
			sc.line++
		}
		sc.advance()

	}

	// Unterminated string.
	if sc.isAtEnd() {
		panic(fmt.Sprintf("ERROR: Line: %d, Unterminated string.", sc.line))
	}

	// The closing ".
	sc.advance()

	// Trim the surrounding quotes.
	value := sc.source[sc.start+1 : sc.current-1]
	sc.addToken(tokenString, string(value))
}
