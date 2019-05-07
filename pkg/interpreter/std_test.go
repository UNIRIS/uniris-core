package interpreter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNowFunc(t *testing.T) {
	res, err := nowFunc(nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, float64(time.Now().Unix()), res)
}
