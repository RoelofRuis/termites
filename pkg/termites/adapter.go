package termites

import (
	"errors"
	"reflect"
)

var SkipElement = errors.New("skip element")

type adapter struct {
	name        string
	transform   func(interface{}) (interface{}, error)
	inDataType  reflect.Type
	outDataType reflect.Type
}

func (a *adapter) ref() AdapterRef {
	info, err := determineFunctionInfo(a.transform)
	if err != nil {
		info = FunctionInfo{}
	}

	return AdapterRef{
		Name:          a.name,
		TransformInfo: info,
	}
}
