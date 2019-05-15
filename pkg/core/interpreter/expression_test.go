package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLiteralExpression(t *testing.T) {
	e := literalExpression{
		value: 10,
	}
	val, err := e.evaluate(NewScope(nil))
	assert.Nil(t, err)
	assert.Equal(t, 10, val)
}

func TestAssignExpression(t *testing.T) {
	e := assignExpression{
		exp: literalExpression{
			value: 10,
		},
		op: token{
			Lexeme: "a",
		},
	}

	env := NewScope(nil)
	_, err := e.evaluate(env)
	assert.Nil(t, err)
	val, err := env.GetValue("a")
	assert.Nil(t, err)
	assert.Equal(t, 10, val)
}

func TestVariableExpression(t *testing.T) {

	env := NewScope(nil)
	env.SetValue("a", 10)

	e := variableExpression{
		op: token{
			Lexeme: "a",
		},
	}

	val, err := e.evaluate(env)
	assert.Nil(t, err)
	assert.Equal(t, 10, val)
}

func TestBinaryMinusExpression(t *testing.T) {
	e := binaryExpression{
		left: literalExpression{
			value: float64(10),
		},
		right: literalExpression{
			value: float64(10),
		},
		op: token{Type: tokenMinus},
	}

	val, err := e.evaluate(NewScope(nil))
	assert.Nil(t, err)
	assert.Equal(t, float64(0), val)
}

func TestBinaryPlusExpression(t *testing.T) {
	e := binaryExpression{
		left: literalExpression{
			value: float64(10),
		},
		right: literalExpression{
			value: float64(10),
		},
		op: token{Type: tokenPlus},
	}

	val, err := e.evaluate(NewScope(nil))
	assert.Nil(t, err)
	assert.Equal(t, float64(20), val)

	e = binaryExpression{
		left: literalExpression{
			value: "hello ",
		},
		right: literalExpression{
			value: "world",
		},
		op: token{Type: tokenPlus},
	}
	val, err = e.evaluate(NewScope(nil))
	assert.Nil(t, err)
	assert.Equal(t, "hello world", val)
}

func TestBinaryStarExpression(t *testing.T) {
	e := binaryExpression{
		left: literalExpression{
			value: float64(10),
		},
		right: literalExpression{
			value: float64(10),
		},
		op: token{Type: tokenStar},
	}

	val, err := e.evaluate(NewScope(nil))
	assert.Nil(t, err)
	assert.Equal(t, float64(100), val)
}

func TestBinarySlashExpression(t *testing.T) {
	e := binaryExpression{
		left: literalExpression{
			value: float64(10),
		},
		right: literalExpression{
			value: float64(10),
		},
		op: token{Type: tokenSlash},
	}

	val, err := e.evaluate(NewScope(nil))
	assert.Nil(t, err)
	assert.Equal(t, float64(1), val)
}

func TestBinaryGreaterExpression(t *testing.T) {
	e := binaryExpression{
		left: literalExpression{
			value: float64(11),
		},
		right: literalExpression{
			value: float64(10),
		},
		op: token{Type: tokenGreater},
	}

	val, err := e.evaluate(NewScope(nil))
	assert.Nil(t, err)
	assert.Equal(t, true, val)
}

func TestBinaryGreaterEqualExpression(t *testing.T) {
	e := binaryExpression{
		left: literalExpression{
			value: float64(10),
		},
		right: literalExpression{
			value: float64(10),
		},
		op: token{Type: tokenGreaterEqual},
	}

	val, err := e.evaluate(NewScope(nil))
	assert.Nil(t, err)
	assert.Equal(t, true, val)
}

func TestBinaryLessExpression(t *testing.T) {
	e := binaryExpression{
		left: literalExpression{
			value: float64(10),
		},
		right: literalExpression{
			value: float64(10),
		},
		op: token{Type: tokenLess},
	}

	val, err := e.evaluate(NewScope(nil))
	assert.Nil(t, err)
	assert.Equal(t, false, val)
}

func TestBinaryLessEqualExpression(t *testing.T) {
	e := binaryExpression{
		left: literalExpression{
			value: float64(10),
		},
		right: literalExpression{
			value: float64(10),
		},
		op: token{Type: tokenLessEqual},
	}

	val, err := e.evaluate(NewScope(nil))
	assert.Nil(t, err)
	assert.Equal(t, true, val)
}

func TestBinaryEqualEqualExpression(t *testing.T) {
	e := binaryExpression{
		left: literalExpression{
			value: float64(10),
		},
		right: literalExpression{
			value: float64(10),
		},
		op: token{Type: tokenEqualEqual},
	}

	val, err := e.evaluate(NewScope(nil))
	assert.Nil(t, err)
	assert.Equal(t, true, val)
}

func TestBinaryBangEqualExpression(t *testing.T) {
	e := binaryExpression{
		left: literalExpression{
			value: float64(10),
		},
		right: literalExpression{
			value: float64(10),
		},
		op: token{Type: tokenBangEqual},
	}

	val, err := e.evaluate(NewScope(nil))
	assert.Nil(t, err)
	assert.Equal(t, false, val)
}

func TestGroupExpression(t *testing.T) {
	e := groupingExpression{
		exp: literalExpression{
			value: 10,
		},
	}

	val, err := e.evaluate(NewScope(nil))
	assert.Nil(t, err)
	assert.Equal(t, 10, val)
}

func TestUnaryBangExpression(t *testing.T) {
	e := unaryExpression{
		op: token{
			Type: tokenBang,
		},
		right: literalExpression{
			value: true,
		},
	}

	val, err := e.evaluate(NewScope(nil))
	assert.Nil(t, err)
	assert.Equal(t, false, val)

}

func TestUnaryMinusExpression(t *testing.T) {
	e := unaryExpression{
		op: token{
			Type: tokenMinus,
		},
		right: literalExpression{
			value: float64(10),
		},
	}

	val, err := e.evaluate(NewScope(nil))
	assert.Nil(t, err)
	assert.Equal(t, float64(-10), val)
}

func TestLogicalANDExpression(t *testing.T) {
	e := logicalExpression{
		left: literalExpression{
			value: true,
		},
		right: literalExpression{
			value: 10,
		},
		op: token{
			Type: tokenAnd,
		},
	}

	val, err := e.evaluate(NewScope(nil))
	assert.Nil(t, err)
	assert.Equal(t, 10, val)

	e = logicalExpression{
		left: literalExpression{
			value: false,
		},
		right: literalExpression{
			value: 10,
		},
		op: token{
			Type: tokenAnd,
		},
	}

	val, err = e.evaluate(NewScope(nil))
	assert.Nil(t, err)
	assert.Equal(t, false, val)
}

func TestLogicalORExpression(t *testing.T) {
	e := logicalExpression{
		left: literalExpression{
			value: true,
		},
		right: literalExpression{
			value: 10,
		},
		op: token{
			Type: tokenOr,
		},
	}

	val, err := e.evaluate(NewScope(nil))
	assert.Nil(t, err)
	assert.Equal(t, true, val)

	e = logicalExpression{
		left: literalExpression{
			value: false,
		},
		right: literalExpression{
			value: 10,
		},
		op: token{
			Type: tokenOr,
		},
	}

	val, err = e.evaluate(NewScope(nil))
	assert.Nil(t, err)
	assert.Equal(t, 10, val)
}

func TestCallExpression(t *testing.T) {
	e := callExpression{
		args: []expression{
			literalExpression{
				value: 10,
			},
		},
		callee: testFuncExpression{},
	}

	val, err := e.evaluate(NewScope(nil))
	assert.Nil(t, err)
	assert.Equal(t, 10, val)
}

func TestCallNotCallableExpression(t *testing.T) {
	e := callExpression{
		callee: literalExpression{},
	}

	_, err := e.evaluate(NewScope(nil))
	assert.Error(t, err, "Can only call functions")
}

type testFuncExpression struct {
}

func (f testFuncExpression) evaluate(sc *Scope) (interface{}, error) {
	return testFuncCallable{}, nil
}

type testFuncCallable struct{}

func (f testFuncCallable) call(sc *Scope, args ...interface{}) (res interface{}, err error) {
	return args[0], nil
}

func TestCollectionExpression(t *testing.T) {

	env := NewScope(nil)

	a := map[int]int{
		0: 1,
		1: 2,
		2: 3,
	}
	env.SetValue("a", a)
	e := collectionExpression{
		op: token{
			Lexeme: "a",
		},
		index: token{Literal: 2},
	}
	val, err := e.evaluate(env)
	assert.Nil(t, err)
	assert.Equal(t, 3, val)

	b := map[string]string{
		"hello": "world",
	}

	env.SetValue("b", b)
	e = collectionExpression{
		op: token{
			Lexeme: "b",
		},
		index: token{Literal: "hello"},
	}

	val, err = e.evaluate(env)
	assert.Nil(t, err)
	assert.Equal(t, "world", val)

	env.SetValue("c", func() {})
	e = collectionExpression{
		op: token{
			Lexeme: "c",
		},
		index: token{Literal: func() {}},
	}

	_, err = e.evaluate(env)
	assert.EqualError(t, err, "c has not the map type - required for the collections")
}

func TestAssignCollectionExpression(t *testing.T) {
	env := NewScope(nil)

	a := map[string]string{
		"hello": "world",
	}
	env.SetValue("a", a)
	e := collectionExpression{
		op: token{
			Lexeme: "a",
		},
		index: token{Literal: "hello"},
	}

	val, err := e.evaluate(env)
	assert.Nil(t, err)
	assert.Equal(t, "world", val)

	assign := collectionAssignmentExpression{
		op: token{
			Lexeme: "a",
		},
		index: token{Literal: "hello"},
		val:   literalExpression{value: "john"},
	}

	_, err = assign.evaluate(env)
	assert.Nil(t, err)

	e = collectionExpression{
		op: token{
			Lexeme: "a",
		},
		index: token{Literal: "hello"},
	}
	val, err = e.evaluate(env)
	assert.Nil(t, err)
	assert.Equal(t, "john", val)
}
