// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package set

import (
	"fmt"
	"sort"
)

// Hash represents the output type of a Hash() function defined on a type.
//
// A Hash could be string-like or int-like. A string hash could be something like
// and md5, sha1, or GoString() representation of a type. An int hash could be
// something like the prime multiple hash code of a type.
type Hash interface {
	~string | ~int | ~uint | ~int64 | ~uint64 | ~int32 | ~uint32 | ~int16 | ~uint16 | ~int8 | ~uint8
}

// HashFunc is a generic type constraint for any type that implements a Hash()
// method with a Hash return type.
type HashFunc[H Hash] interface {
	Hash() H
}

// HashSet is a generic implementation of the mathematical data structure, oriented
// around the use of a HashFunc to make hash values from other types.
type HashSet[T HashFunc[H], H Hash] struct {
	items map[H]T
}

// NewHashSet creates a HashSet with underlying capacity of size.
//
// A HashSet will automatically grow or shrink its capacity as items are added
// or removed.
//
// T must implement HashFunc[H], where H is of Hash type. This allows custom types
// that include non-comparable fields to provide their own hash algorithm.
func NewHashSet[T HashFunc[H], H Hash](size int) *HashSet[T, H] {
	return &HashSet[T, H]{
		items: make(map[H]T, max(0, size)),
	}
}

// HashSetFrom creates a new HashSet containing each item in items.
//
// T must implement HashFunc[H], where H is of type Hash. This allows custom types
// that include non-comparable fields to provide their own hash algorithm.
func HashSetFrom[T HashFunc[H], H Hash](items []T) *HashSet[T, H] {
	s := NewHashSet[T, H](len(items))
	s.InsertSlice(items)
	return s
}

// Insert item into s.
//
// Return true if s was modified (item was not already in s), false otherwise.
func (s *HashSet[T, H]) Insert(item T) bool {
	key := item.Hash()
	if _, exists := s.items[key]; exists {
		return false
	}
	s.items[key] = item
	return true
}

// InsertAll will insert each item in items into s.
//
// Return true if s was modified (at least one item was not already in s), false otherwise.
//
// Deprecated: use InsertSlice instead.
func (s *HashSet[T, H]) InsertAll(items []T) bool {
	return s.InsertSlice(items)
}

// InsertSlice will insert each item in items into s.
//
// Return true if s was modified (at least one item was not already in s), false otherwise.
func (s *HashSet[T, H]) InsertSlice(items []T) bool {
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
func (s *HashSet[T, H]) InsertSet(o *HashSet[T, H]) bool {
	modified := false
	for key, value := range o.items {
		if _, exists := s.items[key]; !exists {
			modified = true
		}
		s.items[key] = value
	}
	return modified
}

// Remove will remove item from s.
//
// Return true if s was modified (item was present), false otherwise.
func (s *HashSet[T, H]) Remove(item T) bool {
	key := item.Hash()
	if _, exists := s.items[key]; !exists {
		return false
	}
	delete(s.items, key)
	return true
}

// RemoveAll will remove each item in items from s.
//
// Return true if s was modified (any item was present), false otherwise.
//
// Deprecated: use RemoveSlice instead.
func (s *HashSet[T, H]) RemoveAll(items []T) bool {
	return s.RemoveSlice(items)
}

// RemoveSlice will remove each item in items from s.
//
// Return true if s was modified (any item was present), false otherwise.
func (s *HashSet[T, H]) RemoveSlice(items []T) bool {
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
func (s *HashSet[T, H]) RemoveSet(o *HashSet[T, H]) bool {
	modified := false
	for key := range o.items {
		if _, exists := s.items[key]; exists {
			modified = true
			delete(s.items, key)
		}
	}
	return modified
}

// RemoveFunc will remove each element from s that satisfies condition f.
//
// Return true if s was modified, false otherwise.
func (s *HashSet[T, H]) RemoveFunc(f func(item T) bool) bool {
	modified := false
	for _, item := range s.items {
		if applies := f(item); applies {
			s.Remove(item)
			modified = true
		}
	}
	return modified
}

// Contains returns whether item is present in s.
func (s *HashSet[T, H]) Contains(item T) bool {
	_, exists := s.items[item.Hash()]
	return exists
}

// ContainsAll returns whether s contains at least every item in items.
func (s *HashSet[T, H]) ContainsAll(items []T) bool {
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
func (s *HashSet[T, H]) ContainsSlice(items []T) bool {
	return s.Equal(HashSetFrom[T, H](items))
}

// Subset returns whether o is a subset of s.
func (s *HashSet[T, H]) Subset(o *HashSet[T, H]) bool {
	if len(s.items) < len(o.items) {
		return false
	}
	for _, item := range o.items {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// Size returns the cardinality of s.
func (s *HashSet[T, H]) Size() int {
	return len(s.items)
}

// Empty returns true if s contains no elements, false otherwise.
func (s *HashSet[T, H]) Empty() bool {
	return s.Size() == 0
}

// Union returns a set that contains all elements of s and o combined.
func (s *HashSet[T, H]) Union(o *HashSet[T, H]) *HashSet[T, H] {
	result := NewHashSet[T, H](s.Size())
	for key, item := range s.items {
		result.items[key] = item
	}
	for key, item := range o.items {
		result.items[key] = item
	}
	return result
}

// Difference returns a set that contains elements of s that are not in o.
func (s *HashSet[T, H]) Difference(o *HashSet[T, H]) *HashSet[T, H] {
	result := NewHashSet[T, H](max(0, s.Size()-o.Size()))
	for key, item := range s.items {
		if _, exists := o.items[key]; !exists {
			result.items[key] = item
		}
	}
	return result
}

// Intersect returns a set that contains elements that are present in both s and o.
func (s *HashSet[T, H]) Intersect(o *HashSet[T, H]) *HashSet[T, H] {
	result := NewHashSet[T, H](0)
	big, small := s, o
	if s.Size() < o.Size() {
		big, small = o, s
	}
	for _, item := range small.items {
		if big.Contains(item) {
			result.Insert(item)
		}
	}
	return result
}

// Copy creates a shallow copy of s.
func (s *HashSet[T, H]) Copy() *HashSet[T, H] {
	result := NewHashSet[T, H](s.Size())
	for key, item := range s.items {
		result.items[key] = item
	}
	return result
}

// Slice creates a copy of s as a slice.
//
// The result is not ordered.
func (s *HashSet[T, H]) Slice() []T {
	result := make([]T, 0, s.Size())
	for _, item := range s.items {
		result = append(result, item)
	}
	return result
}

// List creates a copy of s as a slice.
//
// Deprecated: use Slice() instead.
func (s *HashSet[T, H]) List() []T {
	return s.Slice()
}

// String creates a string representation of s, using "%v" printf formatting to transform
// each element into a string. The result contains elements sorted by their lexical
// string order.
func (s *HashSet[T, H]) String() string {
	return s.StringFunc(func(element T) string {
		return fmt.Sprintf("%v", element)
	})
}

// StringFunc creates a string representation of s, using f to transform each element
// into a string. The result contains elements sorted by their string order.
func (s *HashSet[T, H]) StringFunc(f func(element T) string) string {
	l := make([]string, 0, s.Size())
	for _, item := range s.items {
		l = append(l, f(item))
	}
	sort.Strings(l)
	return fmt.Sprintf("%s", l)
}

// Equal returns whether s and o contain the same elements.
func (s *HashSet[T, H]) Equal(o *HashSet[T, H]) bool {
	if len(s.items) != len(o.items) {
		return false
	}
	for _, item := range s.items {
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
func (s *HashSet[T, H]) EqualSlice(items []T) bool {
	if len(s.items) != len(items) {
		return false
	}
	return s.ContainsAll(items)
}

// MarshalJSON implements the json.Marshaler interface.
func (s *HashSet[T, H]) MarshalJSON() ([]byte, error) {
	return marshalJSON[T](s)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *HashSet[T, H]) UnmarshalJSON(data []byte) error {
	return unmarshalJSON[T](s, data)
}

func (s *HashSet[T, H]) ForEach(visit func(T) bool) {
	for _, item := range s.items {
		if !visit(item) {
			return
		}
	}
}
