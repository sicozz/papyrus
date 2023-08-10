package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/sicozz/papyrus/utils/constants"
)

const (
	flags       = log.LstdFlags
	infoPrefix  = `INF: `
	warnPrefix  = `WRN: `
	errorPrefix = `ERR: `
)

type AggregatedLogger interface {
	// TODO: rename functions to shorthand
	Info(...any)
	Warn(...any)
	Error(...any)
}

type AgLog struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	layer       constants.Layer
	domain      constants.Domain
}

func (f AgLog) log(la *log.Logger, v []any) {
	l := make([]any, 1, 1)
	l[0] = fmt.Sprint("[", f.layer, ":", f.domain, "]")
	l = append(l, v...)
	la.Println(l...)
	return
}

func (f AgLog) Info(v ...any) {
	f.log(f.infoLogger, v)
}

func (f AgLog) Warn(v ...any) {
	f.log(f.warnLogger, v)
}

func (f AgLog) Error(v ...any) {
	f.log(f.errorLogger, v)
}

func NewAggregatedLogger(layer constants.Layer, domain constants.Domain) AggregatedLogger {
	infoLogger := log.New(os.Stdout, infoPrefix, flags)
	warnLogger := log.New(os.Stdout, warnPrefix, flags)
	errorLogger := log.New(os.Stderr, errorPrefix, flags)

	return AgLog{
		infoLogger,
		warnLogger,
		errorLogger,
		layer,
		domain,
	}
}
