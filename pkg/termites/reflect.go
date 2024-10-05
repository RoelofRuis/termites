package termites

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
)

func extractFunc[A any, B any](f func(A) (B, error)) (func(interface{}) (interface{}, error), reflect.Type, reflect.Type) {
	fnType := reflect.TypeOf(f)

	inType := fnType.In(0)
	outType := fnType.Out(0)

	untypedFunc := func(arg interface{}) (interface{}, error) {
		val, ok := arg.(A)
		if !ok {
			return nil, errors.New("invalid argument type")
		}

		result, err := f(val)

		return result, err
	}

	return untypedFunc, inType, outType
}

func determineFunctionInfo(f interface{}) (FunctionInfo, error) {
	if reflect.TypeOf(f).Kind() != reflect.Func {
		return FunctionInfo{}, fmt.Errorf("cannot determine info of non-function")
	}

	v := reflect.ValueOf(f)
	if v.IsNil() {
		return FunctionInfo{}, fmt.Errorf("cannot determine info of `nil` function")
	}

	// FIXME: this was broken in go 1.18 due to update to FuncForPC
	file, line := runtime.FuncForPC(v.Pointer()).FileLine(v.Pointer())

	return FunctionInfo{
		file,
		line,
	}, nil
}
