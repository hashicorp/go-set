// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package set

// Common is a minimal interface that all sets implement.
type Common[T any] interface {

	// Slice returns a slice of all elements in the set.
	//
	// Note: order of elements depends on the underlying implementation.
	Slice() []T

	// Insert an element into the set.
	//
	// Returns true if the set is modified as a result.
	Insert(T) bool

	// InsertSlice inserts all elements from the slice into the set.
	//
	// Returns true if the set was modified as a result.
	InsertSlice([]T) bool

	// Size returns the number of elements in the set.
	Size() int

	// ForEach will call the callback function for each element in the set.
	// If the callback returns false, the iteration will stop.
	//
	// Note: iteration order depends on the underlying implementation.
	ForEach(func(T) bool)
}

// InsertSliceFunc inserts all elements from the slice into the set
func InsertSliceFunc[T, E any](s Common[T], items []E, f func(element E) T) {
	for _, item := range items {
		s.Insert(f(item))
	}
}

// InsertSetFunc inserts the elements of a into b, applying the transform function
// to each element before insertion.
//
// Returns true if b was modified as a result.
func InsertSetFunc[T, E any](a Common[T], b Common[E], transform func(T) E) bool {
	modified := false
	a.ForEach(func(item T) bool {
		if b.Insert(transform(item)) {
			modified = true
		}
		return true
	})
	return modified
}

// SliceFunc produces a slice of the elements in s, applying the transform
// function to each element first.
func SliceFunc[T, E any](s Common[T], transform func(T) E) []E {
	slice := make([]E, 0, s.Size())
	s.ForEach(func(item T) bool {
		slice = append(slice, transform(item))
		return true
	})
	return slice
}
