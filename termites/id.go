package termites

import (
	"fmt"
	"math/rand"
)

type ObjectName string

type Identifier struct {
	ObjName ObjectName
	Id      string
}

func (o Identifier) String() string {
	return fmt.Sprintf("%s-%s", o.ObjName, o.Id)
}

func NewIdentifier(n ObjectName) Identifier {
	return Identifier{
		ObjName: n,
		Id:      RandomID(),
	}
}

var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandomID() string {
	b := make([]rune, 16)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}
