package logger

import (
	"log"
	"os"
)

// Default ...
type Default struct {
	debug    *log.Logger
	info     *log.Logger
	warning  *log.Logger
	err      *log.Logger
	critical *log.Logger
}

// NewDefault ...
func NewDefault() *Default {
	return &Default{
		debug:    log.New(os.Stdout, "DEBUG: ", 0),
		info:     log.New(os.Stdout, "INFO: ", 0),
		warning:  log.New(os.Stdout, "WARNING: ", 0),
		err:      log.New(os.Stdout, "ERROR: ", 0),
		critical: log.New(os.Stdout, "CRITICAL: ", 0),
	}
}

// Debugf ...
func (l *Default) Debugf(format string, args ...interface{}) {
	l.debug.Printf(format, args...)
}

// Infof ...
func (l *Default) Infof(format string, args ...interface{}) {
	l.info.Printf(format, args...)
}

// Warningf ...
func (l *Default) Warningf(format string, args ...interface{}) {
	l.warning.Printf(format, args...)
}

// Errorf ...
func (l *Default) Errorf(format string, args ...interface{}) {
	l.err.Printf(format, args...)
}

// Criticalf ...
func (l *Default) Criticalf(format string, args ...interface{}) {
	l.critical.Printf(format, args...)
}
