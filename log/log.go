package log

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

// Level type encode the log level.
type Level int

// The following constants represent all possible values for the log level.
const (
	DebugLevel Level = iota
	InfoLevel
	ErrorLevel
	Quiet
)

// Logger is a very basic log interface.
type Logger interface {
	Error(msg string, a ...interface{})
	Info(msg string, a ...interface{})
	Debug(msg string, a ...interface{})
}

// log type variable satisfy Logger interface.
type log struct {
	out   io.Writer
	level Level
}

// Option is a log configuration function.
type Option func(*log)

// New creates a new Logger with Options.
func New(options ...Option) Logger {
	l := new(log)
	// set defaults
	l.out = os.Stderr
	l.level = ErrorLevel
	// options
	for _, opt := range options {
		opt(l)
	}

	return l
}

// WithWriter set the log io.Writer.
func WithWriter(w io.Writer) Option {
	return func(l *log) {
		l.out = w
	}
}

// WithLevel set the log level.
func WithLevel(lev Level) Option {
	return func(l *log) {
		l.level = lev
	}
}

// WithColor indicates to use colors when logging.
func WithColor() Option {
	return func(l *log) {
		switch l.out {
		case os.Stderr:
			l.out = color.Error
		case os.Stdout:
			l.out = color.Output
		}
		// color.NoColor = false
	}
}

// WithColor indicates to not use colors when logging.
func WithoutColor() Option {
	return func(l *log) {
		switch l.out {
		case color.Error:
			l.out = os.Stderr
		case color.Output:
			l.out = os.Stdout
		}
		color.NoColor = true
	}
}

// Some color variables.
var (
	// nocolor    = func(a ...interface{}) string { return fmt.Sprint(a...) }
	msgcolor   = color.New(color.FgWhite).SprintFunc()
	errcolor   = color.New(color.FgRed, color.Bold).SprintFunc()
	infocolor  = color.New(color.FgYellow, color.Bold).SprintFunc()
	debugcolor = color.New(color.FgCyan, color.Bold).SprintFunc()
)

// printLog prints to l.out if level is high enough
// level and tag specify the type (INFO, ERROR, DEBUG).
func printLog(l *log, level Level, tag, msg string, a ...interface{}) {
	if l.level > level {
		return
	}
	fmt.Fprintln(l.out, tag, fmt.Sprintf(msgcolor(msg), a...))
}

// Error method for Logger interface.
func (l *log) Error(msg string, a ...interface{}) {
	printLog(l, ErrorLevel, errcolor("ERROR:"), msg, a...)
}

// Infof method for Logger interface.
func (l *log) Info(msg string, a ...interface{}) {
	printLog(l, InfoLevel, infocolor("INFO:"), msg, a...)
}

// Debug method for Logger interface.
func (l *log) Debug(msg string, a ...interface{}) {
	printLog(l, DebugLevel, debugcolor("DEBUG:"), msg, a...)
}
