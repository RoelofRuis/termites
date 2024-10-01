package termites

import (
	"fmt"
	"reflect"
	"runtime"
)

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
