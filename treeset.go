// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package set

import (
	"fmt"
)

// Compare represents a function that compares two elements.
//
// Must return
// < 0 if the first parameter is less than the second parameter
// 0 if the two parameters are equal
// > 0 if the first parameters is greater than the second parameter
type Compare[T any] func(T, T) int

// BuiltIn types compatible with Cmp
type BuiltIn interface {
	~string | ~int | ~uint | ~int64 | ~uint64 | ~int32 | ~uint32 | ~int16 | ~uint16 | ~int8 | ~uint8
}

// Cmp is a Compare function for the specified builtin type B.
//
// Common to use with string, int, etc.
func Cmp[B BuiltIn](x, y B) int {
	switch {
	case x < y:
		return -1
	case x > y:
		return 1
	default:
		return 0
	}
}

// TreeSet provides a generic sortable set implementation for Go.
// Enables fast storage and retrieval of ordered information. Most effective
// in cases where data is regularly being added and/or removed and fast
// lookup properties must be maintained.
//
// The underlying data structure is a Red-Black Binary Search Tree.
// https://en.wikipedia.org/wiki/Red–black_tree
//
// Not thread safe, and not safe for concurrent modification.
type TreeSet[T any, C Compare[T]] struct {
	comparison C
	root       *node[T]
	marker     *node[T]
	size       int
}

// NewTreeSet creates a TreeSet of type T, comparing elements via C.
//
// T may be any type.
//
// C is an implementation of Compare[T]. For builtin types, Cmp provides
// a convenient Compare implementation.
func NewTreeSet[T any, C Compare[T]](compare C) *TreeSet[T, C] {
	return &TreeSet[T, C]{
		comparison: compare,
		root:       nil,
		marker:     &node[T]{color: black},
		size:       0,
	}
}

// TreeSetFrom creates a new TreeSet containing each item in items.
//
// T may be any type.
//
// C is an implementation of Compare[T]. For builtin types, Cmp provides a
// convenient Compare implementation.
func TreeSetFrom[T any, C Compare[T]](items []T, compare C) *TreeSet[T, C] {
	s := NewTreeSet[T](compare)
	s.InsertSlice(items)
	return s
}

// Insert item into s.
//
// Returns true if s was modified (item was not already in s), false otherwise.
func (s *TreeSet[T, C]) Insert(item T) bool {
	return s.insert(&node[T]{
		element: item,
		color:   red,
	})
}

// InsertSlice will insert each item in items into s.
//
// Return true if s was modified (at least one item was not already in s), false otherwise.
func (s *TreeSet[T, C]) InsertSlice(items []T) bool {
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
func (s *TreeSet[T, C]) InsertSet(o *TreeSet[T, C]) bool {
	modified := false
	insert := func(item T) bool {
		if s.Insert(item) {
			modified = true
		}
		return true
	}
	o.ForEach(insert)
	return modified
}

// Remove item from s.
//
// Returns true if s was modified (item was in s), false otherwise.
func (s *TreeSet[T, C]) Remove(item T) bool {
	return s.delete(item)
}

// RemoveSlice will remove each item in items from s.
//
// Return true if s was modified (any item was in s), false otherwise.
func (s *TreeSet[T, C]) RemoveSlice(items []T) bool {
	modified := false
	for _, item := range items {
		if s.Remove(item) {
			modified = true
		}
	}
	return modified
}

// RemoveSet will remove each element in o from s.
//
// Returns true if s was modified (at least one item in o was in s), false otherwise.
func (s *TreeSet[T, C]) RemoveSet(o *TreeSet[T, C]) bool {
	modified := false
	remove := func(item T) bool {
		if s.Remove(item) {
			modified = true
		}
		return true
	}
	o.ForEach(remove)
	return modified
}

// RemoveFunc will remove each element from s that satisifies condition f.
//
// Return true if s was modified, false otherwise.
func (s *TreeSet[T, C]) RemoveFunc(f func(T) bool) bool {
	removeIds := s.FilterSlice(f)
	return s.RemoveSlice(removeIds)
}

// Min returns the smallest item in the set.
//
// Must not be called on an empty set.
func (s *TreeSet[T, C]) Min() T {
	if s.root == nil {
		panic("min: tree is empty")
	}
	n := s.min(s.root)
	return n.element
}

// Max returns the largest item in s.
//
// Must not be called on an empty set.
func (s *TreeSet[T, C]) Max() T {
	if s.root == nil {
		panic("max: tree is empty")
	}
	n := s.max(s.root)
	return n.element
}

// TopK returns the top n (smallest) elements in s, in ascending order.
func (s *TreeSet[T, C]) TopK(n int) []T {
	result := make([]T, 0, n)
	s.fillLeft(s.root, &result)
	return result
}

// BottomK returns the bottom n (largest) elements in s, in descending order.
func (s *TreeSet[T, C]) BottomK(n int) []T {
	result := make([]T, 0, n)
	s.fillRight(s.root, &result)
	return result
}

// FirstBelow returns the first element strictly below item.
//
// A zero value and false are returned if no such element exists.
func (s *TreeSet[T, C]) FirstBelow(item T) (T, bool) {
	var candidate *node[T] = nil
	var n = s.root
	for n != nil {
		c := s.comparison(item, n.element)
		switch {
		case c > 0:
			candidate = n
			n = n.right
		case c <= 0:
			n = n.left
		}
	}
	return candidate.get()
}

// FirstBelowEqual returns the first element below item (or item itself if present).
//
// A zero value and false are returned if no such element exists.
func (s *TreeSet[T, C]) FirstBelowEqual(item T) (T, bool) {
	var candidate *node[T] = nil
	var n = s.root
	for n != nil {
		c := s.comparison(item, n.element)
		switch {
		case c == 0:
			return n.get()
		case c > 0:
			candidate = n
			n = n.right
		case c < 0:
			n = n.left
		}
	}
	return candidate.get()
}

// Below returns a TreeSet containing the elements of s that are < item.
func (s *TreeSet[T, C]) Below(item T) *TreeSet[T, C] {
	result := NewTreeSet[T](s.comparison)
	s.filterLeft(s.root, func(element T) bool {
		return s.comparison(element, item) < 0
	}, result)
	return result
}

// BelowEqual returns a TreeSet containing the elements of s that are ≤ item.
func (s *TreeSet[T, C]) BelowEqual(item T) *TreeSet[T, C] {
	result := NewTreeSet[T](s.comparison)
	s.filterLeft(s.root, func(element T) bool {
		return s.comparison(element, item) <= 0
	}, result)
	return result
}

// FirstAbove returns the first element strictly above item.
//
// A zero value and false are returned if no such element exists.
func (s *TreeSet[T, C]) FirstAbove(item T) (T, bool) {
	var candidate *node[T] = nil
	var n = s.root
	for n != nil {
		c := s.comparison(item, n.element)
		switch {
		case c < 0:
			candidate = n
			n = n.left
		case c >= 0:
			n = n.right
		}
	}
	return candidate.get()
}

// FirstAboveEqual returns the first element above item (or item itself if present).
//
// A zero value and false are returned if no such element exists.
func (s *TreeSet[T, C]) FirstAboveEqual(item T) (T, bool) {
	var candidate *node[T]
	var n = s.root
	for n != nil {
		c := s.comparison(item, n.element)
		switch {
		case c == 0:
			return n.get()
		case c < 0:
			candidate = n
			n = n.left
		case c > 0:
			n = n.right
		}
	}
	return candidate.get()
}

// After returns a TreeSet containing the elements of s that are > item.
func (s *TreeSet[T, C]) Above(item T) *TreeSet[T, C] {
	result := NewTreeSet[T](s.comparison)
	s.filterRight(s.root, func(element T) bool {
		return s.comparison(element, item) > 0
	}, result)
	return result
}

// AfterEqual returns a TreeSet containing the elements of s that are ≥ item.
func (s *TreeSet[T, C]) AboveEqual(item T) *TreeSet[T, C] {
	result := NewTreeSet[T](s.comparison)
	s.filterRight(s.root, func(element T) bool {
		return s.comparison(element, item) >= 0
	}, result)
	return result
}

// Contains returns whether item is present in s.
func (s *TreeSet[T, C]) Contains(item T) bool {
	return s.locate(s.root, item) != nil
}

// ContainsSlice returns whether s contains the same set of elements that are in
// items. The items slice may contain duplicate elements.
//
// If the items slice is known to be set-like (no duplicates), EqualSlice provides
// a more efficient implementation.
func (s *TreeSet[T, C]) ContainsSlice(items []T) bool {
	for _, item := range items {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// Size returns the number of elements in s.
func (s *TreeSet[T, C]) Size() int {
	return s.size
}

// Empty returns true if there are no elements in s.
func (s *TreeSet[T, C]) Empty() bool {
	return s.Size() == 0
}

// Slice returns the elements of s as a slice, in order.
func (s *TreeSet[T, C]) Slice() []T {
	result := make([]T, 0, s.Size())
	s.ForEach(func(element T) bool {
		result = append(result, element)
		return true
	})
	return result
}

// FilterSlice returns the elements of s that satisfy the predicate f.
func (s *TreeSet[T, C]) FilterSlice(filter func(T) bool) []T {
	result := make([]T, 0, s.Size())
	s.ForEach(func(t T) bool {
		if filter != nil && filter(t) {
			result = append(result, t)
		}
		return true
	})
	return result
}

// Subset returns whether o is a subset of s.
func (s *TreeSet[T, C]) Subset(o *TreeSet[T, C]) bool {
	// try the fast paths
	if o.Empty() {
		return true
	}
	if s.Empty() {
		return false
	}
	if s.Size() < o.Size() {
		return false
	}

	// iterate o, and increment s finding each element
	// i.e. merge algorithm but with channels
	iterO := o.iterate()
	iterS := s.iterate()

	idxO := 0
	idxS := 0

next:
	for ; idxO < o.Size(); idxO++ {
		nextO := iterO()
		for idxS < s.Size() {
			idxS++
			nextS := iterS()
			cmp := s.compare(nextS, nextO)
			switch {
			case cmp > 0:
				return false
			case cmp < 0:
				continue
			default:
				continue next
			}
		}
		return false
	}
	return true
}

// Union returns a set that contains all elements of s and o combined.
func (s *TreeSet[T, C]) Union(o *TreeSet[T, C]) *TreeSet[T, C] {
	tree := NewTreeSet[T](s.comparison)
	f := func(n *node[T]) { tree.Insert(n.element) }
	s.prefix(f, s.root)
	o.prefix(f, o.root)
	return tree
}

// Difference returns a set that contains elements of s that are not in o.
func (s *TreeSet[T, C]) Difference(o *TreeSet[T, C]) *TreeSet[T, C] {
	tree := NewTreeSet[T](s.comparison)
	f := func(n *node[T]) {
		if !o.Contains(n.element) {
			tree.Insert(n.element)
		}
	}
	s.prefix(f, s.root)
	return tree
}

// Intersect returns a set that contains elements that are present in both s and o.
func (s *TreeSet[T, C]) Intersect(o *TreeSet[T, C]) *TreeSet[T, C] {
	tree := NewTreeSet[T](s.comparison)
	f := func(n *node[T]) {
		if o.Contains(n.element) {
			tree.Insert(n.element)
		}
	}
	s.prefix(f, s.root)
	return tree
}

// Copy creates a copy of s.
//
// Individual elements are reference copies.
func (s *TreeSet[T, C]) Copy() *TreeSet[T, C] {
	tree := NewTreeSet[T](s.comparison)
	f := func(n *node[T]) {
		tree.Insert(n.element)
	}
	s.prefix(f, s.root)
	return tree
}

// Equal return whether s and o contain the same elements.
func (s *TreeSet[T, C]) Equal(o *TreeSet[T, C]) bool {
	// try the fast fail paths
	if s.Empty() || o.Empty() {
		return s.Size() == o.Size()
	}
	switch {
	case s.Size() != o.Size():
		return false
	case s.comparison(s.Min(), o.Min()) != 0:
		return false
	case s.comparison(s.Max(), o.Max()) != 0:
		return false
	}

	iterS := s.iterate()
	iterO := o.iterate()
	for i := 0; i < s.Size(); i++ {
		nextS := iterS()
		nextO := iterO()
		if s.compare(nextS, nextO) != 0 {
			return false
		}
	}

	return true
}

// EqualSlice returns whether s and items contain the same elements.
func (s *TreeSet[T, C]) EqualSlice(items []T) bool {
	if s.Size() != len(items) {
		return false
	}
	return s.ContainsSlice(items)
}

// String creates a string representation of s, using "%v" printf formatting
// each element into a string. The result contains elements in order.
func (s *TreeSet[T, C]) String() string {
	return s.StringFunc(func(element T) string {
		return fmt.Sprintf("%v", element)
	})
}

// StringFunc creates a string representation of s, using f to transform each
// element into a string. The result contains elements in order.
func (s *TreeSet[T, C]) StringFunc(f func(element T) string) string {
	l := make([]string, 0, s.Size())
	s.ForEach(func(element T) bool {
		l = append(l, f(element))
		return true
	})
	return fmt.Sprintf("%s", l)
}

func (s *TreeSet[T, C]) ForEach(visit func(T) bool) {
	s.infix(func(n *node[T]) (next bool) {
		return visit(n.element)
	}, s.root)
}

// Red-Black Tree Invariants
//
// 1. each node is either red or black
// 2. the root node is always black
// 3. nil leaf nodes are always black
// 4. a red node must not have red children
// 5. all simple paths from a node to nil leaf contain the same number of
// black nodes

type color bool

const (
	red   color = false
	black color = true
)

type node[T any] struct {
	element T
	color   color
	parent  *node[T]
	left    *node[T]
	right   *node[T]
}

func (n *node[T]) black() bool {
	return n == nil || n.color == black
}

func (n *node[T]) red() bool {
	return n != nil && n.color == red
}

func (n *node[T]) get() (T, bool) {
	if n == nil {
		var zero T
		return zero, false
	}
	return n.element, true
}

func (s *TreeSet[T, C]) locate(start *node[T], target T) *node[T] {
	n := start
	for {
		if n == nil {
			return nil
		}
		cmp := s.compare(n, &node[T]{element: target})
		switch {
		case cmp < 0:
			n = n.right
		case cmp > 0:
			n = n.left
		default:
			return n
		}
	}
}

func (s *TreeSet[T, C]) rotateRight(n *node[T]) {
	parent := n.parent
	leftChild := n.left

	n.left = leftChild.right
	if leftChild.right != nil {
		leftChild.right.parent = n
	}

	leftChild.right = n
	n.parent = leftChild

	s.replaceChild(parent, n, leftChild)
}

func (s *TreeSet[T, C]) rotateLeft(n *node[T]) {
	parent := n.parent
	rightChild := n.right

	n.right = rightChild.left
	if rightChild.left != nil {
		rightChild.left.parent = n
	}

	rightChild.left = n
	n.parent = rightChild

	s.replaceChild(parent, n, rightChild)
}

func (s *TreeSet[T, C]) replaceChild(parent, previous, next *node[T]) {
	switch {
	case parent == nil:
		s.root = next
	case parent.left == previous:
		parent.left = next
	case parent.right == previous:
		parent.right = next
	default:
		panic("node is not child of its parent")
	}

	if next != nil {
		next.parent = parent
	}
}

func (s *TreeSet[T, C]) insert(n *node[T]) bool {
	var (
		parent *node[T] = nil
		tmp    *node[T] = s.root
	)

	for tmp != nil {
		parent = tmp

		cmp := s.compare(n, tmp)
		switch {
		case cmp < 0:
			tmp = tmp.left
		case cmp > 0:
			tmp = tmp.right
		default:
			// already exists in tree
			return false
		}
	}

	n.color = red
	switch {
	case parent == nil:
		s.root = n
	case s.compare(n, parent) < 0:
		parent.left = n
	default:
		parent.right = n
	}
	n.parent = parent

	s.rebalanceInsertion(n)
	s.size++
	return true
}

func (s *TreeSet[T, C]) rebalanceInsertion(n *node[T]) {
	parent := n.parent

	// case 1: parent is nil
	// - means we are the root
	// - our color must be black
	if parent == nil {
		n.color = black
		return
	}

	// if parent is black there is nothing to do
	if parent.black() {
		return
	}

	// case 2: no grandparent
	// - implies the parent is root
	// - we must now be black
	grandparent := parent.parent
	if grandparent == nil {
		parent.color = black
		return
	}

	uncle := s.uncleOf(parent)

	switch {
	// case 3: uncle is red
	// - fix color of parent, grandparent, uncle
	// - recurse upwards as necessary
	case uncle != nil && uncle.red():
		parent.color = black
		grandparent.color = red
		uncle.color = black
		s.rebalanceInsertion(grandparent)

	case parent == grandparent.left:
		// case 4a: uncle is black
		// + node is left->right child of its grandparent
		if n == parent.right {
			s.rotateLeft(parent)
			parent = n // recolor in case 5a
		}

		// case 5a: uncle is black
		// + node is left->left child of its grandparent
		s.rotateRight(grandparent)

		// fix color of original parent and grandparent
		parent.color = black
		grandparent.color = red

		// parent is right child of grandparent
	default:
		// case 4b: uncle is black
		// + node is right->left child of its grandparent
		if n == parent.left {
			s.rotateRight(parent)
			// points to root of rotated sub tree
			parent = n // recolor in case 5b
		}

		// case 5b: uncle is black
		// + node is right->right child of its grandparent
		s.rotateLeft(grandparent)

		// fix color of original parent and grandparent
		parent.color = black
		grandparent.color = red
	}
}

func (s *TreeSet[T, C]) delete(element T) bool {
	n := s.locate(s.root, element)
	if n == nil {
		return false
	}

	var (
		moved   *node[T]
		deleted color
	)

	if n.left == nil || n.right == nil {
		// case where deleted node had zero or one child
		moved = s.delete01(n)
		deleted = n.color
	} else {
		// case where node has two children

		// find minimum of right subtree
		successor := s.min(n.right)

		// copy successor data into n
		n.element = successor.element

		// delete successor
		moved = s.delete01(successor)
		deleted = successor.color
	}

	// re-balance if the node was black
	if deleted == black {
		s.rebalanceDeletion(moved)

		// remove marker
		if moved == s.marker {
			s.replaceChild(moved.parent, moved, nil)
		}
	}

	// element was removed
	s.size--
	s.marker.color = black
	s.marker.left = nil
	s.marker.right = nil
	s.marker.parent = nil
	return true
}

func (s *TreeSet[T, C]) delete01(n *node[T]) *node[T] {
	// node only has left child, replace by left child
	if n.left != nil {
		s.replaceChild(n.parent, n, n.left)
		return n.left
	}

	// node only has right child, replace by right child
	if n.right != nil {
		s.replaceChild(n.parent, n, n.right)
		return n.right
	}

	// node has both children
	// if node is black replace with marker
	// if node is red we just remove it
	if n.black() {
		s.replaceChild(n.parent, n, s.marker)
		return s.marker
	} else {
		s.replaceChild(n.parent, n, nil)
		return nil
	}
}

func (s *TreeSet[T, C]) rebalanceDeletion(n *node[T]) {
	// base case: node is root
	if n == s.root {
		n.color = black
		return
	}

	sibling := s.siblingOf(n)

	// case: sibling is red
	if sibling.red() {
		s.fixRedSibling(n, sibling)
		sibling = s.siblingOf(n)
	}

	// case: black sibling with two black children
	if sibling.left.black() && sibling.right.black() {
		sibling.color = red

		// case: black sibling with to black children and a red parent
		if n.parent.red() {
			n.parent.color = black
		} else {
			// case: black sibling with two black children and black parent
			s.rebalanceDeletion(n.parent)
		}
	} else {
		// case: black sibling with at least one red child
		s.fixBlackSibling(n, sibling)
	}
}

func (s *TreeSet[T, C]) fixRedSibling(n *node[T], sibling *node[T]) {
	sibling.color = black
	n.parent.color = red

	switch {
	case n == n.parent.left:
		s.rotateLeft(n.parent)
	default:
		s.rotateRight(n.parent)
	}
}

func (s *TreeSet[T, C]) fixBlackSibling(n, sibling *node[T]) {
	isLeftChild := n == n.parent.left

	if isLeftChild && sibling.right.black() {
		sibling.left.color = black
		sibling.color = red
		s.rotateRight(sibling)
		sibling = n.parent.right
	} else if !isLeftChild && sibling.left.black() {
		sibling.right.color = black
		sibling.color = red
		s.rotateLeft(sibling)
		sibling = n.parent.left
	}

	sibling.color = n.parent.color
	n.parent.color = black
	if isLeftChild {
		sibling.right.color = black
		s.rotateLeft(n.parent)
	} else {
		sibling.left.color = black
		s.rotateRight(n.parent)
	}
}

func (s *TreeSet[T, C]) siblingOf(n *node[T]) *node[T] {
	parent := n.parent
	switch {
	case n == parent.left:
		return parent.right
	case n == parent.right:
		return parent.left
	default:
		panic("bug: parent is not a child of its grandparent")
	}
}

func (*TreeSet[T, C]) uncleOf(n *node[T]) *node[T] {
	grandparent := n.parent
	switch {
	case grandparent.left == n:
		return grandparent.right
	case grandparent.right == n:
		return grandparent.left
	default:
		panic("bug: parent is not a child of our childs grandparent")
	}
}

func (s *TreeSet[T, C]) min(n *node[T]) *node[T] {
	for n.left != nil {
		n = n.left
	}
	return n
}

func (s *TreeSet[T, C]) max(n *node[T]) *node[T] {
	for n.right != nil {
		n = n.right
	}
	return n
}

func (s *TreeSet[T, C]) compare(a, b *node[T]) int {
	return s.comparison(a.element, b.element)
}

// TreeNodeVisit is a function that is called for each node in the tree.
type TreeNodeVisit[T any] func(*node[T]) (next bool)

func (s *TreeSet[T, C]) infix(visit TreeNodeVisit[T], n *node[T]) (next bool) {
	if n == nil {
		return true
	}
	if next = s.infix(visit, n.left); !next {
		return
	}
	if next = visit(n); !next {
		return
	}
	return s.infix(visit, n.right)
}

func (s *TreeSet[T, C]) fillLeft(n *node[T], k *[]T) {
	if n == nil {
		return
	}

	if len(*k) < cap(*k) {
		s.fillLeft(n.left, k)
	}

	if len(*k) < cap(*k) {
		*k = append(*k, n.element)
	}

	if len(*k) < cap(*k) {
		s.fillLeft(n.right, k)
	}
}

func (s *TreeSet[T, C]) fillRight(n *node[T], k *[]T) {
	if n == nil {
		return
	}

	if len(*k) < cap(*k) {
		s.fillRight(n.right, k)
	}

	if len(*k) < cap(*k) {
		*k = append(*k, n.element)
	}

	if len(*k) < cap(*k) {
		s.fillRight(n.left, k)
	}
}

func (s *TreeSet[T, C]) prefix(visit func(*node[T]), n *node[T]) {
	if n == nil {
		return
	}
	visit(n)
	s.prefix(visit, n.left)
	s.prefix(visit, n.right)
}

func (s *TreeSet[T, C]) iterate() func() *node[T] {
	stck := makeStack[*node[T]]()

	for n := s.root; n != nil; n = n.left {
		stck.push(n)
	}

	return func() *node[T] {
		if stck.empty() {
			return nil
		}
		n := stck.pop()
		for r := n.right; r != nil; r = r.left {
			stck.push(r)
		}
		return n
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (s *TreeSet[T, C]) MarshalJSON() ([]byte, error) {
	return marshalJSON[T](s)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *TreeSet[T, C]) UnmarshalJSON(data []byte) error {
	return unmarshalJSON[T](s, data)
}

func (s *TreeSet[T, C]) filterLeft(n *node[T], accept func(element T) bool, result *TreeSet[T, C]) {
	if n == nil {
		return
	}

	s.filterLeft(n.left, accept, result)

	if accept(n.element) {
		result.Insert(n.element)
		s.filterLeft(n.right, accept, result)
	}
}

func (s *TreeSet[T, C]) filterRight(n *node[T], accept func(element T) bool, result *TreeSet[T, C]) {
	if n == nil {
		return
	}

	s.filterRight(n.right, accept, result)

	if accept(n.element) {
		result.Insert(n.element)
		s.filterRight(n.left, accept, result)
	}
}
