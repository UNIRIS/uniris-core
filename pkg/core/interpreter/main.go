package main

import "fmt"

//Execute runs the given code by scanning, parsing and evaluating its statements
func Execute(code string, scope interface{}) (string, error) {
	globals := NewScope(nil)
	globals.SetValue("timestamp", timestampFunc{})

	var sc *Scope
	if scope == nil {
		sc = NewScope(nil)
	} else {
		sc = scope.(*Scope)
	}
	sc.parent = globals

	scan := newScanner(code)
	tokens, err := scan.scanTokens()
	if err != nil {
		return "", err
	}
	p := parser{
		tokens: tokens,
	}
	stmt, err := p.parse()
	if err != nil {
		return "", err
	}

	res := ""

	for _, st := range stmt {
		val, err := st.evaluate(sc)
		if err != nil {
			return "", err
		}
		if val != nil {
			res += fmt.Sprintf("%v\n", val)
		}
	}

	return res, nil
}
