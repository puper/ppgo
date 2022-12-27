package log

import (
	"bytes"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Log struct {
	logs map[string]*logrus.Logger
}

type Config map[string]map[string]string

type JsonFormatter struct {
	logrus.JSONFormatter
}

func (me *JsonFormatter) For(entry *logrus.Entry) ([]byte, error) {
	b, err := me.JSONFormatter.Format(entry)
	if err == nil {
		buf := bytes.NewBufferString(time.Now().Format("2006-01-02 15:04:05.000") + " ")
		buf.WriteString(entry.Level.String() + " ")
		buf.Write(b)
		return buf.Bytes(), nil
	}
	return b, err
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
			l.Formatter = &JsonFormatter{}
		} else {
			l.Formatter = &logrus.TextFormatter{}
		}
		if config["out"] != "std" {
			maxAge, _ := strconv.Atoi(config["maxage"])
			if config["maxage"] == "" {
				maxAge = 10
			}
			maxSize, _ := strconv.Atoi(config["maxsize"])
			if config["maxsize"] == "" {
				maxSize = 500
			}
			maxBackups, _ := strconv.Atoi(config["maxbackups"])
			if config["maxbackups"] == "" {
				maxBackups = 10
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
			stdout, _ := strconv.Atoi(config["skipstdout"])
			if stdout <= 0 {
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
