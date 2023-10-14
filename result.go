package itisadb

type Result[V any] struct {
	val V
	err error
}

func (r Result[V]) Val() V {
	return r.val
}

func (r Result[V]) ValueAndErr() (V, error) {
	return r.val, r.err
}

func (r Result[V]) Err() error {
	return r.err
}

type RAM struct {
	Total     uint64
	Available uint64
}
