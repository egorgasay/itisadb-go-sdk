package itisadb

type Result[V any] struct {
	value V
	err   error
}

func (r Result[V]) Value() V {
	return r.value
}

func (r Result[V]) ValueAndErr() (V, error) {
	return r.value, r.err
}

func (r Result[V]) Err() error {
	return r.err
}

type RAM struct {
	Total     uint64
	Available uint64
}
