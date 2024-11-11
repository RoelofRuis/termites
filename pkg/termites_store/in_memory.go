package termites_store

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

// InMemory is a thread safe in-memory Store implementation.
type InMemory[A any] struct {
	mx     sync.RWMutex
	lastId int
	rows   []A
}

func NewInMemory[A any]() *InMemory[A] {
	return &InMemory[A]{lastId: 0}
}

func (i *InMemory[A]) Put(record A) (RecordId, error) {
	i.mx.Lock()
	i.rows = append(i.rows, record)
	recordId := RecordId(fmt.Sprintf("%d", i.lastId))
	i.lastId++
	i.mx.Unlock()

	return recordId, nil
}

func (i *InMemory[A]) Get(id RecordId) (A, error) {
	var val A
	intId, err := strconv.ParseInt(string(id), 10, 64)
	if err != nil {
		return val, errors.New("invalid id")
	}
	i.mx.RLock()
	defer i.mx.RUnlock()
	if int(intId) >= len(i.rows) {
		return val, errors.New("id out of range")
	}
	return i.rows[int(intId)], nil
}

func (i *InMemory[A]) GetAll() []A {
	i.mx.RLock()
	res := make([]A, len(i.rows))
	copy(res, i.rows)
	i.mx.RUnlock()

	return res
}

func (i *InMemory[A]) Clear() error {
	i.mx.Lock()
	i.rows = []A{}
	i.mx.Unlock()

	return nil
}
