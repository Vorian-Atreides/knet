package knet

import (
	"log"
	"os"
)

// Logger define the interface needed by knet to log its internal errors.
type Logger interface {
	Warningf(format string, args ...interface{})
}

type logger struct {
	*log.Logger
}

func (l logger) Warningf(format string, args ...interface{}) {
	log.Printf("%s", args...)
}

var l Logger = logger{
	log.New(os.Stderr, "knet", log.LstdFlags),
}

// SetLogger to be used instead of the default logger.
func SetLogger(logger Logger) {
	l = logger
}
