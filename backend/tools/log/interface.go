package log

import (
	"github.com/sirupsen/logrus"
)

type Ilog interface {
	New(args ...interface{}) Ilog
	Fatal(msg interface{}, args ...interface{})
	Error(msg interface{}, args ...interface{})
	Warn(msg interface{}, args ...interface{})
	Info(msg interface{}, args ...interface{})
	Debug(msg interface{}, args ...interface{})
	Trace(msg interface{}, args ...interface{})
}

type Level uint8

const (
	FatalLevel Level = iota
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

var lg Ilog

func init() {
	defaultLevel := logrus.Level(InfoLevel + 1)
	lg = &log_logrus{level: defaultLevel}
	logrus.SetLevel(defaultLevel)
}
func SetLevel(l Level) {
	lv := logrus.Level(l + 1)
	lg.(*log_logrus).level = lv
	logrus.SetLevel(lv)
}

func New(args ...interface{}) Ilog {
	return lg.New(args...)
}
func Fatal(msg interface{}, args ...interface{}) {
	lg.Fatal(msg, args...)
}
func Error(msg interface{}, args ...interface{}) {
	lg.Error(msg, args...)
}
func Warn(msg interface{}, args ...interface{}) {
	lg.Warn(msg, args...)
}
func Info(msg interface{}, args ...interface{}) {
	lg.Info(msg, args...)
}
func Debug(msg interface{}, args ...interface{}) {
	lg.Debug(msg, args...)
}
func Trace(msg interface{}, args ...interface{}) {
	lg.Trace(msg, args...)
}
