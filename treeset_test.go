// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package set

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/shoenig/test/must"
	"go.uber.org/goleak"
)

const (
	size = 1000
)

type token struct {
	id string
}

func (t *token) String() string {
	return t.id
}

func compareTokens(a, b *token) int {
	return Cmp(a.id, b.id)
}

var (
	tokenA = &token{id: "A"}
	tokenB = &token{id: "B"}
	tokenC = &token{id: "C"}
	tokenD = &token{id: "D"}
	tokenE = &token{id: "E"}
	tokenF = &token{id: "F"}
	tokenG = &token{id: "G"}
	tokenH = &token{id: "H"}
)

func TestNewTreeSet(t *testing.T) {
	ts := NewTreeSet[*token, Compare[*token]](compareTokens)
	must.NotNil(t, ts)
	ts.dump()
}

func TestTreeSetFrom(t *testing.T) {
	s := shuffle(ints(10))
	ts := TreeSetFrom[int, Compare[int]](s, Cmp[int])
	must.NotEmpty(t, ts)
}

func TestTreeSet_Empty(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		must.Empty(t, ts)
	})

	t.Run("not empty", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		ts.Insert(1)
		must.NotEmpty(t, ts)
	})
}

func TestTreeSet_Size(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		must.Size(t, 0, ts)
	})
	t.Run("one", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		ts.Insert(42)
		must.Size(t, 1, ts)
	})
	t.Run("ten", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		s := shuffle(ints(10))
		for i := 0; i < len(s); i++ {
			ts.Insert(s[i])
			must.Size(t, i+1, ts)
		}
		// insert again (all duplicates)
		s = shuffle(s)
		for i := 0; i < len(s); i++ {
			ts.Insert(s[i])
			must.Size(t, 10, ts)
		}
	})
}

func TestTreeSet_Insert_token(t *testing.T) {
	ts := NewTreeSet[*token, Compare[*token]](compareTokens)

	ts.Insert(tokenA)
	invariants(t, ts, compareTokens)

	ts.Insert(tokenB)
	invariants(t, ts, compareTokens)

	ts.Insert(tokenC)
	invariants(t, ts, compareTokens)

	ts.Insert(tokenD)
	invariants(t, ts, compareTokens)

	ts.Insert(tokenE)
	invariants(t, ts, compareTokens)

	ts.Insert(tokenF)
	invariants(t, ts, compareTokens)

	ts.Insert(tokenG)
	invariants(t, ts, compareTokens)

	ts.Insert(tokenH)
	invariants(t, ts, compareTokens)

	t.Log("dump: insert token")
	t.Log(ts.dump())
}

func TestTreeSet_Insert_int(t *testing.T) {
	cmp := Cmp[int]
	ts := NewTreeSet[int, Compare[int]](cmp)

	numbers := ints(size)
	random := shuffle(numbers)

	for _, i := range random {
		ts.Insert(i)
		invariants(t, ts, cmp)
	}

	t.Log("dump: insert int")
	t.Log(ts.dump())
}

func TestTreeSet_InsertSlice(t *testing.T) {
	cmp := Cmp[int]

	numbers := ints(size)
	random := shuffle(numbers)

	ts := NewTreeSet[int, Compare[int]](cmp)
	must.True(t, ts.InsertSlice(random))
	must.Eq(t, numbers, ts.Slice())
	must.False(t, ts.InsertSlice(numbers))
}

func TestTreeSet_InsertSet(t *testing.T) {
	cmp := Cmp[int]

	ts1 := TreeSetFrom[int, Compare[int]]([]int{1, 3, 5, 7, 9}, cmp)
	ts2 := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3}, cmp)

	must.True(t, ts1.InsertSet(ts2))
	must.Eq(t, []int{1, 2, 3, 5, 7, 9}, ts1.Slice())
	must.Eq(t, []int{1, 2, 3}, ts2.Slice())
}

func TestTreeSet_Remove_int(t *testing.T) {
	cmp := Cmp[int]
	ts := NewTreeSet[int, Compare[int]](cmp)

	numbers := ints(size)
	rnd := shuffle(numbers)

	// insert in random order
	for _, i := range rnd {
		ts.Insert(i)
	}

	invariants(t, ts, cmp)

	// reshuffle
	rnd = shuffle(rnd)

	// remove every element in random order
	for _, i := range rnd {
		removed := ts.Remove(i)
		t.Log("dump: remove", i)
		t.Log(ts.dump())
		must.True(t, removed)
		invariants(t, ts, cmp)
	}

	// all gone
	must.Empty(t, ts)
}

func TestTreeSet_RemoveSlice(t *testing.T) {
	cmp := Cmp[int]
	ts := NewTreeSet[int, Compare[int]](cmp)

	numbers := ints(size)
	random := shuffle(numbers)
	ts.InsertSlice(random)

	must.True(t, ts.RemoveSlice(numbers))
	must.Empty(t, ts)
}

func TestTreeSet_RemoveSet(t *testing.T) {
	cmp := Cmp[int]

	ts1 := NewTreeSet[int, Compare[int]](cmp)
	ts2 := NewTreeSet[int, Compare[int]](cmp)

	numbers := ints(size)
	random := shuffle(numbers)
	ts1.InsertSlice(random)

	random2 := shuffle(numbers[5:])
	ts2.InsertSlice(random2)

	ts1.RemoveSet(ts2)
	result := ts1.Slice()
	must.Eq(t, []int{1, 2, 3, 4, 5}, result)
}

func TestTreeSet_RemoveFunc(t *testing.T) {
	cmp := Cmp[byte]

	ts := TreeSetFrom[byte, Compare[byte]]([]byte{
		'a', 'b', '1', 'c', '2', 'd',
	}, cmp)

	notAlpha := func(c byte) bool {
		return c < 'a' || c > 'z'
	}

	ts.RemoveFunc(notAlpha)

	must.Eq(t, []byte{'a', 'b', 'c', 'd'}, ts.Slice())
}

func TestTreeSet_Contains(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		must.False(t, ts.Contains(42))
	})

	t.Run("exists", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])
		must.Contains[int](t, 1, ts)
		must.Contains[int](t, 2, ts)
		must.Contains[int](t, 3, ts)
		must.Contains[int](t, 4, ts)
		must.Contains[int](t, 5, ts)
	})

	t.Run("absent", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])
		must.NotContains[int](t, 0, ts)
		must.NotContains[int](t, 6, ts)
	})
}

func TestTreeSet_ContainsSlice(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		must.False(t, ts.ContainsSlice([]int{42, 43, 44}))
	})

	t.Run("exists", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])
		must.True(t, ts.ContainsSlice([]int{2, 1, 3}))
		must.True(t, ts.ContainsSlice([]int{5, 4, 3, 2, 1}))
	})

	t.Run("absent", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])
		must.False(t, ts.ContainsSlice([]int{6, 7, 8}))
		must.False(t, ts.ContainsSlice([]int{4, 5, 6}))
	})
}

func TestTreeSet_Subset(t *testing.T) {
	t.Run("empty empty", func(t *testing.T) {
		t1 := NewTreeSet[int, Compare[int]](Cmp[int])
		t2 := NewTreeSet[int, Compare[int]](Cmp[int])
		must.True(t, t1.Subset(t2))
	})

	t.Run("empty full", func(t *testing.T) {
		t1 := NewTreeSet[int, Compare[int]](Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3}, Cmp[int])
		must.False(t, t1.Subset(t2))
	})

	t.Run("full empty", func(t *testing.T) {
		t1 := NewTreeSet[int, Compare[int]](Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3}, Cmp[int])
		must.True(t, t2.Subset(t1))
	})

	t.Run("same", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{2, 1, 3}, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3}, Cmp[int])
		must.True(t, t1.Subset(t2))
		must.True(t, t2.Subset(t1))
	})

	t.Run("subset", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{2, 1, 3}, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{5, 4, 1, 2, 3}, Cmp[int])
		must.False(t, t1.Subset(t2))
	})

	t.Run("superset", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{9, 7, 8, 5, 4, 2, 1, 3}, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{5, 1, 2, 8, 3}, Cmp[int])
		must.True(t, t1.Subset(t2))
	})
	t.Run("diff set", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{6, 7, 8, 9, 10}, Cmp[int])
		must.False(t, t1.Subset(t2))
	})

	t.Run("exhaust s1", func(t *testing.T) {
		s1 := TreeSetFrom[string, Compare[string]]([]string{"a", "b", "c", "d", "e"}, Cmp[string])
		s2 := TreeSetFrom[string, Compare[string]]([]string{"a", "z"}, Cmp[string])
		must.False(t, s1.Subset(s2))
	})
}

func TestTreeSet_Union(t *testing.T) {
	t.Run("empty empty", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]](nil, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]](nil, Cmp[int])
		result := t1.Union(t2)
		must.Empty(t, result)
	})

	t.Run("empty full", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]](nil, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{3, 1, 2}, Cmp[int])
		result := t1.Union(t2)
		must.NotEmpty(t, result)
		must.Eq(t, []int{1, 2, 3}, result.Slice())
	})

	t.Run("full empty", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{2, 3, 1}, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]](nil, Cmp[int])
		result := t1.Union(t2)
		must.NotEmpty(t, result)
		must.Eq(t, []int{1, 2, 3}, result.Slice())
	})

	t.Run("subset", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{2, 3, 1}, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{2}, Cmp[int])
		result := t1.Union(t2)
		must.NotEmpty(t, result)
		must.Eq(t, []int{1, 2, 3}, result.Slice())
	})

	t.Run("superset", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{2, 3, 1}, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{2, 5, 1, 2, 4}, Cmp[int])
		result := t1.Union(t2)
		must.NotEmpty(t, result)
		must.Eq(t, []int{1, 2, 3, 4, 5}, result.Slice())
	})
}

func TestTreeSet_Difference(t *testing.T) {
	t.Run("empty empty", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]](nil, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]](nil, Cmp[int])
		result := t1.Difference(t2)
		must.Empty(t, result)
	})

	t.Run("empty full", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]](nil, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3}, Cmp[int])
		result := t1.Difference(t2)
		must.Empty(t, result)
	})

	t.Run("full empty", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{2, 1, 3}, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]](nil, Cmp[int])
		result := t1.Difference(t2)
		must.NotEmpty(t, result)
		must.Eq(t, []int{1, 2, 3}, result.Slice())
	})

	t.Run("subset", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{3, 2, 4}, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])
		result := t1.Difference(t2)
		must.Empty(t, result)
	})

	t.Run("superset", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{2, 1, 3, 4, 5}, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{1, 2, 5}, Cmp[int])
		result := t1.Difference(t2)
		must.NotEmpty(t, result)
		must.Eq(t, []int{3, 4}, result.Slice())
	})
}

func TestTreeSet_Intersect(t *testing.T) {
	t.Run("empty empty", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]](nil, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]](nil, Cmp[int])
		result := t1.Intersect(t2)
		must.Empty(t, result)
	})

	t.Run("empty full", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]](nil, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3}, Cmp[int])
		result := t1.Intersect(t2)
		must.Empty(t, result)
	})

	t.Run("full empty", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3}, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]](nil, Cmp[int])
		result := t1.Intersect(t2)
		must.Empty(t, result)
	})

	t.Run("overlap", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5, 6}, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{0, 4, 5, 7}, Cmp[int])
		result := t1.Intersect(t2)
		must.NotEmpty(t, result)
		must.Eq(t, []int{4, 5}, result.Slice())
	})
}

func TestTreeSet_Copy(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		t1 := NewTreeSet[int, Compare[int]](Cmp[int])
		c := t1.Copy()
		must.Empty(t, c)
	})

	t.Run("full", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3}, Cmp[int])
		c := t1.Copy()
		must.NotEmpty(t, c)
		must.Eq(t, []int{1, 2, 3}, c.Slice())
	})

	t.Run("modify", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3}, Cmp[int])
		c := t1.Copy()
		c.Insert(4)
		t1.Remove(2)
		must.Eq(t, []int{1, 3}, t1.Slice())
		must.Eq(t, []int{1, 2, 3, 4}, c.Slice())
	})
}

func TestTreeSet_EqualSlice(t *testing.T) {
	t.Run("empty empty", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]](nil, Cmp[int])
		must.True(t, ts.EqualSlice(nil))
	})

	t.Run("empty full", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]](nil, Cmp[int])
		must.False(t, ts.EqualSlice([]int{1, 2, 3}))
	})

	t.Run("matching", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5, 6}, Cmp[int])
		must.True(t, ts.EqualSlice([]int{3, 2, 1, 6, 5, 4}))
	})

	t.Run("different middle", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 5, 6}, Cmp[int])
		must.False(t, ts.EqualSlice([]int{3, 2, 9, 6, 5, 4}))
	})
}

func TestTreeSet_Equal(t *testing.T) {
	t.Run("empty empty", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]](nil, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]](nil, Cmp[int])
		must.Equal(t, t1, t2)
	})

	t.Run("empty full", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]](nil, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3}, Cmp[int])
		must.NotEqual(t, t1, t2)
	})

	t.Run("matching", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5, 6}, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{6, 5, 4, 3, 2, 1}, Cmp[int])
		must.Equal(t, t1, t2)
	})

	t.Run("different min", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4}, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{0, 2, 3, 4}, Cmp[int])
		must.NotEqual(t, t1, t2)
	})

	t.Run("different max", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4}, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{5, 3, 2, 1}, Cmp[int])
		must.NotEqual(t, t1, t2)
	})

	t.Run("different middle", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 5, 6}, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{1, 2, 4, 5, 6}, Cmp[int])
		must.NotEqual(t, t1, t2)
	})
}

func TestTreeSet_TopK(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		result := ts.TopK(5)
		must.Eq(t, []int{}, result)
	})

	t.Run("same size", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{3, 9, 1, 7, 5}, Cmp[int])
		result := ts.TopK(5)
		must.Eq(t, []int{1, 3, 5, 7, 9}, result)
	})

	t.Run("smaller k", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{3, 9, 1, 7, 5}, Cmp[int])
		result := ts.TopK(3)
		must.Eq(t, []int{1, 3, 5}, result)
	})
}

func TestTreeSet_BottomK(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		result := ts.BottomK(5)
		must.Eq(t, []int{}, result)
	})

	t.Run("same size", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{3, 9, 1, 7, 5}, Cmp[int])
		result := ts.BottomK(5)
		must.Eq(t, []int{9, 7, 5, 3, 1}, result)
	})

	t.Run("smaller k", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{3, 9, 1, 7, 5}, Cmp[int])
		result := ts.BottomK(3)
		must.Eq(t, []int{9, 7, 5}, result)
	})
}

func TestTreeSet_FirstBelow(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		_, exists := ts.FirstBelow(5)
		must.False(t, exists)
	})

	t.Run("basic", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{1, 3, 4, 5, 7, 8}, Cmp[int])
		v, exists := ts.FirstBelow(5)
		must.True(t, exists)
		must.Eq(t, 4, v)
	})

	t.Run("many", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		nums := shuffle(ints(100))
		ts.InsertSlice(nums)
		for i := 2; i < 100; i++ {
			v, exists := ts.FirstBelow(i)
			must.True(t, exists)
			must.Eq(t, i-1, v)
		}
	})
}

func TestTreeSet_FirstBelowEqual(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		_, exists := ts.FirstBelowEqual(5)
		must.False(t, exists)
	})

	t.Run("basic", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{1, 3, 4, 5, 7, 8}, Cmp[int])
		v, exists := ts.FirstBelowEqual(5)
		must.True(t, exists)
		must.Eq(t, 5, v)
	})

	t.Run("many", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		nums := shuffle(ints(100))
		ts.InsertSlice(nums)
		for i := 1; i < 100; i++ {
			v, exists := ts.FirstBelowEqual(i)
			must.True(t, exists)
			must.Eq(t, i, v)
		}
	})
}

func TestTreeSet_Below(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{5, 6, 7, 8, 9}, Cmp[int])
		b := ts.Below(5)
		must.Empty(t, b)
	})

	t.Run("basic", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{4, 7, 1, 5, 2, 8, 9, 3}, Cmp[int])
		b := ts.Below(5)
		result := b.Slice()
		must.Eq(t, []int{1, 2, 3, 4}, result)
	})

	t.Run("many", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		nums := shuffle(ints(100))
		ts.InsertSlice(nums)
		for i := 2; i < 100; i++ {
			below := ts.Below(i)
			must.Size(t, i-1, below)
			must.Min(t, 1, below)
			must.Max(t, i-1, below)
		}
	})
}

func TestTreeSet_BelowEqual(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{5, 6, 7, 8, 9}, Cmp[int])
		b := ts.BelowEqual(4)
		must.Empty(t, b)
	})

	t.Run("basic", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{4, 7, 1, 5, 2, 8, 9, 3}, Cmp[int])
		b := ts.BelowEqual(5)
		result := b.Slice()
		must.Eq(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("many", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		nums := shuffle(ints(100))
		ts.InsertSlice(nums)
		for i := 1; i < 100; i++ {
			below := ts.BelowEqual(i)
			must.Size(t, i, below)
			must.Min(t, 1, below)
			must.Max(t, i, below)
		}
	})
}

func TestTreeSet_FirstAbove(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{2, 1, 3, 5, 4}, Cmp[int])
		_, exists := ts.FirstAbove(5)
		must.False(t, exists)
	})

	t.Run("basic", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{2, 1, 4, 6, 5, 7, 8}, Cmp[int])
		v, exists := ts.FirstAbove(5)
		must.True(t, exists)
		must.Eq(t, 6, v)
	})

	t.Run("many", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		nums := shuffle(ints(100))
		ts.InsertSlice(nums)
		for i := 1; i < 100; i++ {
			v, exists := ts.FirstAbove(i)
			must.True(t, exists)
			must.Eq(t, i+1, v)
		}
	})
}

func TestTreeSet_FirstAboveEqual(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{2, 1, 3, 4}, Cmp[int])
		_, exists := ts.FirstAboveEqual(5)
		must.False(t, exists)
	})

	t.Run("basic", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{2, 1, 4, 6, 5, 7, 8}, Cmp[int])
		v, exists := ts.FirstAboveEqual(5)
		must.True(t, exists)
		must.Eq(t, 5, v)
	})

	t.Run("many", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		nums := shuffle(ints(100))
		ts.InsertSlice(nums)
		for i := 1; i < 100; i++ {
			v, exists := ts.FirstAboveEqual(i)
			must.True(t, exists)
			must.Eq(t, i, v)
		}
	})
}

func TestTreeSet_Above(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{5, 6, 7, 8, 9}, Cmp[int])
		b := ts.Above(9)
		must.Empty(t, b)
	})

	t.Run("basic", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{4, 7, 1, 5, 2, 8, 9, 3}, Cmp[int])
		b := ts.Above(5)
		result := b.Slice()
		must.Eq(t, []int{7, 8, 9}, result)
	})

	t.Run("many", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		nums := shuffle(ints(100))
		ts.InsertSlice(nums)
		for i := 1; i < 100; i++ {
			above := ts.Above(i)
			must.Size(t, 100-i, above)
			must.Min(t, i+1, above)
			must.Max(t, 100, above)
		}
	})
}

func TestTreeSet_AboveEqual(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{5, 6, 7, 8, 9}, Cmp[int])
		b := ts.AboveEqual(10)
		must.Empty(t, b)
	})

	t.Run("basic", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{4, 7, 1, 5, 2, 8, 9, 3}, Cmp[int])
		b := ts.AboveEqual(5)
		result := b.Slice()
		must.Eq(t, []int{5, 7, 8, 9}, result)
	})

	t.Run("many", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		nums := shuffle(ints(100))
		ts.InsertSlice(nums)
		for i := 1; i < 100; i++ {
			above := ts.AboveEqual(i)
			must.Size(t, 100-i+1, above)
			must.Min(t, i, above)
			must.Max(t, 100, above)
		}
	})
}

func TestTreeSet_Slice(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		result := ts.Slice()
		must.Eq(t, []int{}, result)
	})

	t.Run("full", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{4, 2, 6, 1}, Cmp[int])
		result := ts.Slice()
		must.Eq(t, []int{1, 2, 4, 6}, result)
	})
}

func TestTreeSet_String(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		result := ts.String()
		must.Eq(t, "[]", result)
	})

	t.Run("full", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{4, 2, 6, 1}, Cmp[int])
		result := ts.String()
		must.Eq(t, "[1 2 4 6]", result)
	})
}

func TestTreeSet_StringFunc(t *testing.T) {
	f := func(i int) string { return fmt.Sprintf("%02d", i) }
	t.Run("empty", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		result := ts.StringFunc(f)
		must.Eq(t, "[]", result)
	})

	t.Run("full", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{4, 2, 6, 1}, Cmp[int])
		result := ts.StringFunc(f)
		must.Eq(t, "[01 02 04 06]", result)
	})
}

// create a colorful representation of the element in node
func (n *node[T]) String() string {
	if n.red() {
		return fmt.Sprintf("\033[1;31m%v\033[0m", n.element)
	}
	return fmt.Sprintf("%v", n.element)
}

// output creates a colorful string representation of s
func (s *TreeSet[T, C]) output(prefix, cprefix string, n *node[T], sb *strings.Builder) {
	if n == nil {
		return
	}

	sb.WriteString(prefix)
	sb.WriteString(n.String())
	sb.WriteString("\n")

	if n.right != nil && n.left != nil {
		s.output(cprefix+"├── ", cprefix+"│   ", n.right, sb)
	} else if n.right != nil {
		s.output(cprefix+"└── ", cprefix+"    ", n.right, sb)
	}
	if n.left != nil {
		s.output(cprefix+"└── ", cprefix+"    ", n.left, sb)
	}
	if n.left == nil && n.right == nil {
		return
	}
}

// dump the output of s along with the slice string
func (s *TreeSet[T, C]) dump() string {
	var sb strings.Builder
	sb.WriteString("\ntree:\n")
	s.output("", "", s.root, &sb)
	sb.WriteString("string:")
	sb.WriteString(s.String())
	return sb.String()
}

// invariants makes basic assertions about tree
func invariants[T any, C Compare[T]](t *testing.T, tree *TreeSet[T, C], cmp C) {
	// assert Slice elements are ascending
	slice := tree.Slice()
	must.AscendingCmp(t, slice, cmp)

	// assert size of tree
	size := tree.Size()
	must.Eq(t, size, len(slice), must.Sprint("tree is wrong size"))

	if size == 0 {
		return
	}

	// assert slice[0] is the minimum
	must.Min(t, slice[0], tree)

	// assert slice[len(slice)-1] is the maximum
	must.Max(t, slice[len(slice)-1], tree)
}

// ints will create a []int from 1 to n
func ints(n int) []int {
	s := make([]int, n)
	for i := 0; i < n; i++ {
		s[i] = i + 1
	}
	return s
}

// create a copy of s and shuffle
func shuffle(s []int) []int {
	c := make([]int, len(s))
	copy(c, s)

	n := len(c)
	for i := 0; i < n; i++ {
		swp := rand.Int31n(int32(n))
		c[i], c[swp] = c[swp], c[i]
	}
	return c
}

func TestTreeSet_infix(t *testing.T) {
	ts := TreeSetFrom[int, Compare[int]]([]int{4, 7, 1, 5, 2, 8, 9, 3, 11, 13}, Cmp[int])
	isOdd := func(n *node[int]) bool {
		return n.element%2 == 1
	}
	odds := make([]int, 0, 5)
	ts.infix(func(n *node[int]) bool {
		if n.element > 8 {
			return false
		}
		if isOdd(n) {
			odds = append(odds, n.element)
		}

		return true
	}, ts.root)
	must.Eq(t, []int{1, 3, 5, 7}, odds)
}
func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestTreeSet_iterate(t *testing.T) {
	s := TreeSetFrom[int, Compare[int]]([]int{4, 7, 1, 5, 2, 8, 9, 3, 11}, Cmp[int])
	ctx, cl := context.WithCancel(context.Background())
	defer cl()
	ret := make([]int, 0, 9)
	ch := s.iterate(ctx)
	for n := range ch {
		if n.element > 3 {
			break
		}
		ret = append(ret, n.element)
	}
	must.Eq(t, []int{1, 2, 3}, ret)
}
