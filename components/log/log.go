package log

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Log struct {
	logs map[string]*logrus.Logger
}

type Config map[string]map[string]string

type LogConfig struct {
	Level         string        `json:"level"`
	TraceLevel    string        `json:"traceLevel"`
	Output        string        `json:"output"`
	MaxSize       int           `json:"maxSize"`
	MaxAge        int           `json:"maxAge"`
	MaxBackups    int           `json:"maxBackups"`
	Compress      bool          `json:"compress"`
	InitialFields []interface{} `json:"initialFields"`
	Format        string        `json:"format"`
}

func New(cfg *Config) (*Log, error) {
	instance := &Log{
		logs: make(map[string]*logrus.Logger),
	}
	instance.logs["default"] = logrus.StandardLogger()
	for name, config := range *cfg {
		l := logrus.New()
		level, err := logrus.ParseLevel(config["level"])
		if err != nil {
			return nil, err
		}
		l.Level = level
		if config["format"] == "json" {
			l.Formatter = &logrus.JSONFormatter{}
		} else {
			l.Formatter = &logrus.TextFormatter{}
		}
		if config["out"] != "std" {
			maxAge, _ := strconv.Atoi(config["maxage"])
			if config["maxage"] == "" {
				maxAge = 30
			}
			maxSize, _ := strconv.Atoi(config["maxsize"])
			if config["maxsize"] == "" {
				maxSize = 1000
			}
			maxBackups, _ := strconv.Atoi(config["maxbackups"])
			if config["maxbackups"] == "" {
				maxBackups = 30
			}
			w := &lumberjack.Logger{
				Filename:   config["out"],
				MaxSize:    maxSize, // megabytes
				MaxAge:     maxAge,  // days
				MaxBackups: maxBackups,
				LocalTime:  true,
				Compress:   true,
			}
			l.Out = w
			if config["stdout"] != "true" {
				l.AddHook(&writer.Hook{
					Writer: os.Stdout,
					LogLevels: []logrus.Level{
						logrus.PanicLevel,
						logrus.FatalLevel,
						logrus.ErrorLevel,
						logrus.WarnLevel,
						logrus.InfoLevel,
						logrus.DebugLevel,
						logrus.TraceLevel,
					},
				})
			}
		}
		instance.logs[name] = l
	}
	return instance, nil
}

func (this *Log) Get(name string) *logrus.Logger {
	l, ok := this.logs[name]
	if ok {
		return l
	}
	return this.logs["default"]
}
