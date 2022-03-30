package zaplog

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Log struct {
	logs map[string]*zap.SugaredLogger
}

type Config map[string]*LogConfig

func New(cfg Config) (*Log, error) {
	instance := &Log{
		logs: make(map[string]*zap.SugaredLogger),
	}
	for logName, logCfg := range cfg {
		instance.logs[logName] = NewLog(logCfg)
	}
	if _, ok := instance.logs["default"]; !ok {
		instance.logs["default"] = NewLog(&LogConfig{})
	}
	return instance, nil
}

func (this *Log) Get(name string) *zap.SugaredLogger {
	l, ok := this.logs[name]
	if ok {
		return l
	}
	return this.logs["default"]
}

func NewLog(logCfg *LogConfig) *zap.SugaredLogger {
	levelMap := map[string]zapcore.Level{
		"debug":  zapcore.DebugLevel,
		"info":   zapcore.InfoLevel,
		"warn":   zapcore.WarnLevel,
		"error":  zapcore.ErrorLevel,
		"dpanic": zapcore.DPanicLevel,
		"panic":  zapcore.PanicLevel,
		"fatal":  zapcore.FatalLevel,
	}
	cfg := zap.NewProductionConfig()
	if lvl, ok := levelMap[logCfg.Level]; ok {
		cfg.Level = zap.NewAtomicLevelAt(lvl)
	} else {
		cfg.Level = zap.NewAtomicLevelAt(levelMap["info"])
	}
	cfg.EncoderConfig.LineEnding = zapcore.DefaultLineEnding
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	var enc zapcore.Encoder
	if logCfg.Format == "json" {
		enc = zapcore.NewJSONEncoder(cfg.EncoderConfig)
	} else {
		enc = zapcore.NewConsoleEncoder(cfg.EncoderConfig)
	}
	var w zapcore.WriteSyncer
	if logCfg.Output == "" {
		w = zapcore.AddSync(os.Stderr)
	} else {
		w = zapcore.AddSync(&lumberjack.Logger{
			Filename:   logCfg.Output,
			MaxSize:    logCfg.MaxSize, // megabytes
			MaxAge:     logCfg.MaxAge,  // days
			MaxBackups: logCfg.MaxBackups,
			LocalTime:  true,
			Compress:   logCfg.Compress,
		})
	}
	log := zap.New(
		zapcore.NewCore(enc, w, cfg.Level),
	)
	opts := []zap.Option{}
	opts = append(opts, zap.AddCaller())
	if lvl, ok := levelMap[logCfg.TraceLevel]; ok {
		opts = append(opts, zap.AddStacktrace(lvl))
	} else {
		opts = append(opts, zap.AddStacktrace(zap.ErrorLevel))
	}
	log = log.WithOptions(opts...)
	return log.Sugar().With(logCfg.InitialFields...)
}

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
