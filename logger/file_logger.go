package logger

import (
	"fmt"
	"os"
	"sync"
)

var _ LoggerInterface = (*FileLogger)(nil)

type FileLogger struct {
	path  string
	label string
	level Level

	file *os.File
	mu   sync.Mutex
}

func NewFileLogger(
	path,
	label string,
	level ...Level,
) *FileLogger {
	l := &FileLogger{
		path:  path,
		label: label,
	}
	if len(level) > 0 {
		l.level = level[0]
	}
	return l
}

func (l *FileLogger) SetLevel(level Level) {
	l.level = level
}

func (l *FileLogger) GetLevel() Level {
	return l.level
}

func (l *FileLogger) Init() error {
	f, err := os.OpenFile(
		l.path,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)

	if err != nil {
		return err
	}

	l.file = f

	return nil
}

func (l *FileLogger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}

	return nil
}

func (l *FileLogger) log(prefix, msg string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file == nil {
		return fmt.Errorf("logger %s not initialized. Call Init() first", l.label)
	}

	_, err := fmt.Fprintf(l.file, "%s %s\n", prefix, msg)
	if err != nil {
		return fmt.Errorf("logger %s: write failed", l.label)
	}

	return nil
}

func (l *FileLogger) GetIdentifier() string {
	return l.label
}

func (l *FileLogger) Debug(msg string) error {
	if !LevelDebug.ShouldLog(l.level) {
		return nil
	}
	return l.log("[DEBUG]", msg)
}

func (l *FileLogger) Info(msg string) error {
	if !LevelInfo.ShouldLog(l.level) {
		return nil
	}
	return l.log("[INFO]", msg)
}

func (l *FileLogger) Warn(msg string) error {
	if !LevelWarn.ShouldLog(l.level) {
		return nil
	}
	return l.log("[WARNING]", msg)
}

func (l *FileLogger) Error(msg string) error {
	if !LevelError.ShouldLog(l.level) {
		return nil
	}
	return l.log("[ERROR]", msg)
}
