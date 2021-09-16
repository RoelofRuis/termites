package termites

type Adapter struct {
	name        string
	inDataType  string
	outDataType string
	transform   func(interface{}) (interface{}, error)
}

func NewAdapter(
	name string,
	exampleMessageIn interface{},
	exampleMessageOut interface{},
	transform func(interface{}) (interface{}, error),
) *Adapter {
	inDataType := determineDataType(exampleMessageIn)
	outDataType := determineDataType(exampleMessageOut)

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
