package logger

import (
	"encoding/json"
	"io"
	"log"
	"sync"
	"time"
)

// Logger provides structured JSON logging with levels.
type Logger struct {
	mu sync.Mutex
	l  *log.Logger
}

// New constructs a Logger that writes JSON entries to the supplied writer.
func New(w io.Writer) *Logger {
	return &Logger{
		l: log.New(w, "", 0),
	}
}

// Info logs an informational message.
func (l *Logger) Info(msg string, fields map[string]any) {
	l.log("INFO", msg, fields)
}

// Error logs an error message.
func (l *Logger) Error(msg string, fields map[string]any) {
	l.log("ERROR", msg, fields)
}

func (l *Logger) log(level, msg string, fields map[string]any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry := map[string]any{
		"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		"level":     level,
		"message":   msg,
	}
	for k, v := range fields {
		entry[k] = v
	}
	data, err := json.Marshal(entry)
	if err != nil {
		l.l.Printf(`{"timestamp":"%s","level":"ERROR","message":"logger marshal error","error":"%v"}`, entry["timestamp"], err)
		return
	}
	l.l.Println(string(data))
}
