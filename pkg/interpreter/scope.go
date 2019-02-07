package interpreter

import (
	"fmt"
)

//Scope defines the interpreter scope including variables and parent(linked) scope
type Scope struct {
	parent    *Scope
	variables map[string]interface{}
}

//NewScope creates a new interpreter context
func NewScope(parent *Scope) *Scope {
	return &Scope{
		variables: make(map[string]interface{}, 0),
		parent:    parent,
	}
}

//SetValue add variable to the context or its parent associated by the provided name
func (env *Scope) SetValue(name string, value interface{}) {
	if env.parent != nil {
		_, err := env.parent.GetValue(name)
		if err != nil {
			if err.Error() == fmt.Sprintf("Undefined variables %s", name) {
				env.variables[name] = value
				return
			}
			panic(err)
		}
		env.parent.SetValue(name, value)
	} else {
		env.variables[name] = value
	}
}

//GetValue returns the variable value from its name. A recursive retrival is done if the value is present in the parent contexts
func (env *Scope) GetValue(name string) (interface{}, error) {
	v, exist := env.variables[name]
	if exist {
		return v, nil
	}

	if env.parent != nil {
		return env.parent.GetValue(name)
	}

	return nil, fmt.Errorf("Undefined variables %s", name)
}
