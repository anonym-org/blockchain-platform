package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Info(v ...any)
	Warn(v ...any)
	Error(v ...any)
	Fatal(v ...any)
}

type logger struct {
	log *logrus.Logger
}

func NewLogger() Logger {
	return &logger{
		log: &logrus.Logger{
			Out:   os.Stderr,
			Level: logrus.DebugLevel,
			Formatter: &logrus.TextFormatter{
				FullTimestamp: true,
			},
		},
	}
}

func (l *logger) Info(v ...any) {
	l.log.Info(v...)
}

func (l *logger) Warn(v ...any) {
	l.log.Warn(v...)
}

func (l *logger) Error(v ...any) {
	l.log.Error(v...)
}

func (l *logger) Fatal(v ...any) {
	l.log.Fatal(v...)
}
