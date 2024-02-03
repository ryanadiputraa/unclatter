package logger

import (
	"log"
	"os"
	"time"
)

type Logger interface {
	Info(v ...any)
	Warn(v ...any)
	Error(v ...any)
	Fatal(v ...any)
}

type logger struct {
	log *log.Logger
}

func NewLogger() Logger {
	location := time.UTC

	l := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	l.SetFlags(l.Flags() | log.LUTC)
	l.SetPrefix("[" + location.String() + "] ")
	return &logger{log: l}
}

func (l *logger) Info(v ...any) {
	l.log.Println(append([]any{"[INFO]:"}, v...)...)
}

func (l *logger) Warn(v ...any) {
	l.log.Println(append([]any{"[WARN]:"}, v...)...)
}

func (l *logger) Error(v ...any) {
	l.log.Println(append([]any{"[ERROR]:"}, v...)...)
}

func (l *logger) Fatal(v ...any) {
	l.log.Fatal(append([]any{"[FATAL]:"}, v...)...)
}
