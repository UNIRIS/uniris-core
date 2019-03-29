package logging

import (
	"fmt"
	"log"
	"net"
	"time"
)

//Logger describe a logger structure
type Logger struct {
	log      *log.Logger
	appID    string
	hostID   net.IP
	logLevel int
}

//NewLogger create a new logger with the good parameteres
func NewLogger(l *log.Logger, appid string, hostid net.IP, ll int) Logger {

	return Logger{
		log:      l,
		appID:    appid,
		logLevel: ll,
		hostID:   hostid,
	}

}

//Error write an Error message
func (l *Logger) Error(app string, hostID net.IP, data string) {
	line := formatLogLine(app, "ERROR", hostID, data)
	l.log.Println(line)

}

//Info write an Info message
func (l *Logger) Info(app string, hostID net.IP, data string) {
	if l.logLevel >= 1 {
		line := formatLogLine(app, "ERROR", hostID, data)
		l.log.Println(line)
	}
}

//Debug write a debug message
func (l *Logger) Debug(app string, hostID net.IP, data string) {
	if l.logLevel >= 2 {
		line := formatLogLine(app, "ERROR", hostID, data)
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
