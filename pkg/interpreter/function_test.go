package interpreter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctionCallWithoutParameters(t *testing.T) {

	f := function{
		declaration: funcStatement{
			name: token{Type: tokenIdentifier, Lexeme: "myFunc"},
			body: blockStmt{
				statements: []statement{
					returnStatement{value: literalExpression{value: 3}},
				},
			},
		},
	}

	res, err := f.call(nil)
	assert.Nil(t, err)
	assert.Equal(t, 3, res)
}

func TestFunctionCallWithArgsNotAllowed(t *testing.T) {

	f := function{
		declaration: funcStatement{
			name: token{Type: tokenIdentifier, Lexeme: "myFunc"},
			body: blockStmt{
				statements: []statement{
					returnStatement{value: literalExpression{value: 3}},
				},
			},
		},
	}

	_, err := f.call(nil, 1, 2, 3)
	assert.EqualError(t, err, "no parameters allowed for this function")
}

func TestFunctionCallWithMissingArgs(t *testing.T) {

	f := function{
		declaration: funcStatement{
			name: token{Type: tokenIdentifier, Lexeme: "myFunc"},
			params: []token{
				token{Type: tokenIdentifier, Lexeme: "text"},
			},
			body: blockStmt{
				statements: []statement{
					returnStatement{value: literalExpression{value: 3}},
				},
			},
		},
	}

	_, err := f.call(nil)
	assert.EqualError(t, err, "missing function parameters")
}

func TestFunctionCallWithParameters(t *testing.T) {

	f := function{
		declaration: funcStatement{
			name: token{Type: tokenIdentifier, Lexeme: "myFunc"},
			params: []token{
				token{Type: tokenIdentifier, Lexeme: "text"},
			},
			body: blockStmt{
				statements: []statement{
					returnStatement{
						value: variableExpression{
							op: token{Type: tokenIdentifier, Lexeme: "text"},
						},
					},
				},
			},
		},
	}

	res, err := f.call(nil, "hello")
	assert.Nil(t, err)
	assert.Equal(t, "hello", res)
}
