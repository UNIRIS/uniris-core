package logging

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
Scenario: Create stdout logger
	Given some inputs parameters
	When we create logger
	Then logger type is stdout
*/
func TestStdoutLogger(t *testing.T) {
	logger := NewLogger("stdout", "/var/log", "debug", false, "")
	assert.Equal(t, Stdout, logger.logType)
}

/*
Scenario: Create file logger
	Given some inputs parameters
	When we create logger
	Then logger type is file
*/
func TestFileLogger(t *testing.T) {
	logger := NewLogger("file", "/tmp/", "debug", false, "")
	assert.Equal(t, File, logger.logType)
}

/*
Scenario: Create file logger
	Given a dir not existing
	When we create logger
	Then logger type is stdout
*/
func TestFileLoggerError(t *testing.T) {
	logger := NewLogger("file", "/775f995a0d8fc7c5b35642308bbd37ae", "debug", false, "")
	assert.Equal(t, Stdout, logger.logType)
}
