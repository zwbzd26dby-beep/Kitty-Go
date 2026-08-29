// Package observability provides logging, metrics, tracing, usage/cost
// tracking and audit for the whole agent (Master Architecture §23).
package observability

import (
	"fmt"
	"io"
	"log"
	"sync"
	"time"
)

// Level is a log severity.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	}
	return "UNKNOWN"
}

// Entry is a structured log record.
type Entry struct {
	Time  time.Time
	Level Level
	Msg   string
	Attrs map[string]any
}

// Logger is a minimal structured logger.
type Logger struct {
	mu  sync.Mutex
	lg  *log.Logger
	min Level
}

// NewLogger creates a Logger writing to out, filtering below min.
func NewLogger(out io.Writer, min Level) *Logger {
	return &Logger{lg: log.New(out, "", log.LstdFlags), min: min}
}

// Logf records a message at level with optional key/value attributes.
func (l *Logger) Logf(level Level, msg string, kv ...any) {
	if level < l.min || l.lg == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lg.Printf("%s %s %s", level, msg, flatten(kv))
}

// Debug/Info/Warn/Error helpers.
func (l *Logger) Debug(msg string, kv ...any) { l.Logf(LevelDebug, msg, kv...) }
func (l *Logger) Info(msg string, kv ...any)  { l.Logf(LevelInfo, msg, kv...) }
func (l *Logger) Warn(msg string, kv ...any)  { l.Logf(LevelWarn, msg, kv...) }
func (l *Logger) Error(msg string, kv ...any) { l.Logf(LevelError, msg, kv...) }

func flatten(kv []any) string {
	if len(kv) == 0 {
		return ""
	}
	out := ""
	for i := 0; i+1 < len(kv); i += 2 {
		out += fmt.Sprintf(" %s=%v", kv[i], kv[i+1])
	}
	return out
}
