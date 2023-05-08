// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package set

// Common is the interface that all sets implement
type Common[T any] interface {
	// Slice returns a slice of all elements in the set
	Slice() []T
	// Insert inserts an element into the set
	// if the element already exists, it will return false
	Insert(T) bool
	// InsertSlice inserts all elements from the slice into the set
	InsertSlice([]T) bool
	// Size returns the number of elements in the set
	Size() int
	// ForEach  will call the callback function for each element in the set.
	// If the callback returns false, the iteration will stop.
	// Note: iteration order depends on the underlying implementation;
	ForEach(call func(T) bool)
}

// InsertSliceFunc inserts all elements from the slice into the set
func InsertSliceFunc[T, E any](s Common[T], items []E, f func(element E) T) {
	for _, item := range items {
		s.Insert(f(item))
	}
}

// TransformUnion transforms the set A into another set B
func TransformUnion[T, E any](a Common[T], b Common[E], transform func(T) E) {
	a.ForEach(func(item T) bool {
		_ = b.Insert(transform(item))
		return true
	})
}

// TransformSlice transforms the set into a slice
func TransformSlice[T, E any](s Common[T], transform func(T) E) []E {
	slice := make([]E, 0, s.Size())
	s.ForEach(func(item T) bool {
		slice = append(slice, transform(item))
		return true
	})
	return slice
}
