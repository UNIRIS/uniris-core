package interpreter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializeScanner(t *testing.T) {
	s := newScanner("print")
	assert.Equal(t, 0, s.current)
	assert.Equal(t, 1, s.line)
	assert.Len(t, s.source, 5)
}

func TestScanAdvance(t *testing.T) {
	s := newScanner("hello")
	s.advance()
	assert.Equal(t, 1, s.current)
	s.advance()
	assert.Equal(t, 2, s.current)
}

func TestIsAtEnd(t *testing.T) {
	s := newScanner("")
	s.current = 1
	assert.True(t, s.isAtEnd())

	s = newScanner("print")
	assert.False(t, s.isAtEnd())
}

func TestAddEmptyToken(t *testing.T) {
	s := newScanner("(")
	s.addEmptyToken(tokenLeftParenthesis)
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenLeftParenthesis, s.tokens[0].Type)
}

func TestScanIsAlpha(t *testing.T) {
	s := newScanner("")
	assert.True(t, s.isAlpha(rune('a')))
	assert.False(t, s.isAlpha(rune(10)))
}

func TestScanIsDigit(t *testing.T) {
	s := newScanner("")
	assert.True(t, s.isDigit(rune('1')))
}

func TestScanIsAlphanumeric(t *testing.T) {
	s := newScanner("")
	assert.True(t, s.isAlphaNumeric(rune('a')))
	assert.True(t, s.isAlphaNumeric(rune('1')))
}

func TestScanIdentifier(t *testing.T) {
	s := newScanner("a")
	s.identifier()
	assert.Equal(t, tokenIdentifier, s.tokens[0].Type)

	s = newScanner("if")
	s.identifier()
	assert.Equal(t, tokenIf, s.tokens[0].Type)
}

func TestScanStringDoubleQuotes(t *testing.T) {
	s := newScanner("\"hello\"")
	s.advance()
	s.stringDoubleQuotes()
	assert.Equal(t, tokenString, s.tokens[0].Type)
	assert.Equal(t, "hello", s.tokens[0].Literal)
	assert.Equal(t, "\"hello\"", s.tokens[0].Lexeme)

	s = newScanner("\"hello")
	s.advance()
	assert.Error(t, s.stringDoubleQuotes())

	s = newScanner("\"hello\nWorld\"")
	s.advance()
	s.stringDoubleQuotes()
	assert.Equal(t, 2, s.line)
}

func TestScanStringSingleQuotes(t *testing.T) {
	s := newScanner("'hello'")
	s.advance()
	s.stringSingleQuotes()
	assert.Equal(t, tokenString, s.tokens[0].Type)
	assert.Equal(t, "hello", s.tokens[0].Literal)
	assert.Equal(t, "'hello'", s.tokens[0].Lexeme)

	s = newScanner("'hello\nWorld'")
	s.advance()
	s.stringDoubleQuotes()
	assert.Equal(t, 2, s.line)
}

func TestScanNumber(t *testing.T) {
	s := newScanner("123")
	s.number()
	assert.Equal(t, "123", s.tokens[0].Lexeme)
	assert.Equal(t, tokenNumber, s.tokens[0].Type)
	assert.Equal(t, float64(123), s.tokens[0].Literal)
}

func TestScanTokenParenthesis(t *testing.T) {
	s := newScanner("(")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenLeftParenthesis, s.tokens[0].Type)
	assert.Equal(t, 1, s.tokens[0].Line)

	s = newScanner(")")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenRightParenthesis, s.tokens[0].Type)

}

func TestScanTokenBracket(t *testing.T) {
	s := newScanner("[")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenLeftBracket, s.tokens[0].Type)

	s = newScanner("]")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenRightBracket, s.tokens[0].Type)
}

func TestScanTokenArithmetic(t *testing.T) {
	s := newScanner("+")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenPlus, s.tokens[0].Type)

	s = newScanner("-")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenMinus, s.tokens[0].Type)

	s = newScanner("/")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenSlash, s.tokens[0].Type)

	s = newScanner("*")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenStar, s.tokens[0].Type)
}

func TestScanTokenBang(t *testing.T) {
	s := newScanner("!")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenBang, s.tokens[0].Type)

	s = newScanner("!=")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenBangEqual, s.tokens[0].Type)

}

func TestScanTokenEqual(t *testing.T) {
	s := newScanner("=")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenEqual, s.tokens[0].Type)

	s = newScanner("==")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenEqualEqual, s.tokens[0].Type)
}

func TestScanTokenComparison(t *testing.T) {
	s := newScanner(">")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenGreater, s.tokens[0].Type)

	s = newScanner(">=")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenGreaterEqual, s.tokens[0].Type)

	s = newScanner("<")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenLess, s.tokens[0].Type)

	s = newScanner("<=")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenLessEqual, s.tokens[0].Type)
}

func TestScantokenNumber(t *testing.T) {
	s := newScanner("1")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenNumber, s.tokens[0].Type)

	s = newScanner("1.10")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenNumber, s.tokens[0].Type)
}

func TestScantokenIdentifier(t *testing.T) {
	s := newScanner("a")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenIdentifier, s.tokens[0].Type)
}

func TestScantokenString(t *testing.T) {
	s := newScanner("\"hello\"")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenString, s.tokens[0].Type)
}

func TestScanComments(t *testing.T) {
	s := newScanner("//a")
	s.scanToken()
	assert.Len(t, s.tokens, 0)
}

func TestScanTokenDot(t *testing.T) {
	s := newScanner(".")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenDot, s.tokens[0].Type)

}

func TestScanTokenComma(t *testing.T) {
	s := newScanner(",")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenComma, s.tokens[0].Type)
}

func TestScanTokenSemiColon(t *testing.T) {
	s := newScanner(";")
	s.scanToken()
	assert.Len(t, s.tokens, 1)
	assert.Equal(t, tokenSemiColon, s.tokens[0].Type)
}

func TestScanTokenWhitespace(t *testing.T) {
	s := newScanner(" ")
	s.scanToken()
	assert.Len(t, s.tokens, 0)

	s = newScanner("\r")
	s.scanToken()
	assert.Len(t, s.tokens, 0)

	s = newScanner("\t")
	s.scanToken()
	assert.Len(t, s.tokens, 0)
}

func TestScanTokenBreakline(t *testing.T) {
	s := newScanner("\n")
	s.scanToken()
	assert.Len(t, s.tokens, 0)
	assert.Equal(t, 2, s.line)
}

func TestScanTokenUnexpected(t *testing.T) {
	s := newScanner("°")
	err := s.scanToken
	assert.NotNil(t, err)
}

func TestScanMultipleTokens(t *testing.T) {
	s := newScanner("print 2+2")
	tokens, err := s.scanTokens()
	assert.Nil(t, err)
	assert.Len(t, tokens, 5)
	assert.Equal(t, tokenNumber, tokens[1].Type)
	assert.Equal(t, float64(2), tokens[1].Literal)
	assert.Equal(t, tokenPlus, tokens[2].Type)
	assert.Equal(t, tokenNumber, tokens[3].Type)
	assert.Equal(t, float64(2), tokens[3].Literal)
	assert.Equal(t, tokenEndOfFile, tokens[4].Type)
}
