package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/pkg/logging"
	"net"
	"testing"
)

/*
Scenario: Create file logger
	Given some inputs parameters
	When we create logger
	Then logger level has the good value
*/
func TestLoggerLevel(t *testing.T) {
	logger := createLogger(miningAppID, "file", "/tmp", "debug", net.ParseIP("127.0.0.1"))
	assert.Equal(t, logging.Loglevel(2), logger.Level())
}

/*
Scenario: Create stdout logger
	Given some inputs parameters
	When we create logger
	Then logger type is stdout
*/
func TestStdoutLogger(t *testing.T) {
	logger := createLogger(miningAppID, "stdout", "/tmp", "debug", net.ParseIP("127.0.0.1"))
	assert.Equal(t, "stdout", logger.Writer())
}

/*
Scenario: Create file logger
	Given some inputs parameters
	When we create logger
	Then logger type is file
*/
func TestFileLogger(t *testing.T) {
	logger := createLogger(miningAppID, "file", "/tmp", "debug", net.ParseIP("127.0.0.1"))
	assert.Equal(t, "file", logger.Writer())
}

/*
Scenario: Create file logger
	Given a dir not existing
	When we create logger
	Then logger type is stdout
*/
func TestFileLoggerError(t *testing.T) {
	logger := createLogger(miningAppID, "file", "/007f6d2c61d5769872a6d10e5bf809a6", "debug", net.ParseIP("127.0.0.1"))
	assert.Equal(t, "stdout", logger.Writer())
}
