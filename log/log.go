package log

import (
	"github.com/mattn/go-colorable"
	"io"
	"os"
	"regexp"
	"time"
)

var cleaner = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// Logger represents a Writer which writes the log to the provided file as well as stdout.
type Logger struct {
	file   *os.File
	stdout io.Writer
}

// New creates a new logger to be used with any log package. It is designed to write to a log file as well as
// stdout to allow you to store logs from the proxy.
func New(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{
		file:   f,
		stdout: colorable.NewColorableStdout(),
	}, nil
}

// Write ...
func (l *Logger) Write(p []byte) (int, error) {
	if n, err := l.stdout.Write(p); err != nil {
		return n, err
	}

	cleaned := cleaner.ReplaceAllString(string(p), "")
	return l.file.WriteString(time.Now().Format("2006-1-2") + " " + cleaned)
}
