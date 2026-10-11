// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package set

import (
	"cmp"
	"encoding/json"
	"testing"

	"github.com/shoenig/test/must"
)

func TestSerialization(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		set := New[int](3)
		set.InsertSlice([]int{1, 2, 3})
		bs, err := json.Marshal(set)
		must.NoError(t, err)
		must.StrContains(t, string(bs), "1")
		must.StrContains(t, string(bs), "2")
		must.StrContains(t, string(bs), "3")

		dstSet := New[int](3)
		err = json.Unmarshal(bs, dstSet)
		must.NoError(t, err)
		must.MapEq(t, dstSet.items, set.items)
	})

	t.Run("HashSet", func(t *testing.T) {
		set := NewHashSet[*company, string](10)
		set.InsertSlice([]*company{c1, c2, c3})
		bs, err := json.Marshal(set)
		must.NoError(t, err)
		must.StrContains(t, string(bs), `"street":1`)
		must.StrContains(t, string(bs), `"street":2`)
		must.StrContains(t, string(bs), `"street":3`)

		dstSet := NewHashSet[*company, string](10)
		err = json.Unmarshal(bs, dstSet)
		must.NoError(t, err)
		must.MapEqual(t, dstSet.items, set.items)
	})

	t.Run("TreeSet", func(t *testing.T) {
		set := NewTreeSet[int](cmp.Compare[int])
		set.InsertSlice([]int{10, 3, 13})
		bs, err := json.Marshal(set)
		must.NoError(t, err)
		must.StrContains(t, string(bs), "10")
		must.StrContains(t, string(bs), "3")
		must.StrContains(t, string(bs), "13")

		dstSet := NewTreeSet[int](cmp.Compare[int])
		err = json.Unmarshal(bs, dstSet)
		must.NoError(t, err)
		must.Eq(t, set.Slice(), dstSet.Slice())
	})
}

func TestUnmarshalReplacesContents(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		s := From([]int{1, 2})
		must.NoError(t, json.Unmarshal([]byte(`[3]`), s))
		must.Eq(t, 1, s.Size())
		must.True(t, s.Contains(3))
		must.False(t, s.Contains(1))

		must.NoError(t, json.Unmarshal([]byte(`[]`), s))
		must.True(t, s.Empty())

		s.Insert(1)
		must.Error(t, json.Unmarshal([]byte(`{`), s))
		must.True(t, s.Contains(1))
	})

	t.Run("HashSet", func(t *testing.T) {
		s := NewHashSetFunc[string, string](2, func(v string) string { return v })
		s.InsertSlice([]string{"a", "b"})
		must.NoError(t, json.Unmarshal([]byte(`["c"]`), s))
		must.Eq(t, 1, s.Size())
		must.True(t, s.Contains("c"))
		must.False(t, s.Contains("a"))
	})

	t.Run("TreeSet", func(t *testing.T) {
		s := NewTreeSet[int](cmp.Compare[int])
		s.InsertSlice([]int{1, 2})
		must.NoError(t, json.Unmarshal([]byte(`[3]`), s))
		must.Eq(t, []int{3}, s.Slice())
	})
}
