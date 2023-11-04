// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package set

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/shoenig/test/must"
)

// assertion that HashSet[T] implements Collection[T]
var _ Collection[*company] = (*HashSet[*company, string])(nil)

// company is an example type that is not comparable, and implements Hash() string
type company struct {
	_       func() // not comparable
	address string
	floor   int
}

func (c *company) Equal(o *company) bool {
	return c.address == o.address && c.floor == o.floor
}

func (c *company) Hash() string {
	return fmt.Sprintf("%s:%d", c.address, c.floor)
}

func (c *company) String() string {
	return fmt.Sprintf("<%s %d>", c.address, c.floor)
}
func (c *company) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("{\"%s\":%d}", c.address, c.floor)), nil
}

func (c *company) UnmarshalJSON(data []byte) error {
	s := strings.TrimLeft(string(data), "{")
	s = strings.TrimRight(s, "}")

	splitSlice := strings.Split(s, ":")
	address := strings.TrimLeft(splitSlice[0], "\"")
	address = strings.TrimRight(address, "\"")
	c.address = address
	c.floor, _ = strconv.Atoi(splitSlice[1])
	return nil
}

var (
	c1  = &company{address: "street", floor: 1}
	c2  = &company{address: "street", floor: 2}
	c3  = &company{address: "street", floor: 3}
	c4  = &company{address: "street", floor: 4}
	c5  = &company{address: "street", floor: 5}
	c6  = &company{address: "street", floor: 6}
	c7  = &company{address: "street", floor: 7}
	c8  = &company{address: "street", floor: 8}
	c10 = &company{address: "street", floor: 10}
)

// coded is an example type that maintains its own hash code, implementing Hash() int
type coded struct {
	i int // internal hash code
}

func (c *coded) Hash() int {
	return c.i
}

var (
	s1 = &coded{i: 1}
	s2 = &coded{i: 2}
	s3 = &coded{i: 3}
)

func TestHashSet_New(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		s := NewHashSet[*company, string](1)
		must.MapEmpty(t, s.items)
	})

	t.Run("zero", func(t *testing.T) {
		s := NewHashSet[*company, string](0)
		must.MapEmpty(t, s.items)
	})

	t.Run("negative", func(t *testing.T) {
		s := NewHashSet[*company, string](-1)
		must.MapEmpty(t, s.items)
	})
}

func TestHashSet_Insert(t *testing.T) {
	t.Run("one", func(t *testing.T) {
		s := NewHashSet[*company, string](1)
		must.True(t, s.Insert(c1))
		must.MapContainsKeys(t, s.items, []string{"street:1"})
	})

	t.Run("re-insert", func(t *testing.T) {
		s := NewHashSet[*company, string](1)
		must.True(t, s.Insert(c1))
		must.False(t, s.Insert(c1))
		must.MapContainsKeys(t, s.items, []string{"street:1"})
	})

	t.Run("insert several", func(t *testing.T) {
		s := NewHashSet[*company, string](3)
		must.True(t, s.Insert(c1))
		must.True(t, s.Insert(c2))
		must.True(t, s.Insert(c3))
		must.MapContainsKeys(t, s.items, []string{
			"street:1", "street:2", "street:3",
		})
	})
}

func TestHashSet_InsertSlice(t *testing.T) {
	t.Run("insert none", func(t *testing.T) {
		empty := NewHashSet[*company, string](0)
		must.False(t, empty.InsertSlice(nil))
		must.MapEmpty(t, empty.items)
	})

	t.Run("insert some", func(t *testing.T) {
		s := NewHashSet[*company, string](0)
		must.True(t, s.InsertSlice([]*company{c1, c2, c3}))
		must.MapContainsKeys(t, s.items, []string{
			"street:1", "street:2", "street:3",
		})
	})
}

func TestHashSet_InsertSet(t *testing.T) {
	t.Run("insert empty", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		b := NewHashSet[*company, string](0)
		must.False(t, a.InsertSet(b))
		must.MapContainsKeys(t, a.items, []string{
			"street:1", "street:2", "street:3",
		})
	})

	t.Run("insert some", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		b := HashSetFrom[*company, string]([]*company{c3, c4, c5})
		a.InsertSet(b)
		must.MapContainsKeys(t, a.items, []string{
			"street:1", "street:2", "street:3", "street:4", "street:5",
		})
	})
}

func TestHashSet_Remove(t *testing.T) {
	t.Run("empty remove item", func(t *testing.T) {
		s := NewHashSet[*company, string](10)
		must.False(t, s.Remove(c1))
		must.MapEmpty(t, s.items)
	})

	t.Run("set remove item", func(t *testing.T) {
		s := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		must.True(t, s.Remove(c2))
		must.MapContainsKeys(t, s.items, []string{
			"street:1", "street:3",
		})
	})

	t.Run("set remove missing", func(t *testing.T) {
		s := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		must.False(t, s.Remove(c4))
		must.MapContainsKeys(t, s.items, []string{
			"street:1", "street:2", "street:3",
		})
	})
}

func TestHashSet_RemoveSlice(t *testing.T) {
	t.Run("empty remove all", func(t *testing.T) {
		s := NewHashSet[*company, string](0)
		must.False(t, s.RemoveSlice([]*company{c1, c2, c3}))
		must.MapEmpty(t, s.items)
	})

	t.Run("set remove nothing", func(t *testing.T) {
		s := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		must.False(t, s.RemoveSlice([]*company{c4, c5}))
		must.MapContainsKeys(t, s.items, []string{
			"street:1", "street:2", "street:3",
		})
	})

	t.Run("set remove some", func(t *testing.T) {
		s := HashSetFrom[*company, string]([]*company{c1, c2, c3, c4, c5})
		must.True(t, s.RemoveSlice([]*company{c4, c2}))
		must.MapContainsKeys(t, s.items, []string{
			"street:1", "street:3", "street:5",
		})
	})
}

func TestHashSet_RemoveSet(t *testing.T) {
	t.Run("empty remove empty", func(t *testing.T) {
		a := NewHashSet[*company, string](0)
		b := NewHashSet[*company, string](0)
		must.False(t, a.RemoveSet(b))
		must.MapEmpty(t, a.items)
	})

	t.Run("empty remove some", func(t *testing.T) {
		a := NewHashSet[*company, string](0)
		b := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		must.False(t, a.RemoveSet(b))
		must.MapEmpty(t, a.items)
	})

	t.Run("set remove some", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3, c4, c5})
		b := HashSetFrom[*company, string]([]*company{c2, c3})
		must.True(t, a.RemoveSet(b))
		must.MapContainsKeys(t, a.items, []string{
			"street:1", "street:4", "street:5",
		})
	})
}

func TestHashSet_RemoveFunc(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := NewHashSet[*company, string](10)
		modified := s.RemoveFunc(func(c *company) bool {
			return c.floor > 3
		})
		must.Empty(t, s)
		must.False(t, modified)
	})

	t.Run("none match", func(t *testing.T) {
		s := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		modified := s.RemoveFunc(func(c *company) bool {
			return c.floor > 3
		})
		must.False(t, modified)
		must.Size(t, 3, s)
	})

	t.Run("some match", func(t *testing.T) {
		s := HashSetFrom[*company, string]([]*company{c1, c2, c3, c4, c5, c6})
		modified := s.RemoveFunc(func(c *company) bool {
			return c.floor > 3
		})
		must.True(t, modified)
		must.Size(t, 3, s)
		must.Contains[*company](t, c1, s)
		must.Contains[*company](t, c2, s)
		must.Contains[*company](t, c3, s)
	})

	t.Run("all match", func(t *testing.T) {
		s := HashSetFrom[*company, string]([]*company{c1, c2, c3, c4, c5, c6})
		modified := s.RemoveFunc(func(c *company) bool {
			return c.floor >= 0
		})
		must.True(t, modified)
		must.Empty(t, s)
	})
}

func TestHashSet_Contains(t *testing.T) {
	t.Run("empty contains", func(t *testing.T) {
		a := NewHashSet[*company, string](0)
		must.False(t, a.Contains(c1))
	})

	t.Run("not contains", func(t *testing.T) {
		s := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		must.False(t, s.Contains(c4))
	})

	t.Run("does contain", func(t *testing.T) {
		s := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		must.True(t, s.Contains(c1))
	})
}

func TestHashSet_ContainsSlice(t *testing.T) {
	t.Run("empty empty", func(t *testing.T) {
		a := NewHashSet[*company, string](0)
		b := make([]*company, 0)
		must.True(t, a.ContainsSlice(b))
	})

	t.Run("empty some", func(t *testing.T) {
		a := NewHashSet[*company, string](0)
		b := []*company{c1, c2, c3}
		must.False(t, a.ContainsSlice(b))
	})

	t.Run("some empty", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		b := make([]*company, 0)
		must.False(t, a.ContainsSlice(b))
	})

	t.Run("equal", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		b := []*company{c3, c2, c1}
		must.True(t, a.ContainsSlice(b))
	})

	t.Run("not equal", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		b := []*company{c2, c3, c4}
		must.False(t, a.ContainsSlice(b))
	})

	t.Run("subset", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3, c4, c5})
		b := []*company{c2, c3, c4}
		must.False(t, a.ContainsSlice(b))
	})

	t.Run("superset", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c2, c3, c4})
		b := []*company{c1, c2, c3, c4, c5}
		must.False(t, a.ContainsSlice(b))
	})

	t.Run("duplicates", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3, c4, c5})
		b := []*company{c1, c2, c2, c3, c3, c4, c5}
		must.True(t, a.ContainsSlice(b))
	})
}

func TestHashSet_Subset(t *testing.T) {
	t.Run("empty empty", func(t *testing.T) {
		a := NewHashSet[*company, string](0)
		b := NewHashSet[*company, string](0)
		must.True(t, a.Subset(b))
	})

	t.Run("empty some", func(t *testing.T) {
		a := NewHashSet[*company, string](0)
		b := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		must.False(t, a.Subset(b))
	})

	t.Run("some empty", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		b := NewHashSet[*company, string](0)
		must.True(t, a.Subset(b))
	})

	t.Run("equal", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		b := HashSetFrom[*company, string]([]*company{c2, c3, c1})
		must.True(t, a.Subset(b))
	})

	t.Run("subset", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		b := HashSetFrom[*company, string]([]*company{c3, c1})
		must.True(t, a.Subset(b))
	})

	t.Run("superset", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		b := HashSetFrom[*company, string]([]*company{c3, c1, c2, c4})
		must.False(t, a.Subset(b))
	})
}

func TestHashSet_ProperSubset(t *testing.T) {
	t.Run("empty empty", func(t *testing.T) {
		a := NewHashSet[*company, string](0)
		b := NewHashSet[*company, string](0)
		must.False(t, a.ProperSubset(b))
	})

	t.Run("empty some", func(t *testing.T) {
		a := NewHashSet[*company, string](0)
		b := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		must.False(t, a.ProperSubset(b))
	})

	t.Run("some empty", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		b := NewHashSet[*company, string](0)
		must.True(t, a.ProperSubset(b))
	})

	t.Run("equal", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		b := HashSetFrom[*company, string]([]*company{c2, c3, c1})
		must.False(t, a.ProperSubset(b))
	})

	t.Run("subset", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		b := HashSetFrom[*company, string]([]*company{c3, c1})
		must.True(t, a.ProperSubset(b))
	})

	t.Run("superset", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		b := HashSetFrom[*company, string]([]*company{c3, c1, c2, c4})
		must.False(t, a.ProperSubset(b))
	})
}

func TestHashSet_Size(t *testing.T) {
	t.Run("size empty", func(t *testing.T) {
		s := NewHashSet[*company, string](10)
		must.Zero(t, s.Size())
	})

	t.Run("size not empty", func(t *testing.T) {
		s := NewHashSet[*company, string](10)
		must.True(t, s.Insert(c1))
		must.True(t, s.Insert(c2))
		must.Eq(t, 2, s.Size())
	})
}

func TestHashSet_Empty(t *testing.T) {
	t.Run("is empty", func(t *testing.T) {
		s := NewHashSet[*company, string](10)
		must.Empty(t, s)
	})

	t.Run("is not empty", func(t *testing.T) {
		s := NewHashSet[*company, string](10)
		must.True(t, s.Insert(c1))
		must.True(t, s.Insert(c2))
		must.NotEmpty(t, s)
	})
}

func TestHashSet_Difference(t *testing.T) {
	t.Run("empty \\ empty", func(t *testing.T) {
		a := NewHashSet[*company, string](10)
		b := NewHashSet[*company, string](10)
		diff := a.Difference(b).(*HashSet[*company, string])
		must.MapEmpty(t, diff.items)
	})

	t.Run("empty \\ set", func(t *testing.T) {
		a := NewHashSet[*company, string](10)
		b := HashSetFrom[*company, string]([]*company{c1, c2, c3, c4, c5})
		diff := a.Difference(b).(*HashSet[*company, string])
		must.MapEmpty(t, diff.items)
	})

	t.Run("set \\ empty", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3, c4, c5})
		b := NewHashSet[*company, string](10)
		diff := a.Difference(b).(*HashSet[*company, string])
		must.MapContainsKeys(t, diff.items, []string{
			"street:1", "street:2", "street:3", "street:4", "street:5",
		})
	})

	t.Run("set \\ other", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3, c4, c5, c6, c7, c8})
		b := HashSetFrom[*company, string]([]*company{c2, c4, c6, c8, c10, c10})
		diff := a.Difference(b).(*HashSet[*company, string])
		must.MapContainsKeys(t, diff.items, []string{
			"street:1", "street:3", "street:5", "street:7",
		})
	})
}

func TestHashSet_Intersect(t *testing.T) {
	t.Run("empty ∩ empty", func(t *testing.T) {
		a := NewHashSet[*company, string](10)
		b := NewHashSet[*company, string](10)
		intersect := a.Intersect(b).(*HashSet[*company, string])
		must.MapEmpty(t, intersect.items)
	})

	t.Run("set ∩ empty", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		b := NewHashSet[*company, string](10)
		intersect := a.Intersect(b).(*HashSet[*company, string])
		must.MapEmpty(t, intersect.items)
	})

	t.Run("empty ∩ set", func(t *testing.T) {
		a := NewHashSet[*company, string](10)
		b := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		intersect := a.Intersect(b).(*HashSet[*company, string])
		must.MapEmpty(t, intersect.items)
	})

	t.Run("big ∩ small", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c2, c3, c4, c6, c8})
		b := HashSetFrom[*company, string]([]*company{c4, c5, c6, c7})
		intersect := a.Intersect(b).(*HashSet[*company, string])
		must.MapContainsKeys(t, intersect.items, []string{
			"street:4", "street:6",
		})
	})

	t.Run("small ∩ big", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c4, c5, c6, c7})
		b := HashSetFrom[*company, string]([]*company{c2, c3, c4, c6, c8})
		intersect := a.Intersect(b).(*HashSet[*company, string])
		must.MapContainsKeys(t, intersect.items, []string{
			"street:4", "street:6",
		})
	})
}

type special struct {
	hash    string
	version int // not part of the hash
}

func (s *special) Hash() string {
	return s.hash
}

func TestHashSet_Union(t *testing.T) {
	t.Run("empty ∪ empty", func(t *testing.T) {
		a := NewHashSet[*company, string](10)
		b := NewHashSet[*company, string](10)
		union := a.Union(b).(*HashSet[*company, string])
		must.MapEmpty(t, union.items)
	})

	t.Run("set ∪ empty", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		b := NewHashSet[*company, string](10)
		union := a.Union(b).(*HashSet[*company, string])
		must.MapContainsKeys(t, union.items, []string{
			"street:1", "street:2", "street:3",
		})
	})

	t.Run("empty ∪ set", func(t *testing.T) {
		a := NewHashSet[*company, string](10)
		b := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		union := a.Union(b).(*HashSet[*company, string])
		must.MapContainsKeys(t, union.items, []string{
			"street:1", "street:2", "street:3",
		})
	})

	t.Run("big ∪ small", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c4, c5, c6, c7})
		b := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		union := a.Union(b).(*HashSet[*company, string])
		must.MapContainsKeys(t, union.items, []string{
			"street:1", "street:2", "street:3", "street:4", "street:5", "street:6", "street:7",
		})
	})

	t.Run("overlap", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c4, c5, c6, c7})
		b := HashSetFrom[*company, string]([]*company{c1, c2, c3, c4, c5})
		union := a.Union(b).(*HashSet[*company, string])
		must.MapContainsKeys(t, union.items, []string{
			"street:1", "street:2", "street:3", "street:4", "street:5", "street:6", "street:7",
		})
	})

	t.Run("overlap priority", func(t *testing.T) {
		x1 := &special{hash: "x", version: 1}
		x2 := &special{hash: "x", version: 2}

		a := HashSetFrom[*special, string]([]*special{x1})
		b := HashSetFrom[*special, string]([]*special{x2})

		// receiver elements take priority over argument elements
		union1 := a.Union(b).(*HashSet[*special, string])
		must.MapContainsValues(t, union1.items, []*special{
			x1,
		})

		// checking in the other direction
		union2 := b.Union(a).(*HashSet[*special, string])
		must.MapContainsValues(t, union2.items, []*special{
			x2,
		})
	})
}

func TestHashSet_Copy(t *testing.T) {
	t.Run("copy empty", func(t *testing.T) {
		a := NewHashSet[*company, string](0)
		b := a.Copy()
		must.MapEmpty(t, b.items)
		must.True(t, b.Insert(c3))
		must.MapEmpty(t, a.items)
		must.MapContainsKeys(t, b.items, []string{"street:3"})
	})

	t.Run("copy some", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3, c4})
		b := a.Copy()
		must.MapContainsKeys(t, b.items, []string{
			"street:1", "street:2", "street:3", "street:4",
		})
		must.True(t, b.RemoveSlice([]*company{c1, c3}))
		must.MapContainsKeys(t, b.items, []string{"street:2", "street:4"})
		must.MapContainsKeys(t, a.items, []string{
			"street:1", "street:2", "street:3", "street:4",
		})
	})
}

func TestHashSet_Slice(t *testing.T) {
	t.Run("slice empty", func(t *testing.T) {
		a := NewHashSet[*company, string](10)
		l := a.Slice()
		must.SliceEmpty(t, l)
	})

	t.Run("slice set", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2})
		l := a.Slice()
		must.Len(t, 2, l)
		must.SliceContainsEqual(t, l, c1)
		must.SliceContainsEqual(t, l, c2)
	})
}

func TestHashSet_String(t *testing.T) {
	a := HashSetFrom[*company, string]([]*company{c2, c1})
	result := a.String()
	must.Eq(t, "[<street 1> <street 2>]", result)
}

func TestHashSet_StringFunc(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		a := NewHashSet[*company, string](10)
		s := a.StringFunc(nil)
		must.Eq(t, "[]", s)
	})

	t.Run("some", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2})
		s := a.StringFunc(func(c *company) string {
			return fmt.Sprintf("(%s %d)", c.address, c.floor)
		})
		must.Eq(t, "[(street 1) (street 2)]", s)
	})
}

func TestHashSet_EqualSet(t *testing.T) {
	t.Run("empty empty", func(t *testing.T) {
		a := NewHashSet[*company, string](0)
		b := NewHashSet[*company, string](0)
		must.True(t, a.EqualSet(b))
	})

	t.Run("different", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		b := HashSetFrom[*company, string]([]*company{c1, c2, c4})
		must.False(t, a.EqualSet(b))
		must.False(t, b.EqualSet(a))
	})
}

func TestHashSet_EqualSlice(t *testing.T) {
	t.Run("empty empty", func(t *testing.T) {
		a := NewHashSet[*company, string](0)
		b := make([]*company, 0)
		must.True(t, a.EqualSlice(b))
	})

	t.Run("empty some", func(t *testing.T) {
		a := NewHashSet[*company, string](0)
		b := []*company{c1, c2, c3}
		must.False(t, a.EqualSlice(b))
	})

	t.Run("some empty", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		b := make([]*company, 0)
		must.False(t, a.EqualSlice(b))
	})

	t.Run("equal", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		b := []*company{c3, c2, c1}
		must.True(t, a.EqualSlice(b))
	})

	t.Run("not equal", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3})
		b := []*company{c2, c3, c4}
		must.False(t, a.EqualSlice(b))
	})

	t.Run("subset", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3, c4, c5})
		b := []*company{c2, c3, c4}
		must.False(t, a.EqualSlice(b))
	})

	t.Run("superset", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c2, c3, c4})
		b := []*company{c1, c2, c3, c4, c5}
		must.False(t, a.EqualSlice(b))
	})

	t.Run("duplicates", func(t *testing.T) {
		a := HashSetFrom[*company, string]([]*company{c1, c2, c3, c4, c5})
		b := []*company{c1, c2, c2, c3, c3, c4, c5}
		must.True(t, a.EqualSlice(b))
	})
}

func TestHashSet_HashCode(t *testing.T) {
	a := NewHashSet[*coded, int](0)
	a.Insert(s1)
	a.Insert(s2)
	must.MapContainsKeys(t, a.items, []int{1, 2})
	must.True(t, a.Contains(s1))
	must.True(t, a.Contains(s2))
	must.False(t, a.Contains(s3))
}
