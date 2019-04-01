package logging

import (
	"fmt"
	"log"
	"net"
	"time"
)

//Loglevel describes the log level
type Loglevel int

const (
	errorLogLevel Loglevel = iota
	infoLogLevel
	debugLogLevel
)

//Logger describe a logger structure
type Logger struct {
	ioType   string
	log      *log.Logger
	appID    string
	hostID   net.IP
	logLevel Loglevel
}

//NewLogger create a new logger with the good parameteres
func NewLogger(o string, l *log.Logger, appid string, hostid net.IP, level Loglevel) Logger {

	return Logger{
		ioType:   o,
		log:      l,
		appID:    appid,
		logLevel: level,
		hostID:   hostid,
	}

}

//Level Return the log level
func (l *Logger) Level() Loglevel {
	return l.logLevel
}

//Writer Return the log IO Writer
func (l *Logger) Writer() string {
	return l.ioType
}

//Error write an Error message
func (l *Logger) Error(data string) {
	line := formatLogLine(l.appID, "ERROR", l.hostID, data)
	l.log.Println(line)

}

//Info write an Info message
func (l *Logger) Info(data string) {
	if l.logLevel >= 1 {
		line := formatLogLine(l.appID, "INFO", l.hostID, data)
		l.log.Println(line)
	}
}

//Debug write a debug message
func (l *Logger) Debug(data string) {
	if l.logLevel >= 2 {
		line := formatLogLine(l.appID, "DEBUG", l.hostID, data)
		l.log.Println(line)
	}
}

func formatLogLine(app string, ll string, hostID net.IP, data string) string {
	t := formatTime()
	return fmt.Sprintf("%s - [%s] [%s] [%s] \"%s\"", hostID.String(), t, ll, app, data)
}

func formatTime() string {
	t := time.Now()
	return t.Format("02/Jan/2006:15:04:05 -0700")
}
