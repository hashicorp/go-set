package set

import (
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
		InsertSliceFunc[string, int](s, numbers, transform)
		slices := s.Slice()
		sort.Strings(slices)
		must.SliceEqFunc(t, slices, []string{"1", "2", "3"}, func(a, b string) bool { return a == b })
	})
	t.Run("hashset", func(t *testing.T) {
		s := NewHashSet[*company, string](10)
		InsertSliceFunc[*company, int](s, numbers, func(element int) *company {
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
		s := NewTreeSet[string, Compare[string]](Cmp[string])
		InsertSliceFunc[string, int](s, numbers, func(element int) string {
			return strconv.Itoa(element)
		})
		invariants(t, s, Cmp[string])
		must.SliceEqFunc(t, s.Slice(), []string{"1", "2", "3"}, func(a, b string) bool { return a == b })
	})
}

func TestTransformSlice(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		s := From(ints(3))

		slice := TransformSlice[int, string](s, func(element int) string {
			return strconv.Itoa(element)
		})
		sort.Strings(slice)
		must.SliceEqFunc(t, slice, []string{"1", "2", "3"}, func(a, b string) bool { return a == b })
	})

	t.Run("hashset", func(t *testing.T) {
		s := NewHashSet[*company, string](10)
		s.InsertSlice([]*company{c1, c2, c3})
		slice := TransformSlice[*company, string](s, func(element *company) string {
			return element.Hash()
		})
		sort.Strings(slice)
		must.SliceEqFunc(t, slice, []string{"street:1", "street:2", "street:3"}, func(a, b string) bool { return a == b })
	})
	t.Run("treeSet", func(t *testing.T) {
		s := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3}, Cmp[int])
		slice := TransformSlice[int, string](s, func(element int) string {
			return strconv.Itoa(element)
		})
		sort.Strings(slice)
		must.SliceEqFunc(t, slice, []string{"1", "2", "3"}, func(a, b string) bool { return a == b })
	})
}

func TestTransform(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		a := From(ints(3))
		t.Run("set -> set", func(t *testing.T) {
			b := New[string](3)
			TransformUnion[int, string](a, b, func(element int) string {
				return strconv.Itoa(element)
			})
			slice := b.Slice()
			sort.Strings(slice)
			must.SliceEqFunc(t, slice, []string{"1", "2", "3"}, func(a, b string) bool { return a == b })
		})

		t.Run("set -> hashset", func(t *testing.T) {
			b := NewHashSet[*company, string](10)
			TransformUnion[int, *company](a, b, func(element int) *company {
				return &company{
					address: "street",
					floor:   element,
				}
			})
			must.MapContainsKeys(t, b.items, []string{
				"street:1", "street:2", "street:3",
			})
		})
		t.Run("set -> treeSet", func(t *testing.T) {
			b := NewTreeSet[string, Compare[string]](Cmp[string])
			TransformUnion[int, string](a, b, func(element int) string {
				return strconv.Itoa(element)
			})
			slice := b.Slice()
			sort.Strings(slice)
			must.SliceEqFunc(t, slice, []string{"1", "2", "3"}, func(a, b string) bool { return a == b })
		})
	})
	t.Run("hashSet", func(t *testing.T) {
		a := NewHashSet[*company, string](10)
		a.InsertSlice([]*company{c1, c2, c3})
		t.Run("hashSet -> set", func(t *testing.T) {
			b := New[int](3)
			TransformUnion[*company, int](a, b, func(element *company) int {
				return element.floor
			})
			slice := b.Slice()
			sort.Ints(slice)
			must.SliceEqFunc(t, slice, []int{1, 2, 3}, func(a, b int) bool { return a == b })
		})
		t.Run("hashSet -> hashSet", func(t *testing.T) {
			b := NewHashSet[*company, string](10)
			TransformUnion[*company, *company](a, b, func(element *company) *company {
				return &company{
					address: element.address,
					floor:   element.floor * 5,
				}
			})
			must.MapContainsKeys(t, b.items, []string{
				"street:5", "street:10", "street:15",
			})
		})
		t.Run("hashSet -> treeSet", func(t *testing.T) {
			b := NewTreeSet[int, Compare[int]](Cmp[int])
			TransformUnion[*company, int](a, b, func(element *company) int {
				return element.floor
			})
			slice := b.Slice()
			sort.Ints(slice)
			must.SliceEqFunc(t, slice, []int{1, 2, 3}, func(a, b int) bool { return a == b })
		})
	})
	t.Run("treeSet", func(t *testing.T) {
		a := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3}, Cmp[int])
		t.Run("treeSet -> set", func(t *testing.T) {
			b := New[string](3)
			TransformUnion[int, string](a, b, func(element int) string {
				return strconv.Itoa(element)
			})
			slice := b.Slice()
			sort.Strings(slice)
			must.SliceEqFunc(t, slice, []string{"1", "2", "3"}, func(a, b string) bool { return a == b })
		})
		t.Run("treeSet -> hashSet", func(t *testing.T) {
			b := NewHashSet[*company, string](10)
			TransformUnion[int, *company](a, b, func(element int) *company {
				return &company{
					address: "street",
					floor:   element,
				}
			})
			must.MapContainsKeys(t, b.items, []string{
				"street:1", "street:2", "street:3",
			})
		})
		t.Run("treeSet -> treeSet", func(t *testing.T) {
			b := NewTreeSet[string, Compare[string]](Cmp[string])
			TransformUnion[int, string](a, b, func(element int) string {
				return strconv.Itoa(element)
			})
			slice := b.Slice()
			sort.Strings(slice)
			must.SliceEqFunc(t, slice, []string{"1", "2", "3"}, func(a, b string) bool { return a == b })
		})
	})
}
