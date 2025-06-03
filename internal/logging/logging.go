package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

type LogLevel int

const (
	LevelQuiet LogLevel = iota
	LevelInfo
	LevelDebug
	LevelTrace
)

var (
	currentLevel LogLevel = LevelQuiet
	logger       *log.Logger
)

func init() {
	logger = log.New(os.Stdout, "", 0)
}

func SetLevel(level int) {
	if level < 0 {
		level = 0
	}
	if level > 3 {
		level = 3
	}
	currentLevel = LogLevel(level)
}

func formatMessage(level string, msg string) string {
	return fmt.Sprintf("[%s] [%s] %s", time.Now().Format("2006-01-02 15:04:05.000"), level, msg)
}

func Info(msg string, args ...interface{}) {
	if currentLevel >= LevelInfo {
		logger.Println(formatMessage("INFO", fmt.Sprintf(msg, args...)))
	}
}

func Debug(msg string, args ...interface{}) {
	if currentLevel >= LevelDebug {
		logger.Println(formatMessage("DEBUG", fmt.Sprintf(msg, args...)))
	}
}

func Trace(msg string, args ...interface{}) {
	if currentLevel >= LevelTrace {
		logger.Println(formatMessage("TRACE", fmt.Sprintf(msg, args...)))
	}
}

func Error(msg string, args ...interface{}) {
	logger.Println(formatMessage("ERROR", fmt.Sprintf(msg, args...)))
}

func Warn(msg string, args ...interface{}) {
	logger.Println(formatMessage("WARN", fmt.Sprintf(msg, args...)))
}
