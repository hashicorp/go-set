package set

import (
	"math/rand"
	"sort"
	"testing"
)

type test struct {
	size int
	name string
}

var cases []test = []test{
	{size: 10, name: "10"},
	{size: 1_000, name: "1000"},
	{size: 100_000, name: "100000"},
	{size: 1_000_000, name: "1000000"},
}

func random[I ~int](n int) []I {
	result := make([]I, n)
	for i := 0; i < n; i++ {
		result[i] = I(rand.Int())
	}
	return result
}

func unsort[I ~int](s []I) {
	s[0], s[len(s)-1] = s[len(s)-1], s[0]
}

type hashint int

func (i hashint) Hash() int {
	return int(i)
}

func BenchmarkSlice_Insert(b *testing.B) {
	for _, tc := range cases {
		s := random[int](tc.size)
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s = append(s, i)
			}
		})
	}
}

func BenchmarkSet_Insert(b *testing.B) {
	for _, tc := range cases {
		s := From(random[int](tc.size))
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s.Insert(i)
			}
		})
	}
}

func BenchmarkHashSet_Insert(b *testing.B) {
	for _, tc := range cases {
		hs := HashSetFrom[hashint, int](random[hashint](tc.size))
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				hs.Insert(hashint(i))
			}
		})
	}
}

func BenchmarkTreeSet_Insert(b *testing.B) {
	for _, tc := range cases {
		ts := TreeSetFrom[int, Compare[int]](random[int](tc.size), Cmp[int])
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ts.Insert(i)
			}
		})
	}
}

func BenchmarkSlice_Minimum(b *testing.B) {
	for _, tc := range cases {
		slice := random[int](tc.size)
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				sort.Ints(slice)
				_ = slice[0]
				unsort(slice)
			}
		})
	}
}

func BenchmarkSet_Minimum(b *testing.B) {
	for _, tc := range cases {
		s := From(random[int](tc.size))
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				values := s.Slice()
				sort.Ints(values)
				_ = values[0]
			}
		})
	}
}

func BenchmarkHashSet_Minimum(b *testing.B) {
	for _, tc := range cases {
		hs := HashSetFrom[hashint, int](random[hashint](tc.size))
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				values := hs.Slice()
				sort.Slice(values, func(a, b int) bool { return values[a] < values[b] })
				_ = values[0]
				unsort(values)
			}
		})
	}
}

func BenchmarkTreeSet_Minimum(b *testing.B) {
	for _, tc := range cases {
		ts := TreeSetFrom[int, Compare[int]](random[int](tc.size), Cmp[int])
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = ts.Min()
			}
		})
	}
}

func BenchmarkSlice_Contains(b *testing.B) {
	contains := func(s []int, target int) bool {
		for i := 0; i < len(s); i++ {
			if s[i] == target {
				return true
			}
		}
		return false
	}

	for _, tc := range cases {
		slice := random[int](tc.size)
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = contains(slice, i)
			}
		})
	}
}

func BenchmarkSet_Contains(b *testing.B) {
	for _, tc := range cases {
		s := From(random[int](tc.size))
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = s.Contains(i)
			}
		})
	}
}

func BenchmarkHashSet_Contains(b *testing.B) {
	for _, tc := range cases {
		hs := HashSetFrom[hashint, int](random[hashint](tc.size))
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = hs.Contains(hashint(i))
			}
		})
	}
}

func BenchmarkTreeSet_Contains(b *testing.B) {
	for _, tc := range cases {
		ts := TreeSetFrom[int, Compare[int]](random[int](tc.size), Cmp[int])
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = ts.Contains(i)
			}
		})
	}
}
