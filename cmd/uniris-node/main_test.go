package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/pkg/logging"
	"net"
	"os"
	"testing"
)

/*
Scenario: Create stdout logger
	Given some inputs parameters
	When we create logger
	Then logger type is stdout
*/
func TestStdoutLogger(t *testing.T) {
	logger := createLogger("mining", "stdout", "/tmp", "debug", net.ParseIP("127.0.0.1"))
	assert.Equal(t, os.Stdout, logger.GetWriter())
}

/*
Scenario: Create file logger
	Given some inputs parameters
	When we create logger
	Then logger type is file
*/
func TestFileLogger(t *testing.T) {
	logger := createLogger("mining", "file", "/tmp", "debug", net.ParseIP("127.0.0.1"))
	assert.Equal(t, logging.Loglevel(2), logger.GetLevel())
	f, _ := os.OpenFile("/tmp/test", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	assert.IsType(t, f, logger.GetWriter())
}

/*
Scenario: Create file logger
	Given a dir not existing
	When we create logger
	Then logger type is stdout
*/
func TestFileLoggerError(t *testing.T) {
	logger := createLogger("mining", "file", "/007f6d2c61d5769872a6d10e5bf809a6", "debug", net.ParseIP("127.0.0.1"))
	assert.Equal(t, os.Stdout, logger.GetWriter())
}
