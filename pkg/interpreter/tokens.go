package interpreter

type token struct {
	Type    tokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

type tokenType string

var keywords = map[string]tokenType{
	"triggers":     tokenTriggers,
	"time":         tokenTriggerTime,
	"conditions":   tokenConditions,
	"originFamily": tokenOriginFamily,
	"response":     tokenResponseConditions,
	"postPaidFee":  tokenPostPaidFeeConditions,
	"inherit":      tokenInheritConditions,
	"actions":      tokenActions,

	"if":    tokenIf,
	"else":  tokenElse,
	"or":    tokenOr,
	"and":   tokenAnd,
	"true":  tokenTrue,
	"false": tokenFalse,
	"then":  tokenThen,
	"end":   tokenEnd,
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

	//----------
	//Keywords
	//----------------

	//Triggers
	tokenTriggers    tokenType = "TRIGGERS"
	tokenTriggerTime tokenType = "TRIGGER_TIME"

	//Conditions
	tokenConditions            tokenType = "CONDITIONS"
	tokenResponseConditions    tokenType = "RESPONSE"
	tokenOriginFamily          tokenType = "ORIGIN_FAMILY"
	tokenPostPaidFeeConditions tokenType = "POST_PAID_FEE"
	tokenInheritConditions     tokenType = "INHERIT"

	tokenActions tokenType = "ACTIONS"

	tokenIf        tokenType = "IF"
	tokenElse      tokenType = "ELSE"
	tokenAnd       tokenType = "AND"
	tokenOr        tokenType = "OR"
	tokenEndOfFile tokenType = "EOF"
	tokenTrue      tokenType = "TRUE"
	tokenFalse     tokenType = "FALSE"
	tokenThen      tokenType = "THEN"
	tokenEnd       tokenType = "END"
)
