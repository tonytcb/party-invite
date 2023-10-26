//nolint:golint,errcheck // In the case that the Write method fails, there's not too much to do in the logger package.
package logger

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/tonytcb/party-invite/pkg/infrastructure/config"
)

type Level string

const (
	LevelInfo  Level = "INFO"
	LevelError Level = "ERROR"
	LevelFatal Level = "FATAL"
)

type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
	WithCorrelationID(correlationID string) Logger
	FromContext(ctx context.Context) Logger
}

type emptyWriter struct {
}

func (e emptyWriter) Write(_ []byte) (int, error) {
	return 0, nil
}

type SimpleLogger struct {
	sync.Mutex

	writer io.Writer
	cid    string
}

func NewLogger(w io.Writer) Logger {
	return &SimpleLogger{writer: w}
}

func NewEmptyLogger() Logger {
	return NewLogger(&emptyWriter{})
}

func (s *SimpleLogger) Infof(format string, v ...any) {
	s.write(formatContent(LevelInfo, s.cid, format, v...))
}

func (s *SimpleLogger) Errorf(format string, v ...any) {
	s.write(formatContent(LevelError, s.cid, format, v...))
}

func (s *SimpleLogger) Fatalf(format string, v ...any) {
	content := formatContent(LevelFatal, s.cid, format, v...)

	s.write(content)

	panic(errors.New(string(content)))
}

func (s *SimpleLogger) WithCorrelationID(cid string) Logger {
	return &SimpleLogger{
		writer: s.writer,
		cid:    cid,
	}
}

func (s *SimpleLogger) FromContext(ctx context.Context) Logger {
	if cid, ok := ctx.Value(config.CorrelationIDKeyName).(string); ok {
		return s.WithCorrelationID(cid)
	}

	return s
}

func (s *SimpleLogger) write(content []byte) {
	s.Lock()
	defer s.Unlock()

	s.writer.Write(content)
}

func formatContent(level Level, cid string, format string, v ...any) []byte {
	now := time.Now().Format(time.RFC3339)

	args := []any{now, level, cid}
	args = append(args, v...)

	return []byte(fmt.Sprintf("[%s] [%s] [cid=%s] "+format+"\n", args...))
}
