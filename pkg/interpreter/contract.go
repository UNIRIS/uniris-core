package interpreter

import (
	"errors"
	"fmt"

	"github.com/uniris/uniris-core/pkg/chain"
)

type Trigger struct {
	kind TriggerType
	val  interface{}
}

func (t Trigger) Display() string {
	switch t.kind {
	case timeTrigger:
		return fmt.Sprintf("Type: %s, Value: %d", t.kind.String(), int64(t.val.(float64)))
	default:
		return fmt.Sprintf("Type: %s, Value: %v", t.kind.String(), t.val)
	}
}

type TriggerType int

func (tt TriggerType) String() string {
	switch tt {
	case timeTrigger:
		return "time"
	default:
		return "not supported"
	}
}

const (
	timeTrigger TriggerType = 0
)

type Conditions struct {
	OriginFamily expression
	PostPaidFee  expression
	Response     expression
	Inherit      expression
}

type Contract struct {
	tx         chain.Transaction
	Triggers   []Trigger
	Conditions Conditions
	actions    []statement
}

func (c Contract) analyze() error {
	//TODO: ensure inputs - outputs == 0
	//TODO: formal and static analysis to ensure the correctness of the code
	return nil
}

func (c Contract) execute(sc *Scope) (interface{}, error) {
	out := make([]interface{}, 0)
	if c.Conditions.Response != nil {
		ok, err := c.Conditions.Response.evaluate(sc)
		if err != nil {
			return nil, err
		}
		switch ok.(type) {
		case bool:
			if ok.(bool) {
				for _, a := range c.actions {
					res, err := a.evaluate(sc)
					if err != nil {
						return nil, err
					}
					out = append(out, res)
				}
				return out, nil
			}
			return nil, errors.New("contract not executed. answer conditions are not respected")
		default:
			return nil, errors.New("answer conditions must be as boolean")
		}
	}

	for _, a := range c.actions {
		res, err := a.evaluate(sc)
		if err != nil {
			return nil, err
		}
		out = append(out, res)
	}
	return out, nil
}
