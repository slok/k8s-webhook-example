package log

import (
	"github.com/sirupsen/logrus"
)

// KV is a helper type for structured logging fields usage.
type KV map[string]interface{}

// Logger is the interface that the loggers used by the library will use.
type Logger interface {
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	WithKV(KV) Logger
}

// Dummy logger doesn't log anything.
const Dummy = dummy(0)

var _ Logger = Dummy

type dummy int

func (d dummy) Infof(format string, args ...interface{})    {}
func (d dummy) Warningf(format string, args ...interface{}) {}
func (d dummy) Errorf(format string, args ...interface{})   {}
func (d dummy) Debugf(format string, args ...interface{})   {}
func (d dummy) WithKV(KV) Logger                            { return d }

type logger struct {
	*logrus.Entry
}

// NewLogrus returns a new log.Logger for a logrus implementation.
func NewLogrus(l *logrus.Entry) Logger {
	return logger{Entry: l}
}

func (l logger) WithKV(kv KV) Logger {
	newLogger := l.Entry.WithFields(logrus.Fields(kv))
	return NewLogrus(newLogger)
}
