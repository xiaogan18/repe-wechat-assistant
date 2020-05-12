package log

import (
	"github.com/sirupsen/logrus"
)

type log_logrus struct {
	level logrus.Level
	filed []interface{}
}

func (t *log_logrus) New(args ...interface{}) Ilog {
	checkArgs(len(args))
	return &log_logrus{
		level: t.level,
		filed: args,
	}
}
func (t *log_logrus) Fatal(msg interface{}, args ...interface{}) {
	checkArgs(len(args))
	if len(args) > 0 {
		t.withFiled(args).Fatal(msg)
	} else {
		logrus.Fatal(msg)
	}
}
func (t *log_logrus) Error(msg interface{}, args ...interface{}) {
	if t.level < logrus.ErrorLevel {
		return
	}
	checkArgs(len(args))
	if len(args) > 0 {
		t.withFiled(args).Error(msg)
	} else {
		logrus.Error(msg)
	}
}
func (t *log_logrus) Warn(msg interface{}, args ...interface{}) {
	if t.level < logrus.WarnLevel {
		return
	}
	checkArgs(len(args))
	if len(args) > 0 {
		t.withFiled(args).Info(msg)
	} else {
		logrus.Info(msg)
	}
}
func (t *log_logrus) Info(msg interface{}, args ...interface{}) {
	if t.level < logrus.InfoLevel {
		return
	}
	checkArgs(len(args))
	if len(args) > 0 {
		t.withFiled(args).Info(msg)
	} else {
		logrus.Info(msg)
	}
}
func (t *log_logrus) Debug(msg interface{}, args ...interface{}) {
	if t.level < logrus.DebugLevel {
		return
	}
	checkArgs(len(args))
	if len(args) > 0 {
		t.withFiled(args).Debug(msg)
	} else {
		logrus.Debug(msg)
	}
}
func (t *log_logrus) Trace(msg interface{}, args ...interface{}) {
	if t.level < logrus.TraceLevel {
		return
	}
	checkArgs(len(args))
	if len(args) > 0 {
		t.withFiled(args).Trace(msg)
	} else {
		logrus.Trace(msg)
	}
}
func (t *log_logrus) withFiled(args []interface{}) *logrus.Entry {
	if len(t.filed) > 0 {
		args = append(t.filed, args...)
	}
	var lg *logrus.Entry
	for i := 0; i < len(args); i += 2 {
		if lg == nil {
			lg = logrus.WithField(args[i].(string), args[i+1])
		} else {
			lg = lg.WithField(args[i].(string), args[i+1])
		}
	}
	return lg
}
func checkArgs(length int) {
	if length%2 != 0 {
		panic("log args should be even")
	}
}
