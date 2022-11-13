package core

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type MyFormatter struct{}

var levelList = []string{
	"\033[1;51;91m[PANIC",
	"\033[0;51;91m[FATAL",
	"\033[91m[ERROR",
	"\033[93m[WARN",
	"\033[0m[INFO",
	"\033[95m[DEBUG",
	"\033[1;30m[TRACE",
}

func (mf *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	level := levelList[int(entry.Level)]
	strList := strings.Split(entry.Caller.File, "/")
	fileName := strList[len(strList)-1] + "-" + entry.Caller.Function
	b.WriteString(fmt.Sprintf("%s - %s] %s > %s \033[0m\n",
		entry.Time.Format("2006-01-02 15:04:05"), level, fileName, entry.Message))
	return b.Bytes(), nil
}

var Log = &logrus.Logger{
	Out:          os.Stderr,
	Level:        logrus.DebugLevel,
	ReportCaller: true,
	Formatter:    &MyFormatter{},
}
var log = Log
