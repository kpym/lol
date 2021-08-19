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
	Error(msg string)
	Errorf(msg string, a ...interface{})
	Info(msg string)
	Infof(msg string, a ...interface{})
	Debug(msg string)
	Debugf(msg string, a ...interface{})
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

// Error method for Logger interface.
func (l *log) Error(msg string) {
	if l.level > ErrorLevel {
		return
	}
	fmt.Fprintln(l.out, errcolor("ERROR:"), msgcolor(msg))
}

// Errorf method for Logger interface.
func (l *log) Errorf(msg string, a ...interface{}) {
	if l.level > ErrorLevel {
		return
	}
	fmt.Fprintf(l.out, errcolor("ERROR:")+" "+msgcolor(msg), a...)
}

// Info method for Logger interface.
func (l *log) Info(msg string) {
	if l.level > InfoLevel {
		return
	}
	fmt.Fprintln(l.out, infocolor("INFO:"), msgcolor(msg))
}

// Infof method for Logger interface.
func (l *log) Infof(msg string, a ...interface{}) {
	if l.level > InfoLevel {
		return
	}
	fmt.Fprintf(l.out, infocolor("INFO:")+" "+msgcolor(msg), a...)
}

// Debug method for Logger interface.
func (l *log) Debug(msg string) {
	if l.level > DebugLevel {
		return
	}
	fmt.Fprintln(l.out, debugcolor("DEBUG:"), msgcolor(msg))
}

// Debugf method for Logger interface.
func (l *log) Debugf(msg string, a ...interface{}) {
	if l.level > DebugLevel {
		return
	}
	fmt.Fprintf(l.out, debugcolor("DEBUG:")+" "+msgcolor(msg), a...)
}
