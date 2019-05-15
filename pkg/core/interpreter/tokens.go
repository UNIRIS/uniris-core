package main

type token struct {
	Type    tokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

type tokenType string

var keywords = map[string]tokenType{
	"if":          tokenIf,
	"else":        tokenElse,
	"while":       tokenWhile,
	"for":         tokenFor,
	"or":          tokenOr,
	"and":         tokenAnd,
	"true":        tokenTrue,
	"false":       tokenFalse,
	"function":    tokenFunction,
	"print":       tokenPrint,
	"return":      tokenReturn,
	"transaction": tokenTransaction,
	"contract":    tokenContract,
}

const (
	// Single-character tokens.
	tokenLeftParenthesis  tokenType = "LEFT_PARENTHESIS"
	tokenRightParenthesis tokenType = "RIGHT_PARENTHESIS"
	tokenLeftBrace        tokenType = "LEFT_BRACE"
	tokenRightBrace       tokenType = "RIGHT_BRACE"
	tokenLeftBracket      tokenType = "LEFT_BRACKET"
	tokenRightBracket     tokenType = "RIGHT_BRACKET"
	tokenPlus             tokenType = "PLUS"
	tokenMinus            tokenType = "MINUS"
	tokenStar             tokenType = "STAR"
	tokenSlash            tokenType = "SLASH"
	tokenDot              tokenType = "DOT"
	tokenComma            tokenType = "COMMA"
	tokenSemiColon        tokenType = "SEMICOLON"

	//One or two character tokens
	tokenBang         tokenType = "BANG"
	tokenBangEqual    tokenType = "BANQ_EQUAL"
	tokenEqual        tokenType = "EQUAL"
	tokenEqualEqual   tokenType = "EQUAL_EQUAL"
	tokenLess         tokenType = "LESS"
	tokenGreater      tokenType = "GREATER"
	tokenLessEqual    tokenType = "LESS_EQUAL"
	tokenGreaterEqual tokenType = "GREATER_EQUAL"

	//Literals
	tokenIdentifier tokenType = "IDENTIFIER"
	tokenString     tokenType = "STRING"
	tokenNumber     tokenType = "NUMBER"

	//Keywords
	tokenPrint       tokenType = "PRINT"
	tokenIf          tokenType = "IF"
	tokenElse        tokenType = "ELSE"
	tokenAnd         tokenType = "AND"
	tokenOr          tokenType = "OR"
	tokenWhile       tokenType = "WHILE"
	tokenFor         tokenType = "FOR"
	tokenEndOfFile   tokenType = "EOF"
	tokenTrue        tokenType = "TRUE"
	tokenFalse       tokenType = "FALSE"
	tokenFunction    tokenType = "FUNC"
	tokenReturn      tokenType = "RETURN"
	tokenTransaction tokenType = "TRANSACTION"
	tokenContract    tokenType = "CONTRACT"
)
