package businesserror

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

func IsBusinessError(err error) bool {
	_, ok := err.(*BusinessError)
	return ok
}

type BusinessError struct {
	manager   *Manager
	err       error
	errConfig ErrorConfig
	fileline  string
	fields    map[string]interface{}
}

func NewBusinessError(manager *Manager, err error) *BusinessError {
	b := &BusinessError{
		manager:   manager,
		err:       err,
		errConfig: manager.config.Errors[manager.config.DefaultErrCode],
	}
	b.fileline = FileTraceInfo(3)
	b.fields = map[string]interface{}{}
	return b
}

func (this *BusinessError) Error() string {
	return fmt.Sprintf("%v:%v;%v", this.errConfig.Code, this.errConfig.ErrMsg, this.err.Error())
}

func (this *BusinessError) ParamError() *BusinessError {
	return this
}

func (this *BusinessError) ErrCode(errCode string) *BusinessError {
	if errConfig, ok := this.manager.config.Errors[errCode]; ok {
		this.errConfig = errConfig
	} else {
		this.errConfig.Code = errCode
	}
	return this
}

func (this *BusinessError) StatusCode(statusCode int) *BusinessError {
	this.errConfig.StatusCode = statusCode
	return this
}

func (this *BusinessError) OutMsg(format string, a ...interface{}) *BusinessError {
	this.errConfig.ErrMsg = fmt.Sprintf(format, a...)
	return this
}

func (this *BusinessError) Log(logger *logrus.Logger) *BusinessError {
	logger.WithField("fileline", this.fileline).WithFields(this.fields).WithError(this.err).Error(this.errConfig.ErrMsg)
	return this
}

func (this *BusinessError) ErrorConfig() ErrorConfig {
	return this.errConfig
}

func (this *BusinessError) Arg(v interface{}) *BusinessError {
	this.fields["arg"] = v
	return this
}

func (this *BusinessError) Reply(v interface{}) *BusinessError {
	this.fields["reply"] = v
	return this
}

func (this *BusinessError) Cost(start time.Time) *BusinessError {
	this.fields["cost"] = time.Now().Sub(start)
	return this
}

func (this *BusinessError) Field(name string, value interface{}) *BusinessError {
	this.fields[name] = value
	return this
}
