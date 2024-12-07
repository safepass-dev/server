package logging

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	INFO = iota
	WARN
	ERROR
	FATAL
)

type Logger struct {
	logLevel int
	file     *os.File
}

func NewLogger(logLevel int, logFile string) (*Logger, error) {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{
		logLevel: logLevel,
		file:     file,
	}, nil
}

func getCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func getCallerInfo() (string, int) {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return "unknown", 0
	}

	fileParts := strings.Split(file, "/")
	return fileParts[len(fileParts)-1], line
}

const (
	RESET  = "\033[0m"
	RED    = "\033[31m"
	GREEN  = "\033[32m"
	BLUE   = "\033[34m"
	YELLOW = "\033[33m"
	WHITE  = "\033[97m"
	BGRED  = "\033[41m"
)

func (l *Logger) log(level int, msg string) {
	if level < l.logLevel {
		return
	}

	timestamp := getCurrentTime()

	file, line := getCallerInfo()

	var levelStr string
	var levelColor string

	switch level {
	case INFO:
		levelStr = "info"
		levelColor = BLUE
	case WARN:
		levelStr = "warn"
		levelColor = YELLOW
	case ERROR:
		levelStr = "error"
		levelColor = RED
	case FATAL:
		levelStr = "fatal"
		levelColor = WHITE
	}

	var logMsg string
	if level == INFO {
		logMsg = fmt.Sprintf("%s %s%s%s %s\n", timestamp, levelColor, levelStr, RESET, msg)
	} else {
		logMsg = fmt.Sprintf("%s %s%s%s %s:%d\n%s\n", timestamp, levelColor, levelStr, RESET, file, line, msg)
	}

	log.Println(logMsg)
}

func (l *Logger) Info(msg string) {
	l.log(INFO, msg)
}

func (l *Logger) Warn(msg string) {
	l.log(WARN, msg)
}

func (l *Logger) Error(msg string) {
	l.log(ERROR, msg)
}

func (l *Logger) Fatal(msg string) {
	l.log(FATAL, msg)
	os.Exit(1)
}

func (l *Logger) SetLogLevel(level int) {
	l.logLevel = level
}
