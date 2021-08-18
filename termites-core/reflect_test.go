package termites

import (
	"testing"
)

type TestStruct struct{}

var typecheckTestCases = []struct {
	name   string
	object interface{}
	expect string
}{
	{
		"nil",
		nil,
		"nil",
	},
	{
		"int",
		0,
		"int",
	},
	{
		"float64",
		0.001,
		"float64",
	},
	{
		"bool",
		true,
		"bool",
	},
	{
		"struct",
		TestStruct{},
		"TestStruct",
	},
	{
		"byte array",
		[]byte{},
		"[]uint8",
	},
	{
		"anything named by string",
		"package.Anything",
		"package.Anything",
	},
	{
		"string",
		"",
		"string",
	},
	{
		"map",
		map[int]string{},
		"map[int]string",
	},
}

func TestDetermineDataType(t *testing.T) {
	for _, test := range typecheckTestCases {
		t.Run(test.name, func(t *testing.T) {
			result := determineDataType(test.object)
			if result != test.expect {
				t.Errorf("Incorrect data type for [%v]\nexpected: %s\ngot: %s", test.object, test.expect, result)
				t.FailNow()
			}
		})
	}
}
