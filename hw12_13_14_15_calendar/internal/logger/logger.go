package logger

import (
	"fmt"
	"log"
	"os"
)

const (
	errorLevel = "error"
	warnLevel  = "warn"
	debugLevel = "debug"
)

type Logger struct {
	LogFile  string
	LogLevel string
}

func New(logFile string, logLevel string) *Logger {
	return &Logger{
		logFile,
		logLevel,
	}
}

func (l Logger) Info(msg string) {
	fmt.Println(msg)
}

func (l Logger) Error(msg string) {
	if l.LogLevel == errorLevel {
		err := logToFile(l.LogFile, msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func (l Logger) Warn(msg string) {
	if l.LogLevel == errorLevel || l.LogLevel == warnLevel {
		err := logToFile(l.LogFile, msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func (l Logger) Debug(msg string) {
	if l.LogLevel == errorLevel || l.LogLevel == warnLevel || l.LogLevel == debugLevel {
		err := logToFile(l.LogFile, msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func logToFile(logFile string, msg string) error {
	f, err := os.OpenFile(logFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println(err)
		}
	}(f)
	if _, err := f.WriteString(msg + "\n"); err != nil {
		return err
	}
	return nil
}
