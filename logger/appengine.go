package logger

import (
	"context"

	"google.golang.org/appengine/log"
)

// Appengine ...
type Appengine struct {
	ctx context.Context
}

// NewAppengine ...
func NewAppengine(ctx context.Context) *Appengine {
	return &Appengine{ctx}
}

// Debugf ...
func (l *Appengine) Debugf(format string, args ...interface{}) {
	log.Debugf(l.ctx, format, args...)
}

// Infof ...
func (l *Appengine) Infof(format string, args ...interface{}) {
	log.Infof(l.ctx, format, args...)
}

// Warningf ...
func (l *Appengine) Warningf(format string, args ...interface{}) {
	log.Warningf(l.ctx, format, args...)
}

// Errorf ...
func (l *Appengine) Errorf(format string, args ...interface{}) {
	log.Errorf(l.ctx, format, args...)
}

// Criticalf ...
func (l *Appengine) Criticalf(format string, args ...interface{}) {
	log.Criticalf(l.ctx, format, args...)
}
