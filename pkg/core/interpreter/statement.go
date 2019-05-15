package main

import (
	"fmt"
)

type statement interface {
	evaluate(sc *Scope) (interface{}, error)
}

type expressionStmt struct {
	exp expression
}

func (stmt expressionStmt) evaluate(sc *Scope) (interface{}, error) {
	return stmt.exp.evaluate(sc)
}

type printStmt struct {
	exp expression
}

func (stmt printStmt) evaluate(sc *Scope) (interface{}, error) {
	value, err := stmt.exp.evaluate(sc)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", value)
	return nil, nil
}

type blockStmt struct {
	statements []statement
}

func (stmt blockStmt) evaluate(sc *Scope) (interface{}, error) {
	newscironment := NewScope(sc)

	for _, st := range stmt.statements {
		switch st.(type) {
		case returnStatement:
			val, err := st.evaluate(newscironment)
			if err != nil {
				return nil, err
			}
			return val, nil
		default:
			if _, err := st.evaluate(newscironment); err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

type ifStatement struct {
	cond     expression
	thenStmt statement
	elseStmt statement
}

func (stmt ifStatement) evaluate(sc *Scope) (interface{}, error) {
	cond, err := stmt.cond.evaluate(sc)
	if err != nil {
		return nil, err
	}
	if isTruthy(cond) {
		if _, err := stmt.thenStmt.evaluate(sc); err != nil {
			return nil, err
		}
	} else {
		if stmt.elseStmt != nil {
			elseStmt := stmt.elseStmt
			elseStmt.evaluate(sc)
		}
	}
	return nil, nil
}

type whileStatement struct {
	cond expression
	body statement
}

func (stmt whileStatement) evaluate(sc *Scope) (interface{}, error) {
	for {
		val, err := stmt.cond.evaluate(sc)
		if err != nil {
			return nil, err
		}
		if !isTruthy(val) {
			break
		}
		stmt.body.evaluate(sc)
	}
	return nil, nil
}

type funcStatement struct {
	name   token
	params []token
	body   blockStmt
}

func (stmt funcStatement) evaluate(sc *Scope) (interface{}, error) {
	f := function{
		declaration: stmt,
	}
	sc.SetValue(stmt.name.Lexeme, f)
	return nil, nil
}

type returnStatement struct {
	value expression
}

func (stmt returnStatement) evaluate(sc *Scope) (interface{}, error) {
	value, err := stmt.value.evaluate(sc)
	if err != nil {
		return nil, err
	}
	if value != nil {
		//Back to the top of the stack on the call statement
		//To ensure this, we panic to make like handled exception
		panic(value)
	}
	return nil, nil
}

func isTruthy(val interface{}) bool {
	if val == nil {
		return false
	}
	switch val.(type) {
	case bool:
		return val.(bool)
	}
	return true
}
