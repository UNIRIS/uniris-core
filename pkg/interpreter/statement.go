package interpreter

type statement interface {
	evaluate(sc *Scope) (interface{}, error)
}

type expressionStmt struct {
	exp expression
}

func (stmt expressionStmt) evaluate(sc *Scope) (interface{}, error) {
	return stmt.exp.evaluate(sc)
}

type blockStmt struct {
	statements []statement
}

func (stmt blockStmt) evaluate(sc *Scope) (interface{}, error) {
	nsc := NewScope(sc)

	out := make([]interface{}, 0)
	for _, st := range stmt.statements {
		res, err := st.evaluate(nsc)
		if err != nil {
			return nil, err
		}
		out = append(out, res)
	}

	return out, nil
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
