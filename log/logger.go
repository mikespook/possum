// Package log provides a global logger using zerolog with configurable levels and output destinations.
package log

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	zerolog.Logger
}

// New creates a configured Logger instance based on the provided Config.
// Handles file output setup and log level configuration.
func New(config *Config) (logger Logger) {
	var logFile *os.File
	if config.Filename == "" {
		logFile = os.Stderr
	} else {
		var err error
		logFile, err = os.OpenFile(config.Filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			logFile = os.Stderr
			log.Error().Err(err).Msgf("Failed to open log file %s, using stderr instead", config.Filename)
		}
	}
	if config.Level == "" {
		config.Level = zerolog.LevelTraceValue
	}
	logLevel, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		logLevel = zerolog.TraceLevel
		log.Error().Err(err).Msgf("Failed to parse log level %s, using `trace` instead", config.Level)
	}
	logger.Logger = zerolog.New(logFile).Level(logLevel).With().Timestamp().Logger()
	return
}

// Init initializes the global logger with the provided configuration.
func init() {
	logger = New(defaultConfig)
}

func Init(config *Config) {
	logger = New(config)
}

// Logger is the global logger.
var logger Logger

// Output creates a new logger instance with the specified output writer.
func Output(w io.Writer) zerolog.Logger {
	return logger.Output(w)
}

// With creates a child logger with additional contextual fields.
func With() zerolog.Context {
	return logger.With()
}

// Level creates a child logger with the specified minimum log level.
func Level(level zerolog.Level) zerolog.Logger {
	return logger.Level(level)
}

// Sample creates a logger with the specified sampling configuration.
func Sample(s zerolog.Sampler) zerolog.Logger {
	return logger.Sample(s)
}

// Hook creates a logger with the specified hook for additional processing.
func Hook(h zerolog.Hook) zerolog.Logger {
	return logger.Hook(h)
}

// Err starts a new log message at error level, including the error if not nil.
// with info level if err is nil.
//
// You must call Msg on the returned event in order to send the event.
func Err(err error) *zerolog.Event {
	return logger.Err(err)
}

// Trace starts a new log message at trace level.
//
// You must call Msg on the returned event in order to send the event.
func Trace() *zerolog.Event {
	return logger.Trace()
}

// Debug starts a new log message at debug level.
//
// You must call Msg on the returned event in order to send the event.
func Debug() *zerolog.Event {
	return logger.Debug()
}

// Info starts a new log message at info level.
//
// You must call Msg on the returned event in order to send the event.
func Info() *zerolog.Event {
	return logger.Info()
}

// Warn starts a new log message at warn level.
//
// You must call Msg on the returned event in order to send the event.
func Warn() *zerolog.Event {
	return logger.Warn()
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func Error() *zerolog.Event {
	return logger.Error()
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method.
//
// You must call Msg on the returned event in order to send the event.
func Fatal() *zerolog.Event {
	return logger.Fatal()
}

// Panic starts a new message with panic level. The message is also sent
// to the panic function.
//
// You must call Msg on the returned event in order to send the event.
func Panic() *zerolog.Event {
	return logger.Panic()
}

// WithLevel starts a new message with level.
//
// You must call Msg on the returned event in order to send the event.
func WithLevel(level zerolog.Level) *zerolog.Event {
	return logger.WithLevel(level)
}

// Log starts a new message with no level. Setting zerolog.GlobalLevel to
// zerolog.Disabled will still disable events produced by this method.
//
// You must call Msg on the returned event in order to send the event.
func Log() *zerolog.Event {
	return logger.Log()
}

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	logger.Debug().CallerSkipFrame(1).Msg(fmt.Sprint(v...))
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	logger.Debug().CallerSkipFrame(1).Msgf(format, v...)
}

// Ctx returns the Logger associated with the ctx. If no logger
// is associated, a disabled logger is returned.
func Ctx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}
