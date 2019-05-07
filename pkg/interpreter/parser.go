package interpreter

import (
	"fmt"
)

type parser struct {
	tokens  []token
	current int
}

func (p *parser) parse() (c Contract, err error) {

	if p.match(tokenTriggers) {
		if _, err := p.consume(tokenColon, "expected `:` after triggers"); err != nil {
			return c, err
		}
		c.Triggers = make([]Trigger, 0)
		if p.match(tokenTriggerTime) {
			if _, err := p.consume(tokenColon, "expected `:` after time"); err != nil {
				return c, err
			}
			val := p.peek().Literal
			c.Triggers = append(c.Triggers, Trigger{
				kind: timeTrigger,
				val:  val,
			})
			p.advance()
		}
	}

	if p.match(tokenConditions) {
		if _, err := p.consume(tokenColon, "expected `:` after conditions"); err != nil {
			return c, err
		}

		if p.match(tokenOriginFamily) {
			if _, err := p.consume(tokenColon, "expected `:` after originFamily conditions"); err != nil {
				return c, err
			}
			val, err := p.expression()
			if err != nil {
				return c, err
			}
			c.Conditions.OriginFamily = val
		}

		if p.match(tokenPostPaidFeeConditions) {
			if _, err := p.consume(tokenColon, "expected `:` after postPaidFee conditions"); err != nil {
				return c, err
			}
			val, err := p.expression()
			if err != nil {
				return c, err
			}
			c.Conditions.PostPaidFee = val
		}
		if p.match(tokenResponseConditions) {
			if _, err := p.consume(tokenColon, "expected `:` after response conditions"); err != nil {
				return c, err
			}
			val, err := p.expression()
			if err != nil {
				return c, err
			}
			c.Conditions.Response = val
		}
		if p.match(tokenInheritConditions) {
			if _, err := p.consume(tokenColon, "expected `:` after inherit conditions"); err != nil {
				return c, err
			}
			val, err := p.expression()
			if err != nil {
				return c, err
			}
			c.Conditions.Inherit = val
		}
	}

	if _, err := p.consume(tokenActions, "expected `actions` in the contract code"); err != nil {
		return c, err
	}
	if _, err := p.consume(tokenColon, "expected `:` after actions"); err != nil {
		return c, err
	}

	for !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			return c, err
		}
		c.actions = append(c.actions, stmt)
	}

	return c, nil
}

func (p *parser) statement() (statement, error) {

	if p.match(tokenIf) {
		return p.ifStatement()
	}
	if p.match(tokenThen) {
		return p.blockStatements()
	}
	return p.expressionStatement()
}

func (p *parser) ifStatement() (statement, error) {
	cond, err := p.expression()
	if err != nil {
		return nil, err
	}

	thenStmt, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseStmt statement

	if p.match(tokenElse) {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		elseStmt = stmt
	}

	return ifStatement{
		cond:     cond,
		thenStmt: thenStmt,
		elseStmt: elseStmt,
	}, nil
}

func (p *parser) blockStatements() (statement, error) {
	statements := make([]statement, 0)
	for !p.check(tokenEnd) && !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	if _, err := p.consume(tokenEnd, "Expect 'end' after block"); err != nil {
		return nil, err
	}
	return blockStmt{
		statements: statements,
	}, nil
}

func (p *parser) expressionStatement() (statement, error) {
	exp, err := p.expression()
	if err != nil {
		return nil, err
	}
	return expressionStmt{exp}, nil
}

func (p *parser) expression() (expression, error) {
	return p.assignement()
}

func (p *parser) assignement() (expression, error) {
	exp, err := p.or()
	if err != nil {
		return nil, err
	}
	if p.match(tokenEqual) {
		eq := p.previous()
		val, err := p.assignement()
		if err != nil {
			return nil, err
		}

		return assignExpression{
			op:  eq,
			exp: val,
		}, nil
	}

	return exp, nil
}

func (p *parser) or() (expression, error) {
	exp, err := p.and()
	if err != nil {
		return nil, err
	}
	for p.match(tokenOr) {
		op := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		exp = logicalExpression{
			left:  exp,
			right: right,
			op:    op,
		}
	}
	return exp, nil
}

func (p *parser) and() (expression, error) {
	exp, err := p.equality()
	if err != nil {
		return nil, err
	}
	for p.match(tokenAnd) {
		op := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		exp = logicalExpression{
			op:    op,
			left:  exp,
			right: right,
		}
	}
	return exp, nil
}

func (p *parser) equality() (expression, error) {
	exp, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match(tokenBangEqual, tokenEqualEqual) {
		op := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		exp = binaryExpression{
			left:  exp,
			op:    op,
			right: right,
		}
	}

	return exp, nil
}

func (p *parser) comparison() (expression, error) {
	exp, err := p.addition()
	if err != nil {
		return nil, err
	}
	for p.match(tokenGreater, tokenGreaterEqual, tokenLess, tokenLessEqual) {
		op := p.previous()
		right, err := p.addition()
		if err != nil {
			return nil, err
		}
		exp = binaryExpression{
			left:  exp,
			op:    op,
			right: right,
		}
	}
	return exp, nil
}

func (p *parser) addition() (expression, error) {
	exp, err := p.multiplication()
	if err != nil {
		return nil, err
	}
	for p.match(tokenMinus, tokenPlus) {
		op := p.previous()
		right, err := p.multiplication()
		if err != nil {
			return nil, err
		}
		exp = binaryExpression{
			left:  exp,
			op:    op,
			right: right,
		}
	}
	return exp, nil
}

func (p *parser) multiplication() (expression, error) {
	exp, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(tokenSlash, tokenStar) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		exp = binaryExpression{
			left:  exp,
			op:    op,
			right: right,
		}
	}
	return exp, nil
}

func (p *parser) unary() (expression, error) {
	if p.match(tokenBang, tokenMinus) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return unaryExpression{
			op:    op,
			right: right,
		}, nil
	}

	return p.stdCall()
}

func (p *parser) stdCall() (expression, error) {
	exp, err := p.primary()
	if err != nil {
		return nil, err
	}
	for true {
		if p.match(tokenLeftParenthesis) {
			switch exp.(type) {
			case variableExpression:
				v := exp.(variableExpression)
				if _, exist := stdFunctions[v.op.Lexeme]; !exist {
					return nil, fmt.Errorf("function '%s' undefined", v.op.Lexeme)
				}
			}
			exp, err = p.finishStdCall(exp)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return exp, nil
}

func (p *parser) finishStdCall(callee expression) (expression, error) {
	args := make(map[string]expression, 0)
	if !p.check(tokenRightParenthesis) {
		for {

			argKey := p.advance()
			if argKey.Type == tokenIdentifier {
				if _, err := p.consume(tokenColon, "Expected ':' after the argument"); err != nil {
					return nil, err
				}

				exp, err := p.expression()
				if err != nil {
					return nil, err
				}
				args[argKey.Lexeme] = exp
			}

			if !p.match(tokenComma) {
				break
			}
		}
	}

	_, err := p.consume(tokenRightParenthesis, "Expected ')' after arguments")
	if err != nil {
		return nil, err
	}
	return stdCallExpression{
		args:   args,
		callee: callee,
	}, nil
}

func (p *parser) primary() (expression, error) {
	if p.match(tokenFalse) {
		return literalExpression{value: false}, nil
	}
	if p.match(tokenTrue) {
		return literalExpression{value: true}, nil
	}
	if p.match(tokenNumber, tokenString) {
		return literalExpression{value: p.previous().Literal}, nil
	}
	if p.match(tokenIdentifier) {
		op := p.previous()

		if p.match(tokenEqual) {
			exp, err := p.expression()
			if err != nil {
				return nil, err
			}
			return assignExpression{
				op:  op,
				exp: exp,
			}, nil
		}
		return variableExpression{
			op: op,
		}, nil
	}
	if p.match(tokenLeftParenthesis) {
		exp, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(tokenRightParenthesis, "Expect ')' after expression"); err != nil {
			return nil, err
		}
		return groupingExpression{exp: exp}, nil
	}

	err := p.error(p.peek(), "Expected expression")

	return nil, err
}

func (p *parser) match(ts ...tokenType) bool {
	for _, t := range ts {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *parser) consume(t tokenType, message string) (token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	err := p.error(p.peek(), message)
	return token{}, err
}

func (p *parser) error(tok token, message string) error {
	if tok.Type == tokenEndOfFile {
		return fmt.Errorf("Parsing error at end of line %d - %s", tok.Line, message)
	}
	return fmt.Errorf("Parsing error at %s of line %d - %s", tok.Lexeme, tok.Line, message)
}

func (p *parser) check(t tokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *parser) advance() token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *parser) isAtEnd() bool {
	return p.peek().Type == tokenEndOfFile
}

func (p *parser) peek() token {
	return p.tokens[p.current]
}

func (p *parser) previous() token {
	return p.tokens[p.current-1]
}
