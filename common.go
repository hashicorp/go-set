package set

// CommonSet is the interface that all sets implement
type CommonSet[T any] interface {
	Slice() []T
	Insert(T) bool
	InsertSlice([]T) bool
	Size() int
	// ForEach  will call the callback function for each element in the set
	// If the callback returns false, the iteration will stop
	ForEach(call func(T) bool)
}

// InsertSliceFunc inserts all elements from the slice into the set
func InsertSliceFunc[T, E any](set CommonSet[T], items []E, f func(element E) T) {
	for _, item := range items {
		set.Insert(f(item))
	}
}

// Transform transforms the set a into another set b
func Transform[T, E any](a CommonSet[T], b CommonSet[E], transform func(T) E) {
	a.ForEach(func(item T) bool {
		_ = b.Insert(transform(item))
		return true
	})
}

// TransformSlice transforms the set into a slice
func TransformSlice[T, E any](s CommonSet[T], transform func(T) E) []E {
	slice := make([]E, 0, s.Size())
	s.ForEach(func(item T) bool {
		slice = append(slice, transform(item))
		return true
	})
	return slice
}
