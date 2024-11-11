package termites_store

type RecordId string

type Store[A any] interface {
	Put(record A) (RecordId, error)
	Get(id RecordId) (A, error)
	GetAll() []A
	Clear() error
}
