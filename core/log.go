package core

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/ylamothe/logrustash"
	"os"
	"strings"
)

/*
log.go
Initialize logger and format outputs
*/

type formatter struct{}

var levelListColor = []string{
	"\033[1;51;91m[PANIC",
	"\033[0;51;91m[FATAL",
	"\033[91m[ERROR",
	"\033[93m[WARNING",
	"\033[0m[INFO",
	"\033[95m[DEBUG",
	"\033[1;30m[TRACE",
}
var levelListPlain = []string{
	"[PANIC",
	"[FATAL",
	"[ERROR",
	"[WARNING",
	"[INFO",
	"[DEBUG",
	"[TRACE",
}
var logLevel logrus.Level
var Log *logrus.Logger

func GetLogger() *logrus.Logger {
	return Log
}

func DefineLogger() {
	log := logrus.New()
	if os.Getenv("D2LIB_sockl") == "true" {
		go addSock(log)
	}
	log.Out = os.Stderr
	log.Level = func() logrus.Level {
		switch os.Getenv("D2LIB_loglv") {
		case "trace":
			logLevel = logrus.TraceLevel
		case "debug":
			logLevel = logrus.DebugLevel
		case "info":
			logLevel = logrus.InfoLevel
		case "warning":
			logLevel = logrus.WarnLevel
		case "error":
			logLevel = logrus.ErrorLevel
		case "panic":
			logLevel = logrus.PanicLevel
		case "fatal":
			logLevel = logrus.FatalLevel
		default:
			fmt.Printf("Unknown log level: %s.\n", os.Getenv("D2LIB_loglv"))
			logLevel = logrus.InfoLevel
		}
		return logLevel
	}()
	log.ReportCaller = true
	log.Formatter = &formatter{}
	Log = log
}

func addSock(log *logrus.Logger) {
	hook, err := logrustash.NewAsyncHook(os.Getenv("D2LIB_sprot"), os.Getenv("D2LIB_saddr"), os.Getenv("D2LIB_sapp"))
	if err != nil {
		log.Fatal(err)
	}
	hook.TimeFormat = "2006-01-02 15:04:05"
	log.Hooks.Add(hook)
}

func (mf *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	var level string
	if os.Getenv("D2LIB_logcl") == "true" {
		level = levelListColor[int(entry.Level)]
	} else {
		level = levelListPlain[int(entry.Level)]
	}
	strList := strings.Split(entry.Caller.File, "/")
	fileName := strList[len(strList)-1] + "-" + entry.Caller.Function
	b.WriteString(fmt.Sprintf("%s - %s] %s > %s \033[0m\n",
		entry.Time.Format("2006-01-02 15:04:05"), level, fileName, entry.Message))
	return b.Bytes(), nil
}
