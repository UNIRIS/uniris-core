package interpreter

import (
	"fmt"
)

//Analyse checks if the code is correct
func Analyse(code string) (Contract, error) {
	scan := newScanner(code)
	tokens, err := scan.scanTokens()
	if err != nil {
		return Contract{}, err
	}

	p := parser{
		tokens: tokens,
	}
	c, err := p.parse()
	if err != nil {
		return Contract{}, err
	}

	return c, c.analyze()
}

//Execute runs the given code by scanning, parsing and evaluating its statements
func Execute(code string, s *Scope) (string, error) {
	globals := NewScope(nil)
	for n, f := range stdFunctions {
		globals.SetValue(n, f)
	}

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
	contract, err := p.parse()
	if err != nil {
		return "", err
	}

	res := ""

	val, err := contract.execute(s)
	if err != nil {
		return "", err
	}
	if val != nil {
		switch val.(type) {
		case []interface{}:
			for _, v := range val.([]interface{}) {
				res += fmt.Sprintf("%v\n", v)
			}
			break
		default:
			res += fmt.Sprintf("%v\n", val)
		}
	}

	return res, nil
}
