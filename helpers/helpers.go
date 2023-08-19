package helpers

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/go-sql-driver/mysql"
)

func FileTraceInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func IsDuplicateKeyError(err error) bool {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
		return true
	}
	return false
}

func IsNull(j json.RawMessage) bool {
	var t interface{}
	err := json.Unmarshal(j, &t)
	if err != nil {
		return false
	}
	return t == nil
}
