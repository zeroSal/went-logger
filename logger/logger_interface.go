package logger

import (
	"io"
)

type LoggerInterface interface {
	io.Closer
	Init() error
	SetLevel(level Level)
	GetLevel() Level
	GetIdentifier() string
	Debug(msg string) error
	Info(msg string) error
	Warn(msg string) error
	Error(msg string) error
}
