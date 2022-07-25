package maybe

import "errors"

var _ Maybe[int, string] = Just[int, string]{}
var _ Maybe[int, string] = Nothing[int, string]{}

var (
	errCouldntAssert = errors.New("unable to type assert as pointer")
)

func Of[T, U any](value *T) Maybe[T, U] {
	if value != nil {
		return Just[T, U]{
			start: value,
		}
	}
	return Nothing[T, U]{}
}

type Maybe[T, U any] interface {
	Map(fn func(*T) *U) Maybe[T, U]
	Get() any
}

type Just[T, U any] struct {
	start       *T
	next        *U
	hasSwitched bool
}

func (j Just[T, U]) Get() any {
	if j.hasSwitched {
		return j.next
	}
	return j.start
}

func (j Just[T, U]) Map(fn func(*T) *U) Maybe[T, U] {
	// once it hasSwitched
	switch j.hasSwitched {
	case true:
		return Nothing[T, U]{}
	}
	if j.start != nil {
		return Just[T, U]{
			start:       j.start,
			next:        fn(j.start),
			hasSwitched: true,
		}
	}
	return Nothing[T, U]{}
}

type Nothing[T, U any] struct{}

func (n Nothing[T, U]) Get() any {
	return n
}

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
			start: j1.next,
		}
	case Nothing[T, U]:
		next = Nothing[U, V]{}
	default:
		return Nothing[U, V]{}
	}

	return next
}

func As[T any](value any) (*T, error) {
	if v1, ok := value.(*T); !ok {
		return nil, errCouldntAssert
	} else {
		return v1, nil
	}
}
