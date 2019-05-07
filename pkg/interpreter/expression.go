package interpreter

import (
	"errors"
	"fmt"
	"reflect"
)

type expression interface {
	evaluate(*Scope) (interface{}, error)
}

//Variable assignation
type assignExpression struct {
	op  token
	exp expression
}

func (e assignExpression) evaluate(sc *Scope) (interface{}, error) {
	value, err := e.exp.evaluate(sc)
	if err != nil {
		return nil, err
	}

	sc.SetValue(e.op.Lexeme, value)
	return nil, nil
}

//Variable execution
type variableExpression struct {
	op token
}

func (e variableExpression) evaluate(sc *Scope) (interface{}, error) {
	return sc.GetValue(e.op.Lexeme)
}

//Arithmetic (+ - * /) and logic (== !=  > < >= <=)
type binaryExpression struct {
	left  expression
	right expression
	op    token
}

func (e binaryExpression) evaluate(sc *Scope) (interface{}, error) {
	left, err := e.left.evaluate(sc)
	if err != nil {
		return nil, err
	}
	right, err := e.right.evaluate(sc)
	if err != nil {
		return nil, err
	}

	switch e.op.Type {
	case tokenMinus:
		return left.(float64) - right.(float64), nil
	case tokenSlash:
		return left.(float64) / right.(float64), nil
	case tokenStar:
		return left.(float64) * right.(float64), nil
	case tokenPlus:
		if reflect.TypeOf(left).String() == "float64" && reflect.TypeOf(right).String() == "float64" {
			return left.(float64) + right.(float64), nil
		}
		return fmt.Sprintf("%v%v", left, right), nil
	case tokenGreater:
		return left.(float64) > right.(float64), nil
	case tokenGreaterEqual:
		return left.(float64) >= right.(float64), nil
	case tokenLess:
		return left.(float64) < right.(float64), nil
	case tokenLessEqual:
		return left.(float64) <= right.(float64), nil
	case tokenEqualEqual:
		return left == right, nil
	case tokenBangEqual:
		return left != right, nil
	default:
		return nil, errors.New("Not supported as binary expression")
	}
}

//Parenthesis and brackets
type groupingExpression struct {
	exp expression
}

func (e groupingExpression) evaluate(sc *Scope) (interface{}, error) {
	return e.exp.evaluate(sc)
}

//Not expression or negative one
type unaryExpression struct {
	op    token
	right expression
}

func (e unaryExpression) evaluate(sc *Scope) (interface{}, error) {
	right, err := e.right.evaluate(sc)
	if err != nil {
		return nil, err
	}
	switch e.op.Type {
	case tokenBang:
		return !isTruthy(right), nil
	case tokenMinus:
		return -right.(float64), nil
	}

	return nil, nil
}

//Number, string, booleans
type literalExpression struct {
	value interface{}
}

func (e literalExpression) evaluate(sc *Scope) (interface{}, error) {
	return e.value, nil
}

//And, OR
type logicalExpression struct {
	left  expression
	op    token
	right expression
}

func (e logicalExpression) evaluate(sc *Scope) (interface{}, error) {
	left, err := e.left.evaluate(sc)
	if err != nil {
		return nil, err
	}
	if e.op.Type == tokenOr {
		if isTruthy(left) {
			return left, nil
		}
	} else {
		if !isTruthy(left) {
			return left, nil
		}
	}
	return e.right.evaluate(sc)
}

type stdCallExpression struct {
	callee expression
	args   map[string]expression
}

func (e stdCallExpression) evaluate(sc *Scope) (interface{}, error) {
	call, err := e.callee.evaluate(sc)
	if err != nil {
		return nil, err
	}

	f, ok := call.(func(*Scope, map[string]interface{}) (interface{}, error))
	if !ok {
		return nil, errors.New("Can only call functions")
	}

	args := make(map[string]interface{}, 0)
	for k, arg := range e.args {
		val, err := arg.evaluate(sc)
		if err != nil {
			return nil, err
		}
		args[k] = val
	}
	return f(sc, args)
}
