//nolint:golint,errcheck // In the case that the Write method fails, there's not too much to do in the logger package.
package logger

import (
	"context"
	"fmt"
	"git.codesubmit.io/sfox/golang-party-invite-ivsjhn/pkg/infrastructure/config"
	"io"
	"time"

	"github.com/pkg/errors"
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

type SimpleLogger struct {
	writer io.Writer
	cid    string
}

func NewLogger(w io.Writer) Logger {
	return &SimpleLogger{writer: w}
}

func (s SimpleLogger) Infof(format string, v ...any) {
	s.writer.Write(formatContent(LevelInfo, s.cid, format, v...))
}

func (s SimpleLogger) Errorf(format string, v ...any) {
	s.writer.Write(formatContent(LevelError, s.cid, format, v...))
}

func (s SimpleLogger) Fatalf(format string, v ...any) {
	content := formatContent(LevelFatal, s.cid, format, v...)

	s.writer.Write(content)

	panic(errors.New(string(content)))
}

func (s SimpleLogger) WithCorrelationID(cid string) Logger {
	return SimpleLogger{
		writer: s.writer,
		cid:    cid,
	}
}

func (s SimpleLogger) FromContext(ctx context.Context) Logger {
	if cid, ok := ctx.Value(config.CorrelationIDKeyName).(string); ok {
		return s.WithCorrelationID(cid)
	}

	return s
}

func formatContent(level Level, cid string, format string, v ...any) []byte {
	now := time.Now().Format(time.RFC3339)

	args := []any{now, level, cid}
	args = append(args, v...)

	return []byte(fmt.Sprintf("[%s] [%s] [cid=%s] "+format+"\n", args...))
}
