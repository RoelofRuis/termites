package termites

import (
	"errors"
	"reflect"
)

var SkipElement = errors.New("skip element")

type Adapter struct {
	name        string
	transform   func(interface{}) (interface{}, error)
	inDataType  reflect.Type
	outDataType reflect.Type
}

func NewAdapter[A any, B any](
	name string,
	transform func(A) (B, error),
) *Adapter {
	untypedTransform, inDataType, outDataType := extractFunc(transform)

	return &Adapter{
		name:        name,
		inDataType:  inDataType,
		outDataType: outDataType,
		transform:   untypedTransform,
	}
}

func (a *Adapter) ref() AdapterRef {
	info, err := determineFunctionInfo(a.transform)
	if err != nil {
		info = FunctionInfo{}
	}

	return AdapterRef{
		Name:          a.name,
		TransformInfo: info,
	}
}
