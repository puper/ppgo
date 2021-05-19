package helpers

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func FileTraceInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func FileLogFields() logrus.Fields {
	return logrus.Fields{
		"fileline": FileTraceInfo(2),
	}
}

func IsDuplicateKeyError(err error) bool {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
		return true
	}
	return false
}

func WrapError(logger *logrus.Logger, err error, args interface{}, context interface{}, msg string) error {
	if err == nil {
		err = errors.Errorf(msg)
	} else {
		err = errors.Wrap(err, msg)
	}
	if msg == "" {
		msg = err.Error()
	}
	if logger != nil {
		logger.WithField("fileline", FileTraceInfo(3)).WithField("args", args).WithField("context", context).WithError(err).Error(msg)
	}
	return err
}

func IsNull(j json.RawMessage) bool {
	var t interface{}
	err := json.Unmarshal(j, &t)
	if err != nil {
		return false
	}
	return t == nil
}
