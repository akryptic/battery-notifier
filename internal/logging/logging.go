package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
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
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// Base format with just timestamp and level
	baseFormat := fmt.Sprintf("[%s] [%s]", timestamp, level)

	// Add location info based on current log level
	if currentLevel >= LevelDebug {
		pc, file, line, ok := runtime.Caller(2)
		if ok {
			if currentLevel >= LevelTrace {
				// Include file, line, and function name for TRACE level
				funcName := runtime.FuncForPC(pc).Name()
				// Extract just the function name (remove package path)
				if lastSlash := filepath.Base(funcName); lastSlash != "" {
					funcName = lastSlash
				}
				location := fmt.Sprintf("%s:%d:%s", filepath.Base(file), line, funcName)
				return fmt.Sprintf("%s [%s] %s", baseFormat, location, msg)
			} else {
				// Include only file and line for DEBUG level
				location := fmt.Sprintf("%s:%d", filepath.Base(file), line)
				return fmt.Sprintf("%s [%s] %s", baseFormat, location, msg)
			}
		}
	}

	// For INFO level and below, just timestamp, level, and message
	return fmt.Sprintf("%s %s", baseFormat, msg)
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
