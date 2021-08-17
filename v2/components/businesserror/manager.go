package businesserror

import (
	"fmt"
	"runtime"

	pplog "github.com/puper/ppgo/v2/components/log"
	"github.com/puper/ppgo/v2/engine"
	"github.com/sirupsen/logrus"
)

func Builder(configBuilder func() *Config) func(*engine.Engine) (interface{}, error) {
	return func(e *engine.Engine) (interface{}, error) {
		cfg := configBuilder()
		m := New(cfg)
		m.logger = e.Get("log").(*pplog.Log).Get("")
		return m, nil
	}
}

func FileTraceInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	}
	return fmt.Sprintf("%s:%d", file, line)
}

type ErrorConfig struct {
	StatusCode int    `json:"statusCode,omitempty"`
	Code       string `json:"code,omitempty"`
	ErrMsg     string `json:"errMsg,omitempty"`
}

type Config struct {
	DefaultErrCode string                 `json:"defaultErrCode,omitempty"`
	Errors         map[string]ErrorConfig `json:"errors,omitempty"`
}

type Manager struct {
	config *Config
	logger *logrus.Logger
}

func (this *Manager) GetDefaultErrConfig() ErrorConfig {
	return this.config.Errors[this.config.DefaultErrCode]
}

func (this *Manager) SetDefaultErrCode(errCode string) {
	this.config.DefaultErrCode = errCode
}

func (this *Manager) AddErrorConfigs(errConfigs []ErrorConfig) {
	for _, errConfig := range errConfigs {
		this.config.Errors[errConfig.Code] = errConfig
	}
}

func New(config *Config) *Manager {
	if config.Errors == nil {
		config.Errors = map[string]ErrorConfig{}
	}
	return &Manager{
		config: config,
	}
}

func (this *Manager) FromError(err error) *BusinessError {
	if IsBusinessError(err) {
		return err.(*BusinessError)
	}
	return NewBusinessError(this, err)
}

func (this *Manager) FromString(format string, a ...interface{}) *BusinessError {
	return NewBusinessError(
		this,
		fmt.Errorf(format, a...),
	)
}
