package interpreter

type token struct {
	Type    tokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

type tokenType string

var keywords = map[string]tokenType{
	"if":         tokenIf,
	"else":       tokenElse,
	"or":         tokenOr,
	"and":        tokenAnd,
	"true":       tokenTrue,
	"false":      tokenFalse,
	"with":       tokenWith,
	"then":       tokenThen,
	"end":        tokenEnd,
	"contract":   tokenContract,
	"conditions": tokenConditions,
	"answer":     tokenAnswerConditions,
	"fees":       tokenFeeConditions,
	"is":         tokenIs,
	"actions":    tokenActions,
	"triggers":   tokenTriggers,
	"time":       tokenTriggerTime,
}

const (
	// Single-character tokens.
	tokenLeftParenthesis  tokenType = "LEFT_PARENTHESIS"
	tokenRightParenthesis tokenType = "RIGHT_PARENTHESIS"
	tokenLeftBracket      tokenType = "LEFT_BRACKET"
	tokenRightBracket     tokenType = "RIGHT_BRACKET"
	tokenPlus             tokenType = "PLUS"
	tokenMinus            tokenType = "MINUS"
	tokenStar             tokenType = "STAR"
	tokenSlash            tokenType = "SLASH"
	tokenDot              tokenType = "DOT"
	tokenComma            tokenType = "COMMA"
	tokenSemiColon        tokenType = "SEMICOLON"
	tokenColon            tokenType = "COLON"

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
	tokenContract         tokenType = "CONTRACT"
	tokenIf               tokenType = "IF"
	tokenElse             tokenType = "ELSE"
	tokenAnd              tokenType = "AND"
	tokenOr               tokenType = "OR"
	tokenWith             tokenType = "WITH"
	tokenEndOfFile        tokenType = "EOF"
	tokenTrue             tokenType = "TRUE"
	tokenFalse            tokenType = "FALSE"
	tokenConditions       tokenType = "CONDITIONS"
	tokenAnswerConditions tokenType = "ANSWER"
	tokenFeeConditions    tokenType = "FEES"
	tokenIs               tokenType = "IS"
	tokenActions          tokenType = "ACTIONS"
	tokenTriggers         tokenType = "TRIGGERS"
	tokenTriggerTime      tokenType = "TRIGGER_TIME"
	tokenThen             tokenType = "THEN"
	tokenEnd              tokenType = "END"
)
