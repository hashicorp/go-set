// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package set

import (
	"fmt"
	"testing"

	"github.com/shoenig/test/must"
)

type employee struct {
	name string
	id   int
}

func (e *employee) String() string {
	return fmt.Sprintf("(%d %s)", e.id, e.name)
}

func TestSet_New(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		s := New[float64](1)
		must.MapEmpty(t, s.items)
	})

	t.Run("zero", func(t *testing.T) {
		s := New[int](0)
		must.MapEmpty(t, s.items)
	})

	t.Run("negative", func(t *testing.T) {
		s := New[string](-1) // assume zero
		must.MapEmpty(t, s.items)
	})
}

func TestSet_From(t *testing.T) {
	t.Run("from nil", func(t *testing.T) {
		s := From[string](nil)
		must.MapEmpty(t, s.items)
	})

	t.Run("from some", func(t *testing.T) {
		s := From[string]([]string{"apple", "banana", "cherry"})
		must.MapContainsKeys(t, s.items, []string{"apple", "banana", "cherry"})
	})
}

func TestSet_FromFunc(t *testing.T) {
	employees := []employee{
		{"alice", 1}, {"bob", 2}, {"bob", 2}, {"carol", 3}, {"dave", 4},
	}
	s := FromFunc(employees, func(e employee) string {
		return e.name
	})
	must.MapContainsKeys(t, s.items, []string{"alice", "bob", "carol", "dave"})
}

func TestSet_Insert(t *testing.T) {
	t.Run("one int", func(t *testing.T) {
		s := New[int](10)
		must.True(t, s.Insert(1))
		must.MapContainsKeys(t, s.items, []int{1})
	})

	t.Run("one string", func(t *testing.T) {
		s := New[string](10)
		must.True(t, s.Insert("apple"))
		must.MapContainsKeys(t, s.items, []string{"apple"})
	})

	t.Run("re-insert", func(t *testing.T) {
		s := New[int](10)
		must.True(t, s.Insert(2))
		must.False(t, s.Insert(2))
		must.MapContainsKeys(t, s.items, []int{2})
	})

	t.Run("insert several", func(t *testing.T) {
		s := New[int](10)
		must.True(t, s.Insert(1))
		must.True(t, s.Insert(2))
		must.True(t, s.Insert(3))
		must.True(t, s.Insert(4))
		must.True(t, s.Insert(5))
		must.MapContainsKeys(t, s.items, []int{1, 2, 3, 4, 5})
	})

	t.Run("insert custom", func(t *testing.T) {
		s := New[employee](10)
		must.True(t, s.Insert(employee{"mitchell", 1}))
		must.True(t, s.Insert(employee{"armon", 2}))
		must.True(t, s.Insert(employee{"jack", 3}))
		must.False(t, s.Insert(employee{"jack", 3}))
		must.False(t, s.Insert(employee{"armon", 2}))
		must.False(t, s.Insert(employee{"mitchell", 1}))
		must.MapContainsKeys(t, s.items, []employee{
			{"mitchell", 1}, {"armon", 2}, {"jack", 3},
		})
	})
}

func TestSet_InsertSlice(t *testing.T) {
	t.Run("insert none", func(t *testing.T) {
		empty := New[int](0)
		must.False(t, empty.InsertSlice(nil))
		must.MapEmpty(t, empty.items)
	})

	t.Run("insert some", func(t *testing.T) {
		s := New[string](0)
		must.True(t, s.InsertSlice([]string{"apple", "banana", "cherry"}))
		must.MapContainsKeys(t, s.items, []string{"apple", "banana", "cherry"})
	})

	t.Run("insert duplicates", func(t *testing.T) {
		s := New[int](0)
		must.True(t, s.InsertSlice([]int{2, 4, 6, 8}))
		must.True(t, s.InsertSlice([]int{4, 5, 6}))
		must.MapContainsKeys(t, s.items, []int{2, 4, 5, 6, 8})
	})
}

func TestSet_InsertSet(t *testing.T) {
	t.Run("insert empty", func(t *testing.T) {
		a := From[int]([]int{1, 2, 3, 4})
		b := New[int](0)
		must.False(t, a.InsertSet(b))
		must.MapContainsKeys(t, a.items, []int{1, 2, 3, 4})
	})

	t.Run("insert some", func(t *testing.T) {
		a := From[int]([]int{1, 2, 3, 4})
		b := From[int]([]int{3, 4, 5, 6, 7})
		must.True(t, a.InsertSet(b))
		must.MapContainsKeys(t, a.items, []int{1, 2, 3, 4, 5, 6, 7})
	})
}

func TestSet_Contains(t *testing.T) {
	t.Run("contains string item", func(t *testing.T) {
		s := New[string](10)
		must.True(t, s.InsertSlice([]string{"apple", "banana", "chery"}))
		must.True(t, s.Contains("apple"))
		must.True(t, s.Contains("banana"))
		must.True(t, s.Contains("chery"))
		must.False(t, s.Contains("zucchini"))
	})

	t.Run("contains custom item", func(t *testing.T) {
		s := New[employee](10)
		must.True(t, s.Insert(employee{"mitchell", 1}))
		must.True(t, s.Insert(employee{"armon", 2}))
		must.True(t, s.Insert(employee{"jack", 3}))
		must.Contains[employee](t, employee{"mitchell", 1}, s)
		must.Contains[employee](t, employee{"armon", 2}, s)
		must.Contains[employee](t, employee{"jack", 3}, s)
		must.NotContains[employee](t, employee{"dave", 27}, s)
	})
}

func TestSet_ContainsAll(t *testing.T) {
	t.Run("contains subset", func(t *testing.T) {
		s := New[int](10)
		must.True(t, s.InsertSlice([]int{1, 2, 3, 4, 5}))
		must.True(t, s.ContainsAll([]int{1, 3, 5}))
	})

	t.Run("contains missing", func(t *testing.T) {
		s := New[int](10)
		must.True(t, s.InsertSlice([]int{1, 2, 3, 4, 5}))
		must.False(t, s.ContainsAll([]int{1, 3, 5, 7}))
	})
}

func TestSet_ContainsSlice(t *testing.T) {
	t.Run("empty empty", func(t *testing.T) {
		a := New[int](0)
		b := make([]int, 0)
		must.True(t, a.ContainsSlice(b))
	})

	t.Run("empty some", func(t *testing.T) {
		a := New[int](0)
		b := []int{1, 2, 3}
		must.False(t, a.ContainsSlice(b))
	})

	t.Run("some empty", func(t *testing.T) {
		a := From[int]([]int{1, 2, 3})
		b := make([]int, 0)
		must.False(t, a.ContainsSlice(b))
	})

	t.Run("equal", func(t *testing.T) {
		a := From[int]([]int{1, 2, 3})
		b := []int{3, 2, 1}
		must.True(t, a.ContainsSlice(b))
	})

	t.Run("not equal", func(t *testing.T) {
		a := From[int]([]int{1, 2, 3})
		b := []int{2, 3, 4}
		must.False(t, a.ContainsSlice(b))
	})

	t.Run("subset", func(t *testing.T) {
		a := From[int]([]int{1, 2, 3, 4, 5})
		b := []int{2, 3, 4}
		must.False(t, a.ContainsSlice(b))
	})

	t.Run("superset", func(t *testing.T) {
		a := From[int]([]int{2, 3, 4})
		b := []int{1, 2, 3, 4, 5}
		must.False(t, a.ContainsSlice(b))
	})

	t.Run("duplicates", func(t *testing.T) {
		a := From[int]([]int{1, 2, 3, 4, 5})
		b := []int{1, 2, 2, 3, 3, 4, 5}
		must.True(t, a.ContainsSlice(b))
	})
}

func TestSet_Size(t *testing.T) {
	t.Run("size empty", func(t *testing.T) {
		s := New[int](10)
		must.Zero(t, s.Size())
	})

	t.Run("size not empty", func(t *testing.T) {
		s := New[int](10)
		must.True(t, s.Insert(1))
		must.True(t, s.Insert(2))
		must.Eq(t, 2, s.Size())
	})
}

func TestSet_Empty(t *testing.T) {
	t.Run("is empty", func(t *testing.T) {
		s := New[int](10)
		must.Empty(t, s)
	})

	t.Run("is not empty", func(t *testing.T) {
		s := From[int]([]int{1, 2, 3, 4})
		must.NotEmpty(t, s)
	})
}

func TestSet_Union(t *testing.T) {
	t.Run("empty ∪ empty", func(t *testing.T) {
		a := New[int](0)
		b := New[int](10)
		union := a.Union(b)
		must.MapEmpty(t, union.items)
	})

	t.Run("empty ∪ set", func(t *testing.T) {
		a := New[int](10)
		b := New[int](10)
		b.InsertSlice([]int{1, 2, 3, 4, 5})
		union := a.Union(b)
		must.MapContainsKeys(t, union.items, []int{1, 2, 3, 4, 5})
	})

	t.Run("set ∪ empty", func(t *testing.T) {
		a := New[int](10)
		a.InsertSlice([]int{1, 2, 3, 4, 5})
		b := New[int](10)
		union := a.Union(b)
		must.MapContainsKeys(t, union.items, []int{1, 2, 3, 4, 5})
	})

	t.Run("set ∪ other", func(t *testing.T) {
		a := New[int](10)
		must.True(t, a.InsertSlice([]int{2, 4, 6, 8}))
		b := New[int](10)
		must.True(t, b.InsertSlice([]int{4, 5, 6}))
		union := a.Union(b)
		must.MapContainsKeys(t, union.items, []int{2, 4, 5, 6, 8})
	})
}

func TestSet_Difference(t *testing.T) {
	t.Run("empty \\ empty", func(t *testing.T) {
		a := New[int](10)
		b := New[int](10)
		diff := a.Difference(b)
		must.MapEmpty(t, diff.items)
	})

	t.Run("empty \\ set", func(t *testing.T) {
		a := New[int](10)
		b := From([]int{1, 2, 3, 4, 5})
		diff := a.Difference(b)
		must.MapEmpty(t, diff.items)
	})

	t.Run("set \\ empty", func(t *testing.T) {
		a := From([]int{1, 2, 3, 4, 5})
		b := New[int](10)
		diff := a.Difference(b)
		must.MapContainsKeys(t, diff.items, []int{1, 2, 3, 4, 5})
	})

	t.Run("set \\ other", func(t *testing.T) {
		a := From([]int{1, 2, 3, 4, 5, 6, 7, 8})
		b := From([]int{2, 4, 6, 8, 10, 12})
		diff := a.Difference(b)
		must.MapContainsKeys(t, diff.items, []int{1, 3, 5, 7})
	})
}

func TestSet_Intersect(t *testing.T) {
	t.Run("empty ∩ empty", func(t *testing.T) {
		a := New[int](10)
		b := New[int](10)
		intersect := a.Intersect(b)
		must.MapEmpty(t, intersect.items)
	})

	t.Run("set ∩ empty", func(t *testing.T) {
		a := From[int]([]int{1, 2, 3})
		b := New[int](10)
		intersect := a.Intersect(b)
		must.MapEmpty(t, intersect.items)
	})

	t.Run("empty ∩ set", func(t *testing.T) {
		a := New[int](10)
		b := From[int]([]int{1, 2, 3})
		intersect := a.Intersect(b)
		must.MapEmpty(t, intersect.items)
	})

	t.Run("big ∩ small", func(t *testing.T) {
		a := From[int]([]int{2, 3, 4, 6, 8})
		b := From[int]([]int{4, 5, 6, 7})
		intersect := a.Intersect(b)
		must.MapContainsKeys(t, intersect.items, []int{4, 6})
	})

	t.Run("small ∩ big", func(t *testing.T) {
		a := From[int]([]int{4, 5, 6, 7})
		b := From[int]([]int{2, 3, 4, 6, 8})
		intersect := a.Intersect(b)
		must.MapContainsKeys(t, intersect.items, []int{4, 6})
	})
}

func TestSet_Remove(t *testing.T) {
	t.Run("empty remove item", func(t *testing.T) {
		s := New[int](10)
		must.False(t, s.Remove(32))
		must.MapEmpty(t, s.items)
	})

	t.Run("set remove item", func(t *testing.T) {
		s := From[string]([]string{"apple", "banana", "cherry"})
		must.True(t, s.Remove("banana"))
		must.MapContainsKeys(t, s.items, []string{"apple", "cherry"})
	})

	t.Run("set remove missing", func(t *testing.T) {
		s := From[string]([]string{"apple", "banana", "cherry"})
		must.False(t, s.Remove("zucchini"))
		must.MapContainsKeys(t, s.items, []string{"apple", "banana", "cherry"})
	})
}

func TestSet_RemoveSlice(t *testing.T) {
	t.Run("empty remove all", func(t *testing.T) {
		s := New[int](10)
		must.False(t, s.RemoveSlice([]int{1, 2, 3}))
		must.MapEmpty(t, s.items)
	})

	t.Run("set remove nothing", func(t *testing.T) {
		s := From[int]([]int{1, 2, 3})
		must.False(t, s.RemoveSlice(nil))
		must.MapContainsKeys(t, s.items, []int{1, 2, 3})
	})

	t.Run("set remove some", func(t *testing.T) {
		s := From[int]([]int{1, 2, 3, 4, 5, 6})
		must.True(t, s.RemoveSlice([]int{5, 6, 7, 8, 9}))
		must.MapContainsKeys(t, s.items, []int{1, 2, 3, 4})
	})
}

func TestSet_RemoveSet(t *testing.T) {
	t.Run("empty remove empty", func(t *testing.T) {
		a := New[int](0)
		b := New[int](0)
		must.False(t, a.RemoveSet(b))
		must.MapEmpty(t, a.items)
	})

	t.Run("empty remove some", func(t *testing.T) {
		a := New[int](0)
		b := From[int]([]int{1, 2, 3, 4})
		must.False(t, a.RemoveSet(b))
		must.MapEmpty(t, a.items)
	})

	t.Run("set remove some", func(t *testing.T) {
		a := From[int]([]int{1, 2, 3, 4, 5, 6, 7, 8})
		b := From[int]([]int{2, 4, 6, 8})
		must.True(t, a.RemoveSet(b))
		must.MapContainsKeys(t, a.items, []int{1, 3, 5, 7})
	})
}

func TestSet_RemoveFunc(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		a := New[int](10)
		modified := a.RemoveFunc(func(i int) bool {
			return i%2 == 0
		})
		must.Empty(t, a)
		must.False(t, modified)
	})

	t.Run("none match", func(t *testing.T) {
		a := From[int]([]int{1, 3, 5, 7, 9})
		modified := a.RemoveFunc(func(i int) bool {
			return i%2 == 0
		})
		must.True(t, a.ContainsSlice([]int{1, 3, 5, 7, 9}))
		must.False(t, modified)
	})

	t.Run("some match", func(t *testing.T) {
		a := From[int]([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})
		modified := a.RemoveFunc(func(i int) bool {
			return i%2 == 0
		})
		must.True(t, a.ContainsSlice([]int{1, 3, 5, 7, 9}))
		must.True(t, modified)
	})

	t.Run("all match", func(t *testing.T) {
		a := From[int]([]int{1, 3, 5, 7, 9})
		modified := a.RemoveFunc(func(i int) bool {
			return i%2 != 0
		})
		must.Empty(t, a)
		must.True(t, modified)
	})
}

func TestSet_Copy(t *testing.T) {
	t.Run("copy empty", func(t *testing.T) {
		a := New[int](0)
		b := a.Copy()
		must.MapEmpty(t, b.items)
		must.True(t, b.Insert(3))
		must.MapEmpty(t, a.items)
		must.MapContainsKeys(t, b.items, []int{3})
	})

	t.Run("copy some", func(t *testing.T) {
		a := From[int]([]int{1, 2, 3, 4})
		b := a.Copy()
		must.MapContainsKeys(t, b.items, []int{1, 2, 3, 4})
		must.True(t, b.RemoveSlice([]int{1, 3}))
		must.MapContainsKeys(t, b.items, []int{2, 4})
		must.MapContainsKeys(t, a.items, []int{1, 2, 3, 4})
	})
}

func TestSet_Slice(t *testing.T) {
	t.Run("slice empty", func(t *testing.T) {
		a := New[string](10)
		l := a.Slice()
		must.SliceEmpty(t, l)
	})

	t.Run("slice set", func(t *testing.T) {
		a := From([]string{"apple", "banana", "cherry"})
		l := a.Slice()
		must.Len(t, 3, l)
		must.SliceContains(t, l, "apple")
		must.SliceContains(t, l, "banana")
		must.SliceContains(t, l, "cherry")
	})
}

func TestSet_String(t *testing.T) {
	t.Run("ints", func(t *testing.T) {
		a := From([]int{1, 2, 3})
		result := a.String()
		must.Eq(t, "[1 2 3]", result)
	})

	t.Run("custom", func(t *testing.T) {
		a := From([]*employee{{"bob", 2}, {"alice", 1}, {"carl", 3}})
		result := a.String()
		must.Eq(t, "[(1 alice) (2 bob) (3 carl)]", result)
	})
}

func TestSet_StringFunc(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		a := New[string](10)
		s := a.StringFunc(nil)
		must.Eq(t, "[]", s)
	})

	t.Run("int", func(t *testing.T) {
		a := From([]int{5, 2, 5, 1, 3})
		s := a.StringFunc(func(i int) string {
			return fmt.Sprintf("%d", i)
		})
		must.Eq(t, "[1 2 3 5]", s)
	})

	t.Run("custom", func(t *testing.T) {
		a := From([]employee{
			{"mitchell", 1},
			{"jack", 3},
			{"armon", 2},
		})
		s := a.StringFunc(func(e employee) string {
			return fmt.Sprintf("(%d %s)", e.id, e.name)
		})
		must.Eq(t, "[(1 mitchell) (2 armon) (3 jack)]", s)
	})
}

func TestSet_Equal(t *testing.T) {
	t.Run("empty empty", func(t *testing.T) {
		a := New[int](0)
		b := New[int](10)
		must.True(t, a.Equal(b))
	})

	t.Run("empty some", func(t *testing.T) {
		a := New[int](0)
		b := From[int]([]int{1, 2, 3})
		must.False(t, a.Equal(b))
	})

	t.Run("same", func(t *testing.T) {
		a := From[int]([]int{3, 2, 1})
		b := From[int]([]int{1, 2, 3})
		must.True(t, a.Equal(b))
	})

	t.Run("subset", func(t *testing.T) {
		a := From[int]([]int{2, 3})
		b := From[int]([]int{1, 2, 3})
		must.False(t, a.Equal(b))
		must.False(t, b.Equal(a))
	})
}

func TestSet_Subset(t *testing.T) {
	t.Run("empty empty", func(t *testing.T) {
		a := New[int](0)
		b := New[int](0)
		must.True(t, a.Subset(b))
	})

	t.Run("empty some", func(t *testing.T) {
		a := New[int](0)
		b := From[int]([]int{1, 2, 3})
		must.False(t, a.Subset(b))
	})

	t.Run("some empty", func(t *testing.T) {
		a := From[int]([]int{1, 2, 3})
		b := New[int](0)
		must.True(t, a.Subset(b))
	})

	t.Run("equal", func(t *testing.T) {
		a := From[int]([]int{1, 2, 3})
		b := From[int]([]int{2, 3, 1})
		must.True(t, a.Subset(b))
	})

	t.Run("subset", func(t *testing.T) {
		a := From[int]([]int{1, 2, 3})
		b := From[int]([]int{3, 1})
		must.True(t, a.Subset(b))
	})

	t.Run("superset", func(t *testing.T) {
		a := From[int]([]int{1, 2, 3})
		b := From[int]([]int{3, 1, 2, 4})
		must.False(t, a.Subset(b))
	})
}

func TestSet_EqualSlice(t *testing.T) {
	t.Run("empty empty", func(t *testing.T) {
		a := New[int](0)
		b := make([]int, 0)
		must.True(t, a.EqualSlice(b))
	})

	t.Run("empty some", func(t *testing.T) {
		a := New[int](0)
		b := []int{1, 2, 3}
		must.False(t, a.EqualSlice(b))
	})

	t.Run("some empty", func(t *testing.T) {
		a := From[int]([]int{1, 2, 3})
		b := make([]int, 0)
		must.False(t, a.EqualSlice(b))
	})

	t.Run("equal", func(t *testing.T) {
		a := From[int]([]int{1, 2, 3})
		b := []int{3, 2, 1}
		must.True(t, a.EqualSlice(b))
	})

	t.Run("not equal", func(t *testing.T) {
		a := From[int]([]int{1, 2, 3})
		b := []int{2, 3, 4}
		must.False(t, a.EqualSlice(b))
	})

	t.Run("subset", func(t *testing.T) {
		a := From[int]([]int{1, 2, 3, 4, 5})
		b := []int{2, 3, 4}
		must.False(t, a.EqualSlice(b))
	})

	t.Run("superset", func(t *testing.T) {
		a := From[int]([]int{2, 3, 4})
		b := []int{1, 2, 3, 4, 5}
		must.False(t, a.EqualSlice(b))
	})

	t.Run("duplicates", func(t *testing.T) {
		a := From[int]([]int{1, 2, 3, 4, 5})
		b := []int{1, 2, 2, 3, 3, 4, 5}
		must.False(t, a.EqualSlice(b))
	})
}

func TestSet_ForEach(t *testing.T) {
	s := From[int]([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	var result []int
	evens := func(i int) bool {
		if i%2 == 0 {
			result = append(result, i)
		}
		return true
	}
	s.ForEach(evens)
	must.Eq(t, []int{0, 2, 4, 6, 8}, result)
}
