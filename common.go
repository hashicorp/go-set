package set

type CommonSet[T any] interface {
	Slice() []T
	Insert(T) bool
	InsertSlice([]T) bool
	Size() int
}
