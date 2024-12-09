package termites

import (
	"reflect"
	"testing"
)

func TestNewInPort(t *testing.T) {
	builder := NewBuilder("Tester")

	port := NewInPort[int](builder)

	if port.dataType != reflect.TypeOf(0) {
		t.Errorf("Port has incorrect data type")
	}
	if len(port.connections) != 0 {
		t.Errorf("Port should not have any connections")
	}
}

type testType interface {
	someFunc() error
}

func TestNewInPortFromInterface(t *testing.T) {
	builder := NewBuilder("Tester")

	port := NewInPort[testType](builder)

	if port.dataType != reflect.TypeFor[testType]() {
		t.Errorf("Port has incorrect data type")
	}
	if len(port.connections) != 0 {
		t.Errorf("Port should not have any connections")
	}
}

func TestNewOutPort(t *testing.T) {
	builder := NewBuilder("Tester")

	port := NewOutPort[int](builder)

	if port.dataType != reflect.TypeOf(0) {
		t.Errorf("Port has incorrect data type")
	}
	if len(port.connections) != 0 {
		t.Errorf("Port should not have any connections")
	}
}

func TestNewOutPortFromInterface(t *testing.T) {
	builder := NewBuilder("Tester")

	port := NewOutPort[testType](builder)

	if port.dataType != reflect.TypeFor[testType]() {
		t.Errorf("Port has incorrect data type")
	}
	if len(port.connections) != 0 {
		t.Errorf("Port should not have any connections")
	}
}
