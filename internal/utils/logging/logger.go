package logging

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Logger struct {
	file *os.File
	mu   sync.Mutex
}

func NewLogger(path string) (*Logger, error) {
	file, err := os.OpenFile(
		path,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, err
	}

	return &Logger{file: file}, nil
}

func (l *Logger) Log(msg string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	line := fmt.Sprintf(
		"[%s] %s\n",
		time.Now().Format(time.RFC3339),
		msg,
	)

	_, err := l.file.WriteString(line)

	return err
}

func (l *Logger) Close() error {
	return l.file.Close()
}

func NewFile(path string) (*os.File, error) {
	return os.OpenFile(
		path,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
}
