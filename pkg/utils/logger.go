package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	// colorCyan   = "\033[36m"
)

type Logger struct {
	logger *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", 0),
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.print(colorGreen, "INFO", format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.print(colorYellow, "WARN", format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.print(colorRed, "ERROR", format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.print(colorBlue, "DEBUG", format, v...)
}

func (l *Logger) print(color, level, format string, v ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, v...)
	output := fmt.Sprintf("%s[%s] [%s] %s%s", color, timestamp, level, msg, colorReset)
	l.logger.Println(output)
}
