package interpreter

import (
	"fmt"
)

//Execute runs the given code by scanning, parsing and evaluating its statements
func Execute(code string, s *Scope) (string, error) {
	globals := NewScope(nil)
	globals.SetValue("timestamp", timestampFunc{})

	if s == nil {
		s = NewScope(nil)
	}
	s.parent = globals

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
		val, err := st.evaluate(s)
		if err != nil {
			return "", err
		}
		if val != nil {
			res += fmt.Sprintf("%v\n", val)
		}
	}

	return res, nil
}
