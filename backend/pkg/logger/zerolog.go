package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
)

type Config struct {

	// Enable console logging
	ConsoleLoggingEnabled bool

	// EncodeLogsAsJson makes the log framework log JSON
	EncodeLogsAsJson bool

	// FileLoggingEnabled makes the framework log to a file.go
	// the fields below can be skipped if this value is false!
	FileLoggingEnabled bool

	// Directory to log to to when filelogging is enabled
	Directory string

	// Filename is the name of the logfile which will be placed inside the directory
	Filename string

	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int

	// MaxBackups the max number of rolled files to keep
	MaxBackups int

	// MaxAge the max age in days to keep a logfile
	MaxAge int

	// CallerSkip the number of directory hierarchy to be skipped
	CallerSkip int
}

type Logger struct {
	serviceLogger *zerolog.Logger
}

// InitZerolog
// Initialize Zerolog, it during cmd main.go
func InitZerolog(config Config) *Logger {
	var writers []io.Writer

	if config.ConsoleLoggingEnabled {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}
	if config.FileLoggingEnabled {
		writers = append(writers, newRollingFile(config))
	}
	if config.CallerSkip != 0 {
		zerolog.CallerSkipFrameCount = config.CallerSkip
	}
	mw := io.MultiWriter(writers...)

	logger := zerolog.New(mw).With().Timestamp().Logger()

	logger.Info().
		Bool("fileLogging", config.FileLoggingEnabled).
		Bool("jsonLogOutput", config.EncodeLogsAsJson).
		Str("logDirectory", config.Directory).
		Str("fileName", config.Filename).
		Int("maxSizeMB", config.MaxSize).
		Int("maxBackups", config.MaxBackups).
		Int("maxAgeInDays", config.MaxAge).
		Msg("logging configured")

	return &Logger{
		serviceLogger: &logger,
	}
}

func newRollingFile(config Config) io.Writer {
	l := &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxBackups: config.MaxBackups, // files
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
	}

	return l
}

func (l *Logger) GetServiceLogger() *zerolog.Logger {
	return l.serviceLogger
}

// Errorf print formatted error message
func (l *Logger) Errorf(format string, a ...interface{}) {
	errs := fmt.Errorf(format, a...)
	l.serviceLogger.Error().Caller().Msgf(errs.Error())
}

// Error print error message
func (l *Logger) Error(msg string) {
	l.serviceLogger.Error().Caller().Msgf(msg)
}

// Infof
// Print log message with format
func (l *Logger) Infof(format string, a ...interface{}) {
	l.serviceLogger.Info().Msgf(format, a...)
}

// Info
// Print log message
func (l *Logger) Info(format string) {
	l.serviceLogger.Info().Msg(format)
}

// Fatalf
// Stop the app after invocation, and
// print Fatal message with format
func (l *Logger) Fatalf(format string, a ...interface{}) {
	l.serviceLogger.Fatal().Msgf(format, a...)
}

// Fatal
// Stop the app after invocation,
// and print the Fatal message
func (l *Logger) Fatal(format string) {
	l.serviceLogger.Fatal().Msg(format)
}

// Panic
// Stop the app after invocation,
// and print the Panic message
func (l *Logger) Panic(format string) {
	l.serviceLogger.Panic().Msg(format)
}
