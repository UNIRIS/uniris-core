package main

import (
	"errors"
)

type callable interface {
	call(*Scope, ...interface{}) (interface{}, error)
}

type function struct {
	declaration funcStatement
}

func (f function) call(sc *Scope, args ...interface{}) (res interface{}, err error) {
	newScope := NewScope(sc)

	if len(f.declaration.params) == 0 && len(args) > 0 {
		return nil, errors.New("no parameters allowed for this function")
	}

	if len(args) != len(f.declaration.params) {
		return nil, errors.New("missing function parameters")
	}

	for i := 0; i < len(f.declaration.params); i++ {
		newScope.SetValue(f.declaration.params[i].Lexeme, args[i])
	}

	defer func() {
		if x := recover(); x != nil {
			res = x
		}
	}()
	res, err = f.declaration.body.evaluate(newScope)
	if err != nil {
		return nil, err
	}
	return res, nil
}
