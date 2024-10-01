package termites

import "reflect"

type Adapter struct {
	name        string
	inDataType  reflect.Type
	outDataType reflect.Type
	transform   func(interface{}) (interface{}, error)
}

func NewAdapter[A any, B any](
	name string,
	exampleMessageIn A,
	exampleMessageOut B,
	transform func(interface{}) (interface{}, error),
) *Adapter {
	inDataType := reflect.TypeOf(exampleMessageIn)
	outDataType := reflect.TypeOf(exampleMessageOut)

	return &Adapter{
		name:        name,
		inDataType:  inDataType,
		outDataType: outDataType,
		transform:   transform,
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
