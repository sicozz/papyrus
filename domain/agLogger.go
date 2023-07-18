package domain

import (
	"log"
	"os"
)

const (
	flags       = log.LstdFlags | log.Lshortfile
	infoPrefix  = `INF: `
	warnPrefix  = `WRN: `
	errorPrefix = `ERR: `
)

type AggregatedLogger struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

func (f *AggregatedLogger) Info(v ...interface{}) {
	f.infoLogger.Println(v...)
}

func (f *AggregatedLogger) Warn(v ...interface{}) {
	f.warnLogger.Println(v...)
}

func (f *AggregatedLogger) Error(v ...interface{}) {
	f.errorLogger.Println(v...)
}

func NewAggregatedLogger() AggregatedLogger {
	infoLogger := log.New(os.Stdout, infoPrefix, flags)
	warnLogger := log.New(os.Stdout, warnPrefix, flags)
	errorLogger := log.New(os.Stderr, errorPrefix, flags)

	return AggregatedLogger{
		infoLogger,
		warnLogger,
		errorLogger,
	}
}

var AgLog AggregatedLogger = NewAggregatedLogger()
