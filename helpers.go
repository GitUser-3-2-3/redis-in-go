package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

type logs struct {
	logInfo  *log.Logger
	logError *log.Logger
}

func initLogs() *logs {
	return &logs{
		logError: log.New(os.Stdout, "ERROR::\t", log.Ldate|log.Ltime|log.Lshortfile),
		logInfo:  log.New(os.Stdout, "INFO::\t", log.Ldate|log.Ltime),
	}
}

func (l *logs) serverError(err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = l.logError.Output(2, trace)
}
