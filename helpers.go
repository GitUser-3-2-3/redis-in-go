package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

type Logs struct {
	logError *log.Logger
	logInfo  *log.Logger
}

func NewLogger() *Logs {
	return &Logs{
		logInfo:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
		logError: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)}
}

func (rl *Logs) ServerError(err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = rl.logError.Output(2, trace)
}
