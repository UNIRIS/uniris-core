package logging

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	system "github.com/uniris/uniris-core/pkg/system"
)

//LogType defines the log type
type LogType int

const (
	//Stdout define an stdout logger
	Stdout LogType = iota

	//File define a file logger
	File
)

//LogLevel defines the log level
type LogLevel int

const (
	//Error loglevel
	Error LogLevel = iota

	//Info loglevel
	Info

	//Debug loglevel
	Debug
)

//AppID defines the application id behind the logging
type AppID int

const (
	//Discoverie logs
	Discoverie AppID = iota

	//Mining logs
	Mining
)

const (
	//DiscoverieLogFile define the name for discoveries logs file
	DiscoverieLogFile = "Discoverie.log"

	//MiningLogFile define the name for mining logs file
	MiningLogFile = "Mining.log"
)

//Logger describe a logger structure
type Logger struct {
	logType  LogType
	logDir   string
	logLevel LogLevel
	hostID   net.IP
}

//NewLogger create a new logger with the good parameteres
func NewLogger(ty string, dir string, level string, privateNetwork bool, privateIface string) Logger {

	//get HostID using the system reader
	sysReader := system.NewReader(privateNetwork, privateIface)
	ip, err := sysReader.IP()
	if err != nil {
		log.Fatal("[Fatal] Cannot get the hostID")
	}

	//check if log type has the good value
	if ty != "stdout" && ty != "file" {
		log.Fatal("[Fatal] log-type value should be (stdout|file)")
	}

	var t LogType

	if ty == "file" {
		t = File
	} else {
		t = Stdout
	}

	//check if log-level has the good value
	if level != "info" && level != "error" && level != "debug" {
		log.Fatal("[Fatal] log-level value should be (info|error|debug)")
	}

	var ll LogLevel

	if level == "info" {
		ll = Info
	} else if level == "error" {
		ll = Error
	} else {
		ll = Debug
	}

	if t == Stdout {
		return Logger{
			logType:  t,
			logDir:   "",
			logLevel: ll,
			hostID:   ip,
		}
	}

	if t == File {

		src, err := os.Stat(dir)

		//check if logdir exist or not
		if os.IsNotExist(err) {
			log.Println("[Error] log-dir" + dir + "does not exist, please create the adequate directory")
			return Logger{
				logType:  Stdout,
				logDir:   "a",
				logLevel: ll,
				hostID:   ip,
			}
		}

		//check if logdir is not a file
		if src.Mode().IsRegular() {
			log.Println("[Erro] log-dir" + dir + "is a file, please create the adequate directory")
			return Logger{
				logType:  Stdout,
				logDir:   "b",
				logLevel: ll,
				hostID:   ip,
			}
		}

		//create Discoverie log files
		err = createFile(dir, DiscoverieLogFile)
		if err != nil {
			log.Print(err.Error())
			return Logger{
				logType:  Stdout,
				logDir:   "c",
				logLevel: ll,
				hostID:   ip,
			}
		}

		//create Mining log files
		err = createFile(dir, MiningLogFile)
		if err != nil {
			log.Print(err.Error())
			return Logger{
				logType:  Stdout,
				logDir:   "d",
				logLevel: ll,
				hostID:   ip,
			}
		}

	}

	return Logger{
		logType:  t,
		logDir:   dir,
		logLevel: ll,
		hostID:   ip,
	}

}

//AppendLog write a log line taking into consideration the loglevel and the Appid and the logger type
func (l *Logger) AppendLog(path string, app AppID, ll LogLevel, hostID net.IP, data string) error {

	if l.logType == Stdout && ll >= l.logLevel {
		writeOnStdout(app, ll, hostID, data)
	}

	if l.logType == File && ll >= l.logLevel {
		writeOnFile(path, app, ll, hostID, data)
	}

	return nil
}

func writeOnStdout(app AppID, ll LogLevel, hostID net.IP, data string) {
	t := formatTime()
	l := formatLogLevel(ll)
	a := formatAppID(app)
	log.Println(fmt.Sprintf("%s - [%s] [%s] [%s] \"%s\"", hostID.String(), t, l, a, data))
}

func writeOnFile(path string, app AppID, ll LogLevel, hostID net.IP, data string) {
	if app == Discoverie {
		err := writeLineOnFile(path, DiscoverieLogFile, ll, hostID, data)
		if err != nil {
			writeOnStdout(app, ll, hostID, data)
		}

	} else {
		err := writeLineOnFile(path, DiscoverieLogFile, ll, hostID, data)
		if err != nil {
			writeOnStdout(app, ll, hostID, data)
		}
	}
}

func formatTime() string {
	t := time.Now()
	return t.Format("02/Jan/2006:15:04:05 -0700")
}

func formatAppID(app AppID) string {
	if app == Discoverie {
		return "DISCOVERIE"
	}
	return "MINING"
}

func formatLogLevel(ll LogLevel) string {
	if ll == Info {
		return "INFO"
	} else if ll == Error {
		return "ERROR"
	}
	return "DEBUG"
}

func createFile(path string, filename string) error {

	_, err := os.Stat(path + "/" + filename)

	if os.IsNotExist(err) {
		file, err := os.Create(path + "/" + filename)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

func writeLineOnFile(path string, filename string, ll LogLevel, hostID net.IP, data string) error {

	f, err := os.OpenFile(path+"/"+filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	t := formatTime()
	l := formatLogLevel(ll)
	line := fmt.Sprintf("%s - [%s] [%s] \"%s\"", hostID.String(), t, l, data)

	logger := log.New(f, "", 0)
	logger.Println(line)
	return nil

}
