package golog

import (
	"fmt"
	"os"
	"strings"
	"sync/atomic"
)

// NameSeperator is the character used to seperate log name parts.
const NameSeperator = ":"

// Wildcard is used to match multiple logs.
const Wildcard = "*"

// The following constants are conventional log levels:
const (
	DebugLevel = "DEBUG"
	InfoLevel  = "INFO"
	WarnLevel  = "WARN"
	ErrorLevel = "ERROR"
)

var manager LogManager

// Value represents a named value that will be logged.
type Value struct {
	Name  string
	Value interface{}
}

// Val is shorthand for creating a value instance.
func Val(name string, value interface{}) Value {
	return Value{
		Name:  name,
		Value: value,
	}
}

// Logger is used to write log messages.
type Logger struct {
	Name   string
	writer atomic.Value
}

// LogWriterFunc is a function to write log messages.
type LogWriterFunc func(name, msg string, values []Value)

// GetWriter gets the writer func for this logger.
func (l *Logger) GetWriter() LogWriterFunc {
	return l.writer.Load().(LogWriterFunc)
}

// SetWriter sets the writer func for this logger.
func (l *Logger) SetWriter(writer LogWriterFunc) {
	l.writer.Store(writer)
}

// Log writes a log message to the logger.
func (l *Logger) Log(msg string, values ...Value) {
	w := l.GetWriter()
	if w != nil {
		w(l.Name, msg, values)
	}
}

// ConsoleWriter is a LogWriterFunc which logs a message to stderr.
func ConsoleWriter(name, msg string, values []Value) {
	fmt.Fprintf(os.Stderr, "%s    \t%s\n", name, msg)
	for _, v := range values {
		fmt.Fprintf(os.Stderr, "\t%s    \t%v\n", v.Name, v.Value)
	}
}

// SetWriter sets the writer for logs with names matching the given pattern.
// Log names are conventionally colon-seperated identifiers, e.g.
// DEBUG:mymodule:mycomponent.  Multiple logs can be selected using an asterisk,
// e.g. DEBUG:*.
func SetWriter(writer LogWriterFunc, filter ...string) {
	manager.SetWriter(writer, MakeNameFilter(filter...))
}

// NewLogger makes a new logger and registers it with the log manager.
func NewLogger(name ...string) *Logger {
	return manager.NewLogger(strings.Join(name, NameSeperator))
}
