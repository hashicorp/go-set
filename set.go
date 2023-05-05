// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package set provides a basic generic set implementation.
//
// https://en.wikipedia.org/wiki/Set_(mathematics)
package set

import (
	"fmt"
	"sort"
)

type nothing struct{}

var sentinel = nothing{}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// New creates a new Set with initial underlying capacity of size.
//
// A Set will automatically grow or shrink its capacity as items are added or
// removed.
//
// T must *not* be of pointer type, nor contain pointer fields, which are comparable
// but not in the way you expect. For these types, use HashSet instead.
func New[T comparable](size int) *Set[T] {
	return &Set[T]{
		items: make(map[T]nothing, max(0, size)),
	}
}

// From creates a new Set containing each item in items.
//
// T must *not* be of pointer type, nor contain pointer fields, which are comparable
// but not in the way you expect. For these types, use HashSet instead.
func From[T comparable](items []T) *Set[T] {
	s := New[T](len(items))
	s.InsertSlice(items)
	return s
}

// FromFunc creates a new Set containing a conversion of each item in items.
//
// T must *not* be of pointer type, nor contain pointer fields, which are comparable
// but not in the way you expect. For these types, use HashSet instead.
func FromFunc[A any, T comparable](items []A, conversion func(A) T) *Set[T] {
	s := New[T](len(items))
	for _, item := range items {
		s.Insert(conversion(item))
	}
	return s
}

// Set is a simple, generic implementation of the set mathematical data structure.
// It is optimized for correctness and convenience, as a replacement for the use
// of map[interface{}]struct{}.
type Set[T comparable] struct {
	items map[T]nothing
}

// Insert item into s.
//
// Return true if s was modified (item was not already in s), false otherwise.
func (s *Set[T]) Insert(item T) bool {
	if _, exists := s.items[item]; exists {
		return false
	}
	s.items[item] = sentinel
	return true
}

// InsertAll will insert each item in items into s.
//
// Return true if s was modified (at least one item was not already in s), false otherwise.
//
// Deprecated: use InsertSlice instead.
func (s *Set[T]) InsertAll(items []T) bool {
	return s.InsertSlice(items)
}

// InsertSlice will insert each item in items into s.
//
// Return true if s was modified (at least one item was not already in s), false otherwise.
func (s *Set[T]) InsertSlice(items []T) bool {
	modified := false
	for _, item := range items {
		if s.Insert(item) {
			modified = true
		}
	}
	return modified
}

// InsertSet will insert each element of o into s.
//
// Return true if s was modified (at least one item of o was not already in s), false otherwise.
func (s *Set[T]) InsertSet(o *Set[T]) bool {
	modified := false
	for item := range o.items {
		if s.Insert(item) {
			modified = true
		}
	}
	return modified
}

// Remove will remove item from s.
//
// Return true if s was modified (item was present), false otherwise.
func (s *Set[T]) Remove(item T) bool {
	if _, exists := s.items[item]; !exists {
		return false
	}
	delete(s.items, item)
	return true
}

// RemoveAll will remove each item in items from s.
//
// Return true if s was modified (any item was present), false otherwise.
//
// Deprecated: use RemoveSlice instead.
func (s *Set[T]) RemoveAll(items []T) bool {
	return s.RemoveSlice(items)
}

// RemoveSlice will remove each item in items from s.
//
// Return true if s was modified (any item was present), false otherwise.
func (s *Set[T]) RemoveSlice(items []T) bool {
	modified := false
	for _, item := range items {
		if s.Remove(item) {
			modified = true
		}
	}
	return modified
}

// RemoveSet will remove each element of o from s.
//
// Return true if s was modified (any item of o was present in s), false otherwise.
func (s *Set[T]) RemoveSet(o *Set[T]) bool {
	modified := false
	for item := range o.items {
		if s.Remove(item) {
			modified = true
		}
	}
	return modified
}

// RemoveFunc will remove each element from s that satisfies condition f.
//
// Return true if s was modified, false otherwise.
func (s *Set[T]) RemoveFunc(f func(T) bool) bool {
	modified := false
	for item := range s.items {
		if applies := f(item); applies {
			s.Remove(item)
			modified = true
		}
	}
	return modified
}

// Contains returns whether item is present in s.
func (s *Set[T]) Contains(item T) bool {
	_, exists := s.items[item]
	return exists
}

// ContainsAll returns whether s contains at least every item in items.
func (s *Set[T]) ContainsAll(items []T) bool {
	if len(s.items) < len(items) {
		return false
	}
	for _, item := range items {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// ContainsSlice returns whether s contains the same set of of elements
// that are in items. The elements of items may contain duplicates.
//
// If the slice is known to be set-like (no duplicates), EqualSlice provides
// a more efficient implementation.
func (s *Set[T]) ContainsSlice(items []T) bool {
	return s.Equal(From(items))
}

// Subset returns whether o is a subset of s.
func (s *Set[T]) Subset(o *Set[T]) bool {
	if len(s.items) < len(o.items) {
		return false
	}
	for item := range o.items {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// Size returns the cardinality of s.
func (s *Set[T]) Size() int {
	return len(s.items)
}

// Empty returns true if s contains no elements, false otherwise.
func (s *Set[T]) Empty() bool {
	return s.Size() == 0
}

// Union returns a set that contains all elements of s and o combined.
func (s *Set[T]) Union(o *Set[T]) *Set[T] {
	result := New[T](s.Size())
	for item := range s.items {
		result.items[item] = sentinel
	}
	for item := range o.items {
		result.items[item] = sentinel
	}
	return result
}

// Difference returns a set that contains elements of s that are not in o.
func (s *Set[T]) Difference(o *Set[T]) *Set[T] {
	result := New[T](max(0, s.Size()-o.Size()))
	for item := range s.items {
		if !o.Contains(item) {
			result.items[item] = sentinel
		}
	}
	return result
}

// Intersect returns a set that contains elements that are present in both s and o.
func (s *Set[T]) Intersect(o *Set[T]) *Set[T] {
	result := New[T](0)
	big, small := s, o
	if s.Size() < o.Size() {
		big, small = o, s
	}
	for item := range small.items {
		if big.Contains(item) {
			result.Insert(item)
		}
	}
	return result
}

// Copy creates a copy of s.
func (s *Set[T]) Copy() *Set[T] {
	result := New[T](s.Size())
	for item := range s.items {
		result.items[item] = sentinel
	}
	return result
}

// Slice creates a copy of s as a slice. Elements are in no particular order.
func (s *Set[T]) Slice() []T {
	result := make([]T, 0, s.Size())
	for item := range s.items {
		result = append(result, item)
	}
	return result
}

// List creates a copy of s as a slice.
//
// Deprecated: use Slice() instead.
func (s *Set[T]) List() []T {
	return s.Slice()
}

// String creates a string representation of s, using "%v" printf formating to transform
// each element into a string. The result contains elements sorted by their lexical
// string order.
func (s *Set[T]) String() string {
	return s.StringFunc(func(element T) string {
		return fmt.Sprintf("%v", element)
	})
}

// StringFunc creates a string representation of s, using f to transform each element
// into a string. The result contains elements sorted by their lexical string order.
func (s *Set[T]) StringFunc(f func(element T) string) string {
	l := make([]string, 0, s.Size())
	for item := range s.items {
		l = append(l, f(item))
	}
	sort.Strings(l)
	return fmt.Sprintf("%s", l)
}

// Equal returns whether s and o contain the same elements.
func (s *Set[T]) Equal(o *Set[T]) bool {
	if len(s.items) != len(o.items) {
		return false
	}

	for item := range s.items {
		if !o.Contains(item) {
			return false
		}
	}

	return true
}

// EqualSlice returns whether s and items contain the same elements.
//
// If items contains duplicates EqualSlice will return false; it is
// assumed that items is itself set-like. For comparing equality with
// a slice that may contain duplicates, use ContainsSlice.
func (s *Set[T]) EqualSlice(items []T) bool {
	if len(s.items) != len(items) {
		return false
	}
	return s.ContainsAll(items)
}

// MarshalJSON implements the json.Marshaler interface.
func (s *Set[T]) MarshalJSON() ([]byte, error) {
	return marshalJSON[T](s)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *Set[T]) UnmarshalJSON(data []byte) error {
	return unmarshalJSON[T](s, data)
}

func (s *Set[T]) ForEach(visit func(T) bool) {
	for item := range s.items {
		if !visit(item) {
			return
		}
	}
}
