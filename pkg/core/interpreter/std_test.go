package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimestampFunction(t *testing.T) {
	f := timestampFunc{}
	res, err := f.call(nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, time.Now().Unix(), res)
}
