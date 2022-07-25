package maybe

var _ Maybe[int, string] = Just[int, string]{}
var _ Maybe[int, string] = Nothing[int, string]{}

type Maybe[T, U any] interface {
	Map(fn func(*T) *U) Maybe[T, U]
}

type Just[T, U any] struct {
	Start       *T
	Next        *U
	hasSwitched bool
}

func (j Just[T, U]) Map(fn func(*T) *U) Maybe[T, U] {
	// once it hasSwitched
	switch j.hasSwitched {
	case true:
		return Nothing[T, U]{}
	}
	if j.Start != nil {
		return Just[T, U]{
			Start:       j.Start,
			Next:        fn(j.Start),
			hasSwitched: true,
		}
	}
	return Nothing[T, U]{}
}

type Nothing[T, U any] struct{}

func (n Nothing[T, U]) Map(_ func(*T) *U) Maybe[T, U] {
	return n
}

func FromMaybeToAnother[T, U, V any](m Maybe[T, U]) Maybe[U, V] {
	var next Maybe[U, V]

	switch m.(type) {
	case Just[T, U]:
		j1, ok := m.(Just[T, U])
		if !ok {
			return Nothing[U, V]{}
		}
		next = Just[U, V]{
			Start: j1.Next,
		}
	case Nothing[T, U]:
		next = Nothing[U, V]{}
	default:
		return Nothing[U, V]{}
	}

	return next
}
