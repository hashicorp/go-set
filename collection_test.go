// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package set

import (
	"cmp"
	"sort"
	"strconv"
	"testing"

	"github.com/shoenig/test/must"
)

func TestInsertSliceFunc(t *testing.T) {
	numbers := ints(3)

	t.Run("set", func(t *testing.T) {
		s := New[string](10)
		transform := func(element int) string { return strconv.Itoa(element) }
		InsertSliceFunc[string](s, numbers, transform)
		slices := s.Slice()
		sort.Strings(slices)
		must.SliceEqFunc(t, slices, []string{"1", "2", "3"}, func(a, b string) bool { return a == b })
	})

	t.Run("hashset", func(t *testing.T) {
		s := NewHashSet[*company, string](10)
		InsertSliceFunc[*company](s, numbers, func(element int) *company {
			return &company{
				address: "InsertSliceFunc",
				floor:   element,
			}
		})
		must.MapContainsKeys(t, s.items, []string{
			"InsertSliceFunc:1", "InsertSliceFunc:2", "InsertSliceFunc:3",
		})
	})

	t.Run("treeSet", func(t *testing.T) {
		s := NewTreeSet[string](cmp.Compare[string])
		InsertSliceFunc[string](s, numbers, func(element int) string {
			return strconv.Itoa(element)
		})
		invariants(t, s, cmp.Compare[string])
		must.SliceEqFunc(t, s.Slice(), []string{"1", "2", "3"}, func(a, b string) bool { return a == b })
	})
}

func TestSliceFunc(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		s := From(ints(3))
		slice := SliceFunc[int](s, func(element int) string {
			return strconv.Itoa(element)
		})
		sort.Strings(slice)
		must.SliceEqFunc(t, slice, []string{"1", "2", "3"}, func(a, b string) bool { return a == b })
	})

	t.Run("hashset", func(t *testing.T) {
		s := NewHashSet[*company, string](10)
		s.InsertSlice([]*company{c1, c2, c3})
		slice := SliceFunc[*company](s, func(element *company) string {
			return element.Hash()
		})
		sort.Strings(slice)
		must.SliceEqFunc(t, slice, []string{"street:1", "street:2", "street:3"}, func(a, b string) bool { return a == b })
	})

	t.Run("treeSet", func(t *testing.T) {
		s := TreeSetFrom[int]([]int{1, 2, 3}, cmp.Compare[int])
		slice := SliceFunc[int](s, func(element int) string {
			return strconv.Itoa(element)
		})
		sort.Strings(slice)
		must.SliceEqFunc(t, slice, []string{"1", "2", "3"}, func(a, b string) bool { return a == b })
	})
}

func TestInsertSetFunc(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		a := From(ints(3))
		t.Run("set -> set", func(t *testing.T) {
			b := New[string](3)
			modified := InsertSetFunc[int, string](a, b, func(element int) string {
				return strconv.Itoa(element)
			})
			must.True(t, modified)
			slice := b.Slice()
			sort.Strings(slice)
			must.SliceEqFunc(t, slice, []string{"1", "2", "3"}, func(a, b string) bool { return a == b })
		})

		t.Run("set -> hashset", func(t *testing.T) {
			b := NewHashSet[*company, string](10)
			modified := InsertSetFunc[int, *company](a, b, func(element int) *company {
				return &company{
					address: "street",
					floor:   element,
				}
			})
			must.True(t, modified)
			must.MapContainsKeys(t, b.items, []string{
				"street:1", "street:2", "street:3",
			})
		})

		t.Run("set -> treeSet", func(t *testing.T) {
			b := NewTreeSet[string](cmp.Compare[string])
			modified := InsertSetFunc[int, string](a, b, func(element int) string {
				return strconv.Itoa(element)
			})
			must.True(t, modified)
			slice := b.Slice()
			sort.Strings(slice)
			must.SliceEqFunc(t, slice, []string{"1", "2", "3"}, func(a, b string) bool { return a == b })
		})

		t.Run("not modified", func(t *testing.T) {
			b := a.Copy()
			modified := InsertSetFunc[int, int](a, b, func(element int) int {
				return element
			})
			must.False(t, modified)
		})
	})

	t.Run("hashSet", func(t *testing.T) {
		a := NewHashSet[*company, string](10)
		a.InsertSlice([]*company{c1, c2, c3})

		t.Run("hashSet -> set", func(t *testing.T) {
			b := New[int](3)
			modified := InsertSetFunc[*company, int](a, b, func(element *company) int {
				return element.floor
			})
			must.True(t, modified)
			slice := b.Slice()
			sort.Ints(slice)
			must.SliceEqFunc(t, slice, []int{1, 2, 3}, func(a, b int) bool { return a == b })
		})

		t.Run("hashSet -> hashSet", func(t *testing.T) {
			b := NewHashSet[*company, string](10)
			modified := InsertSetFunc[*company, *company](a, b, func(element *company) *company {
				return &company{
					address: element.address,
					floor:   element.floor * 5,
				}
			})
			must.True(t, modified)
			must.MapContainsKeys(t, b.items, []string{
				"street:5", "street:10", "street:15",
			})
		})

		t.Run("hashSet -> treeSet", func(t *testing.T) {
			b := NewTreeSet[int](cmp.Compare[int])
			modified := InsertSetFunc[*company, int](a, b, func(element *company) int {
				return element.floor
			})
			must.True(t, modified)
			slice := b.Slice()
			sort.Ints(slice)
			must.SliceEqFunc(t, slice, []int{1, 2, 3}, func(a, b int) bool { return a == b })
		})

		t.Run("not modified", func(t *testing.T) {
			b := a.Copy()
			modified := InsertSetFunc[*company, *company](a, b, func(element *company) *company {
				return element
			})
			must.False(t, modified)
		})
	})

	t.Run("treeSet", func(t *testing.T) {
		a := TreeSetFrom[int]([]int{1, 2, 3}, cmp.Compare[int])

		t.Run("treeSet -> set", func(t *testing.T) {
			b := New[string](3)
			modified := InsertSetFunc[int, string](a, b, func(element int) string {
				return strconv.Itoa(element)
			})
			must.True(t, modified)
			slice := b.Slice()
			sort.Strings(slice)
			must.SliceEqFunc(t, slice, []string{"1", "2", "3"}, func(a, b string) bool { return a == b })
		})

		t.Run("treeSet -> hashSet", func(t *testing.T) {
			b := NewHashSet[*company, string](10)
			modified := InsertSetFunc[int, *company](a, b, func(element int) *company {
				return &company{
					address: "street",
					floor:   element,
				}
			})
			must.True(t, modified)
			must.MapContainsKeys(t, b.items, []string{
				"street:1", "street:2", "street:3",
			})
		})

		t.Run("treeSet -> treeSet", func(t *testing.T) {
			b := NewTreeSet[string](cmp.Compare[string])
			modified := InsertSetFunc[int, string](a, b, func(element int) string {
				return strconv.Itoa(element)
			})
			must.True(t, modified)
			slice := b.Slice()
			sort.Strings(slice)
			must.SliceEqFunc(t, slice, []string{"1", "2", "3"}, func(a, b string) bool { return a == b })
		})

		t.Run("not modified", func(t *testing.T) {
			b := a.Copy()
			modified := InsertSetFunc[int, int](a, b, func(element int) int {
				return element
			})
			must.False(t, modified)
		})
	})
}

func TestEqualSet(t *testing.T) {
	t.Run("equal ok", func(t *testing.T) {
		a := From(ints(3))
		b := From(ints(3))
		must.True(t, a.EqualSet(b))
	})
	t.Run("none equal none", func(t *testing.T) {
		a := New[int](0)
		b := New[int](0)
		must.True(t, a.EqualSet(b))
	})
	t.Run("size not equal", func(t *testing.T) {
		a := From(ints(3))
		b := From(ints(4))
		must.False(t, a.EqualSet(b))
	})
	t.Run("items not equal", func(t *testing.T) {
		a := From(ints(3))
		b := From([]int{1, 2, 4})
		must.False(t, a.EqualSet(b))
	})
}
