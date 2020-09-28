package zerologcpt

import (
	"io"
	"os"
	"path"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	levels = map[string]zerolog.Level{
		"debug": zerolog.DebugLevel,
		"info":  zerolog.InfoLevel,
		"warn":  zerolog.WarnLevel,
		"error": zerolog.ErrorLevel,
		"fatal": zerolog.FatalLevel,
		"panic": zerolog.PanicLevel,
		"":      zerolog.NoLevel,
	}
)

// Configuration for logging
type Config struct {
	// Enable console logging
	ConsoleLoggingEnabled bool `json:"consoleLoggingEnabled"`

	// EncodeLogsAsJson makes the log framework log JSON
	EncodeLogsAsJson bool `json:"encodeLogsAsJson"`
	// FileLoggingEnabled makes the framework log to a file
	// the fields below can be skipped if this value is false!
	FileLoggingEnabled bool `json:"fileLoggingEnabled"`
	// Directory to log to to when filelogging is enabled
	Directory string `json:"directory"`
	// Filename is the name of the logfile which will be placed inside the directory
	Filename string `json:"filename"`
	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int `json:"maxSize"`
	// MaxBackups the max number of rolled files to keep
	MaxBackups int `json:"maxBackups"`
	// MaxAge the max age in days to keep a logfile
	MaxAge int    `json:"maxAge"`
	Level  string `json:"level"`
}

type Logger struct {
	*zerolog.Logger
	mu         sync.Mutex
	fileCloser io.Closer
}

func (this *Logger) Close() error {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.fileCloser != nil {
		err := this.fileCloser.Close()
		if err == nil {
			this.fileCloser = nil
			return nil
		}
		return err
	}
	return nil
}

// Configure sets up the logging framework
//
// In production, the container logs will be collected and file logging should be disabled. However,
// during development it's nicer to see logs as text and optionally write to a file when debugging
// problems in the containerized pipeline
//
// The output log file will be located at /var/log/service-xyz/service-xyz.log and
// will be rolled according to configuration set.
func New(config Config) *Logger {
	var writers []io.Writer

	if config.ConsoleLoggingEnabled {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}
	if config.FileLoggingEnabled {
		writers = append(writers, newRollingFile(config))
	}
	mw := io.MultiWriter(writers...)

	// zerolog.SetGlobalLevel(zerolog.DebugLevel)
	logger := zerolog.New(mw).With().Timestamp().Logger()
	if lvl, ok := levels[config.Level]; ok {
		logger = logger.Level(lvl)
	}
	logger.Info().
		Bool("fileLogging", config.FileLoggingEnabled).
		Bool("jsonLogOutput", config.EncodeLogsAsJson).
		Str("logDirectory", config.Directory).
		Str("fileName", config.Filename).
		Int("maxSizeMB", config.MaxSize).
		Int("maxBackups", config.MaxBackups).
		Int("maxAgeInDays", config.MaxAge).
		Str("level", config.Level).
		Msg("logging configured")

	return &Logger{
		Logger: &logger,
	}
}

func newRollingFile(config Config) io.Writer {
	if err := os.MkdirAll(config.Directory, 0744); err != nil {
		log.Error().Err(err).Str("path", config.Directory).Msg("can't create log directory")
		return nil
	}

	return &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxBackups: config.MaxBackups, // files
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
	}
}
