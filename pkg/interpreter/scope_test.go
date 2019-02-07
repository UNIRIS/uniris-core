package interpreter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetValue(t *testing.T) {
	e := NewScope(nil)
	e.SetValue("a", 2)
	assert.Equal(t, 2, e.variables["a"])
}

func TestSetParentValue(t *testing.T) {

	enc := NewScope(nil)
	enc.SetValue("a", 2)

	e := NewScope(enc)

	assert.Nil(t, e.variables["a"])
	assert.Equal(t, 2, e.parent.variables["a"])

	e.SetValue("a", 5)
	assert.Equal(t, 5, e.parent.variables["a"])

	e.SetValue("b", 10)
	assert.Equal(t, 10, e.variables["b"])
	assert.Nil(t, e.parent.variables["b"])
}

func TestGetValue(t *testing.T) {
	e := NewScope(nil)
	e.SetValue("a", 2)
	val, err := e.GetValue("a")
	assert.Nil(t, err)
	assert.Equal(t, 2, val)
}

func TestGetParentValue(t *testing.T) {
	enc := NewScope(nil)
	enc.SetValue("a", 2)

	e := NewScope(enc)

	val, err := e.GetValue("a")
	assert.Nil(t, err)
	assert.Equal(t, 2, val)
}

func TestGetUndefinedValue(t *testing.T) {
	e := NewScope(nil)
	_, err := e.GetValue("a")
	assert.Error(t, err, "Undefined variable a")
}
