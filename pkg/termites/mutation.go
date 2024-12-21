package termites

import (
	"fmt"
	"reflect"
)

// Mutation defines a mutation on a state S.
type Mutation[S any] interface {
	Mutate(state S) error
}

// AsMutationFor ensures that a Mutation message for state S can be passed to a node with receiving type Mutation[S].
func AsMutationFor[S any]() ConnectionOption {
	return Via(func(a any) (Mutation[S], error) {
		cast, ok := a.(Mutation[S])
		if !ok {
			return nil, fmt.Errorf("value is not a Mutation[%s]", reflect.TypeFor[S]().Name())
		}
		return cast, nil
	})
}
