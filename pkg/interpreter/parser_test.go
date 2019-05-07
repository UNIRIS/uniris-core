package interpreter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserPreviousToken(t *testing.T) {
	p := parser{
		tokens: []token{
			token{
				Type: tokenNumber,
			},
			token{
				Type: tokenPlus,
			},
			token{
				Type: tokenNumber,
			},
		},
	}

	p.current = 2

	assert.Equal(t, tokenPlus, p.previous().Type)
}

func TestParserPeek(t *testing.T) {
	p := parser{
		tokens: []token{
			token{
				Type: tokenNumber,
			},
		},
	}

	assert.Equal(t, tokenNumber, p.peek().Type)
}

func TestParserIsAtEnd(t *testing.T) {
	p := parser{
		tokens: []token{
			token{
				Type: tokenEndOfFile,
			},
		},
	}

	assert.True(t, p.isAtEnd())
}

func TestParserAdvance(t *testing.T) {
	p := parser{
		tokens: []token{
			token{
				Type: tokenNumber,
			},
			token{
				Type: tokenPlus,
			},
			token{
				Type: tokenNumber,
			},
		},
	}

	assert.Equal(t, tokenNumber, p.advance().Type)
	assert.Equal(t, tokenPlus, p.advance().Type)
	assert.Equal(t, tokenNumber, p.advance().Type)
}

func TestParserCheck(t *testing.T) {
	p := parser{
		tokens: []token{
			token{
				Type: tokenNumber,
			},
		},
	}

	assert.True(t, p.check(tokenNumber))
}

func TestParserError(t *testing.T) {
	p := parser{}
	assert.Error(t, p.error(token{Type: tokenEndOfFile, Line: 1}, "Invalid"), "Parsing error at end of line 1 - Invalid")
	assert.Error(t, p.error(token{Type: tokenNumber, Line: 1, Lexeme: "2"}, "Invalid"), "Parsing error at 2 of line 1 - Invalid")
}

func TestParserConsume(t *testing.T) {
	p := parser{
		tokens: []token{
			token{
				Type: tokenNumber,
			},
			token{
				Type: tokenPlus,
			},
			token{
				Type: tokenNumber,
			},
		},
	}

	token, err := p.consume(tokenNumber, "Invalid number")
	assert.Equal(t, tokenNumber, token.Type)
	assert.Nil(t, err)

	token, err = p.consume(tokenNumber, "Invalid Operator")
	assert.NotNil(t, err)
}

func TestParserMatch(t *testing.T) {
	p := parser{
		tokens: []token{
			token{
				Type: tokenNumber,
			},
			token{
				Type: tokenPlus,
			},
			token{
				Type: tokenNumber,
			},
		},
	}

	assert.True(t, p.match(tokenNumber, tokenPlus))

	assert.False(t, p.match(tokenIf))
}

func TestParserPrimaryLiteralExpression(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: tokenFalse},
			token{Type: tokenTrue},
			token{Type: tokenNumber, Literal: "10"},
			token{Type: tokenString, Literal: "hello"},
		},
	}

	exp, err := p.primary()
	assert.Nil(t, err)
	assert.Equal(t, literalExpression{value: false}, exp)

	exp, err = p.primary()
	assert.Nil(t, err)
	assert.Equal(t, literalExpression{value: true}, exp)

	exp, err = p.primary()
	assert.Nil(t, err)
	assert.Equal(t, literalExpression{value: "10"}, exp)

	exp, err = p.primary()
	assert.Nil(t, err)
	assert.Equal(t, literalExpression{value: "hello"}, exp)
}

func TestParserPrimaryAssignExpression(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: tokenIdentifier},
			token{Type: tokenEqual},
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenEndOfFile},
		},
	}

	exp, err := p.primary()
	assert.Nil(t, err)
	assert.Equal(t, assignExpression{
		op: token{
			Type: tokenIdentifier,
		},
		exp: literalExpression{
			value: 10,
		},
	}, exp)
}

func TestParserPrimaryVariableExpression(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: tokenIdentifier},
			token{Type: tokenEndOfFile},
		},
	}

	exp, err := p.primary()
	assert.Nil(t, err)
	assert.Equal(t, variableExpression{
		op: token{
			Type: tokenIdentifier,
		},
	}, exp)
}

func TestParserPrimaryGroupingExpression(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: tokenLeftParenthesis},
			token{Type: tokenTrue},
			token{Type: tokenRightParenthesis},
			token{Type: tokenEndOfFile},
		},
	}

	exp, err := p.primary()
	assert.Nil(t, err)
	assert.Equal(t, groupingExpression{
		exp: literalExpression{
			value: true,
		},
	}, exp)
}

func TestParserFinishCall(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: tokenIdentifier, Lexeme: "val"},
			token{Type: tokenColon},
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenRightParenthesis},
			token{Type: tokenEndOfFile},
		},
	}
	exp, err := p.finishStdCall(testFuncExpression{})
	assert.Nil(t, err)
	assert.Equal(t, stdCallExpression{
		args: map[string]expression{
			"val": literalExpression{
				value: 10,
			},
		},
		callee: testFuncExpression{},
	}, exp)
}

func TestParserCall(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: tokenIdentifier},
			token{Type: tokenLeftParenthesis},
			token{Type: tokenRightParenthesis},
			token{Type: tokenEndOfFile},
		},
	}

	exp, err := p.stdCall()
	assert.Nil(t, err)
	assert.Equal(t, stdCallExpression{
		args: map[string]expression{},
		callee: variableExpression{
			op: token{Type: tokenIdentifier},
		},
	}, exp)
}

func TestParserUnary(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: tokenBang},
			token{Type: tokenTrue},
			token{Type: tokenMinus},
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenEndOfFile},
		},
	}

	exp, err := p.unary()
	assert.Nil(t, err)
	assert.Equal(t, unaryExpression{
		op: token{Type: tokenBang},
		right: literalExpression{
			value: true,
		},
	}, exp)

	exp, err = p.unary()
	assert.Nil(t, err)
	assert.Equal(t, unaryExpression{
		op: token{Type: tokenMinus},
		right: literalExpression{
			value: 10,
		},
	}, exp)
}

func TestParserMultiplication(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenStar},
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenEndOfFile},
		},
	}

	exp, err := p.multiplication()
	assert.Nil(t, err)
	assert.Equal(t, binaryExpression{
		left:  literalExpression{value: 10},
		op:    token{Type: tokenStar},
		right: literalExpression{value: 10},
	}, exp)
}

func TestParserAddition(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenPlus},
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenMinus},
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenEndOfFile},
		},
	}

	exp, err := p.addition()
	assert.Nil(t, err)
	assert.Equal(t, binaryExpression{
		left:  literalExpression{value: 10},
		op:    token{Type: tokenPlus},
		right: literalExpression{value: 10},
	}, exp)

	exp, err = p.addition()
	assert.Nil(t, err)
	assert.Equal(t, binaryExpression{
		left:  literalExpression{value: 10},
		op:    token{Type: tokenMinus},
		right: literalExpression{value: 10},
	}, exp)
}

func TestParserComparison(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenGreater},
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenGreaterEqual},
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenLess},
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenLessEqual},
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenEndOfFile},
		},
	}

	exp, err := p.comparison()
	assert.Nil(t, err)
	assert.Equal(t, binaryExpression{
		left:  literalExpression{value: 10},
		op:    token{Type: tokenGreater},
		right: literalExpression{value: 10},
	}, exp)

	exp, err = p.comparison()
	assert.Nil(t, err)
	assert.Equal(t, binaryExpression{
		left:  literalExpression{value: 10},
		op:    token{Type: tokenGreaterEqual},
		right: literalExpression{value: 10},
	}, exp)

	exp, err = p.comparison()
	assert.Nil(t, err)
	assert.Equal(t, binaryExpression{
		left:  literalExpression{value: 10},
		op:    token{Type: tokenLess},
		right: literalExpression{value: 10},
	}, exp)

	exp, err = p.comparison()
	assert.Nil(t, err)
	assert.Equal(t, binaryExpression{
		left:  literalExpression{value: 10},
		op:    token{Type: tokenLessEqual},
		right: literalExpression{value: 10},
	}, exp)
}

func TestParserEquality(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenEqualEqual},
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenIdentifier},
			token{Type: tokenBangEqual},
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenEndOfFile},
		},
	}

	exp, err := p.equality()
	assert.Nil(t, err)
	assert.Equal(t, binaryExpression{
		left:  literalExpression{value: 10},
		op:    token{Type: tokenEqualEqual},
		right: literalExpression{value: 10},
	}, exp)

	exp, err = p.equality()
	assert.Nil(t, err)
	assert.Equal(t, binaryExpression{
		left:  variableExpression{op: token{Type: tokenIdentifier}},
		op:    token{Type: tokenBangEqual},
		right: literalExpression{value: 10},
	}, exp)
}

func TestParserAnd(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: tokenTrue},
			token{Type: tokenAnd},
			token{Type: tokenTrue},
			token{Type: tokenEndOfFile},
		},
	}

	exp, err := p.and()
	assert.Nil(t, err)
	assert.Equal(t, logicalExpression{
		left:  literalExpression{value: true},
		op:    token{Type: tokenAnd},
		right: literalExpression{value: true},
	}, exp)
}

func TestParserOr(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: tokenTrue},
			token{Type: tokenOr},
			token{Type: tokenTrue},
			token{Type: tokenEndOfFile},
		},
	}

	exp, err := p.or()
	assert.Nil(t, err)
	assert.Equal(t, logicalExpression{
		left:  literalExpression{value: true},
		op:    token{Type: tokenOr},
		right: literalExpression{value: true},
	}, exp)
}

func TestParserAssignement(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: tokenIdentifier},
			token{Type: tokenEqual},
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenEndOfFile},
		},
	}

	exp, err := p.assignement()
	assert.Nil(t, err)
	assert.Equal(t, assignExpression{
		exp: literalExpression{value: 10},
		op:  token{Type: tokenIdentifier},
	}, exp)
}

func TestParserExpressionStatement(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenEndOfFile},
		},
	}

	stmt, err := p.expressionStatement()
	assert.Nil(t, err)
	assert.Equal(t, expressionStmt{
		exp: literalExpression{value: 10},
	}, stmt)

}

func TestParserBlockStatement(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenPlus},
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenEnd},
			token{Type: tokenEndOfFile},
		},
	}

	stmt, err := p.blockStatements()
	assert.Nil(t, err)
	assert.Equal(t, blockStmt{
		statements: []statement{
			expressionStmt{
				binaryExpression{
					left:  literalExpression{value: 10},
					right: literalExpression{value: 10},
					op:    token{Type: tokenPlus},
				},
			},
		},
	}, stmt)
}

func TestParserIfStatement(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: tokenTrue},
			token{Type: tokenNumber},
			token{Type: tokenElse},
			token{Type: tokenNumber},
			token{Type: tokenEndOfFile},
		},
	}

	stmt, err := p.ifStatement()
	assert.Nil(t, err)
	assert.Equal(t, ifStatement{
		cond:     literalExpression{value: true},
		thenStmt: expressionStmt{exp: literalExpression{}},
		elseStmt: expressionStmt{exp: literalExpression{}},
	}, stmt)
}

func TestParserStatementWhenMatchesIf(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: tokenIf},
			token{Type: tokenTrue},
			token{Type: tokenThen},
			token{Type: tokenEnd},
			token{Type: tokenEndOfFile},
		},
	}

	stmt, err := p.statement()
	assert.Nil(t, err)
	assert.Equal(t, ifStatement{
		cond: literalExpression{value: true},
		thenStmt: blockStmt{
			statements: []statement{},
		},
		elseStmt: nil,
	}, stmt)
}

func TestParserStatementWhenMatchesThenBrace(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: tokenThen},
			token{Type: tokenEnd},
			token{Type: tokenEndOfFile},
		},
	}

	stmt, err := p.statement()
	assert.Nil(t, err)
	assert.Equal(t, blockStmt{
		statements: []statement{},
	}, stmt)
}

func TestParse(t *testing.T) {
	p := parser{
		tokens: []token{
			token{Type: tokenActions},
			token{Type: tokenColon},
			token{Type: tokenNumber, Literal: 10},
			token{Type: tokenPlus},
			token{Type: tokenNumber, Literal: 2},
			token{Type: tokenEndOfFile},
		},
	}
	contract, err := p.parse()
	assert.Nil(t, err)
	assert.Equal(t, []statement{
		expressionStmt{
			binaryExpression{
				left:  literalExpression{value: 10},
				right: literalExpression{value: 2},
				op:    token{Type: tokenPlus},
			},
		},
	}, contract.actions)
}
