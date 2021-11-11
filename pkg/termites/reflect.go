package termites

import (
	"fmt"
	"reflect"
	"runtime"
)

func determineDataType(object interface{}) string {
	if object == nil {
		return "nil"
	}

	kind := reflect.TypeOf(object).Kind()
	switch kind {
	case reflect.Slice:
		return "[]" + reflect.TypeOf(object).Elem().Name()

	case reflect.Map:
		return fmt.Sprintf(
			"map[%s]%s",
			reflect.TypeOf(object).Key().Name(),
			reflect.TypeOf(object).Elem().Name(),
		)

	case reflect.String:
		strVal := object.(string)
		if strVal != "" {
			return strVal
		}
		fallthrough

	default:
		return reflect.TypeOf(object).Name()
	}
}

func determineFunctionInfo(f interface{}) (FunctionInfo, error) {
	if reflect.TypeOf(f).Kind() != reflect.Func {
		return FunctionInfo{}, fmt.Errorf("cannot determine info of non-function")
	}

	v := reflect.ValueOf(f)
	if v.IsNil() {
		return FunctionInfo{}, fmt.Errorf("cannot determine info of `nil` function")
	}

	file, line := runtime.FuncForPC(v.Pointer()).FileLine(v.Pointer())

	return FunctionInfo{
		file,
		line,
	}, nil
}