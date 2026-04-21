package logger

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/zeroSal/go-semantic-log/ansi"
)

var _ LoggerInterface = (*ConsoleLogger)(nil)

type ConsoleLogger struct {
	out   io.Writer
	level Level
	mutex sync.Mutex
}

func NewConsoleLogger(level ...Level) *ConsoleLogger {
	l := &ConsoleLogger{out: os.Stderr}
	if len(level) > 0 {
		l.level = level[0]
	}
	return l
}

func (l *ConsoleLogger) GetIdentifier() string {
	return "console"
}

func (l *ConsoleLogger) SetLevel(level Level) {
	l.level = level
}

func (l *ConsoleLogger) GetLevel() Level {
	return l.level
}

func (l *ConsoleLogger) Init() error {
	return nil
}

func (l *ConsoleLogger) Close() error {
	return nil
}

func (l *ConsoleLogger) log(color, prefix, msg string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	_, err := fmt.Fprintf(
		l.out,
		"%s%s%s%s\n",
		color,
		prefix,
		msg,
		ansi.Reset,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to console logger: %v\n", err)
	}
}

func (l *ConsoleLogger) Debug(msg string) {
	if !LevelDebug.ShouldLog(l.level) {
		return
	}
	l.log(ansi.White, "[•] ", msg)
}

func (l *ConsoleLogger) Info(msg string) {
	if !LevelInfo.ShouldLog(l.level) {
		return
	}
	l.log(ansi.Blue, "[i] ", msg)
}

func (l *ConsoleLogger) Warn(msg string) {
	if !LevelWarn.ShouldLog(l.level) {
		return
	}
	l.log(ansi.Yellow, "[!] ", msg)
}

func (l *ConsoleLogger) Error(msg string) {
	if !LevelError.ShouldLog(l.level) {
		return
	}
	l.log(ansi.Red, "[×] ", msg)
}

func (l *ConsoleLogger) Log(msg string) {
	l.log(ansi.Reset, "", msg)
}

func (l *ConsoleLogger) List(msgs []string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	for _, msg := range msgs {
		_, err := fmt.Fprintf(l.out, " · %s\n", msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to console logger: %v\n", err)
		}
	}
}
