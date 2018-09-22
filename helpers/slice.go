package helpers

import (
	"fmt"
	"reflect"
)

func SlicePluck(key string, in, out interface{}) error {
	inVal := reflect.Indirect(reflect.ValueOf(in))
	outVal := reflect.Indirect(reflect.ValueOf(out))
	if inVal.Kind() != reflect.Slice || outVal.Kind() != reflect.Slice {
		return fmt.Errorf("in, out not a slice")
	}
	for i := 0; i < inVal.Len(); i++ {
		outVal.Set(reflect.Append(outVal, reflect.Indirect(inVal.Index(i)).FieldByName(key)))
	}
	return nil
}
