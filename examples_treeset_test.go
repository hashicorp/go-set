// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package set

import (
	"fmt"
)

func ExampleCompare_contestant() {
	type contestant struct {
		name  string
		score int
	}

	compare := func(a, b contestant) int {
		return a.score - b.score
	}

	s := NewTreeSet[contestant, Compare[contestant]](compare)
	s.Insert(contestant{name: "alice", score: 80})
	s.Insert(contestant{name: "dave", score: 90})
	s.Insert(contestant{name: "bob", score: 70})

	fmt.Println(s)

	// Output:
	// [{bob 70} {alice 80} {dave 90}]
}

func ExampleCmp_strings() {
	s := NewTreeSet[string, Compare[string]](Cmp[string])
	s.Insert("red")
	s.Insert("green")
	s.Insert("blue")

	fmt.Println(s)
	fmt.Println("min:", s.Min())
	fmt.Println("max:", s.Max())

	// Output:
	// [blue green red]
	// min: blue
	// max: red
}

func ExampleCmp_ints() {
	s := NewTreeSet[int, Compare[int]](Cmp[int])
	s.Insert(50)
	s.Insert(42)
	s.Insert(100)

	fmt.Println(s)
	fmt.Println("min:", s.Min())
	fmt.Println("max:", s.Max())

	// Output:
	// [42 50 100]
	// min: 42
	// max: 100
}

func ExampleTreeSet_Insert() {
	s := TreeSetFrom[string, Compare[string]]([]string{}, Cmp[string])

	fmt.Println(s)

	s.Insert("red")
	s.Insert("green")
	s.Insert("blue")

	fmt.Println(s)

	// []
	// [blue green red]
}

func ExampleTreeSet_InsertSlice() {
	s := TreeSetFrom[string, Compare[string]]([]string{}, Cmp[string])

	fmt.Println(s)

	s.InsertSlice([]string{"red", "green", "blue"})

	fmt.Println(s)

	// []
	// [blue green red]
}

// InsertSet

func ExampleTreeSet_Remove() {
	s := TreeSetFrom[string, Compare[string]]([]string{"red", "green", "blue"}, Cmp[string])

	fmt.Println(s)

	fmt.Println(s.Remove("green"))
	fmt.Println(s.Remove("orange"))

	fmt.Println(s)

	// Output:
	// [blue green red]
	// true
	// false
	// [blue red]
}

func ExampleTreeSet_RemoveSlice() {
	s := TreeSetFrom[string, Compare[string]]([]string{"red", "green", "blue"}, Cmp[string])

	fmt.Println(s)

	fmt.Println(s.RemoveSlice([]string{"red", "blue"}))
	fmt.Println(s.RemoveSlice([]string{"orange", "white"}))

	fmt.Println(s)

	// Output:
	// [blue green red]
	// true
	// false
	// [green]
}

// RemoveSet

// RemoveFunc

func ExampleTreeSet_Contains() {
	s := TreeSetFrom[string, Compare[string]]([]string{"red", "green", "blue"}, Cmp[string])

	fmt.Println(s.Contains("green"))
	fmt.Println(s.Contains("orange"))

	// Output:
	// true
	// false
}

// ContainsAll

func ExampleTreeSet_ContainsSlice() {
	s := TreeSetFrom[string, Compare[string]]([]string{"red", "green", "blue"}, Cmp[string])

	fmt.Println(s.ContainsSlice([]string{"red", "green"}))
	fmt.Println(s.ContainsSlice([]string{"red", "orange"}))

	// Output:
	// true
	// false
}

// Subset

func ExampleTreeSet_Size() {
	s := TreeSetFrom[string, Compare[string]]([]string{"red", "green", "blue"}, Cmp[string])

	fmt.Println(s.Size())

	// Output:
	// 3
}

func ExampleTreeSet_Empty() {
	s := TreeSetFrom[string, Compare[string]]([]string{}, Cmp[string])

	fmt.Println(s.Empty())

	s.InsertSlice([]string{"red", "green", "blue"})

	fmt.Println(s.Empty())

	// Output:
	// true
	// false
}

func ExampleTreeSet_Union() {
	s := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])
	t := TreeSetFrom[int, Compare[int]]([]int{5, 4, 3, 2, 1}, Cmp[int])
	f := TreeSetFrom[int, Compare[int]]([]int{1, 3, 5, 7, 9}, Cmp[int])

	fmt.Println(s.Union(t))
	fmt.Println(s.Union(f))

	// Output:
	// [1 2 3 4 5]
	// [1 2 3 4 5 7 9]
}

func ExampleTreeSet_Difference() {
	s := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])
	t := TreeSetFrom[int, Compare[int]]([]int{5, 4, 3, 2, 1}, Cmp[int])
	f := TreeSetFrom[int, Compare[int]]([]int{1, 3, 5, 7, 9}, Cmp[int])

	fmt.Println(s.Difference(t))
	fmt.Println(s.Difference(f))

	// Output:
	// []
	// [2 4]
}

func ExampleTreeSet_Intersect() {
	s := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])
	t := TreeSetFrom[int, Compare[int]]([]int{5, 4, 3, 2, 1}, Cmp[int])
	f := TreeSetFrom[int, Compare[int]]([]int{1, 3, 5, 7, 9}, Cmp[int])

	fmt.Println(s.Intersect(t))
	fmt.Println(s.Intersect(f))

	// Output:
	// [1 2 3 4 5]
	// [1 3 5]
}

func ExampleTreeSet_Equal() {
	s := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])
	t := TreeSetFrom[int, Compare[int]]([]int{5, 4, 3, 2, 1}, Cmp[int])
	f := TreeSetFrom[int, Compare[int]]([]int{1, 3, 5, 7, 9}, Cmp[int])

	fmt.Println(s.Equal(t))
	fmt.Println(s.Equal(f))

	// Output:
	// true
	// false
}

// Copy

func ExampleTreeSet_Slice() {
	s := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])
	slice := s.Slice()

	fmt.Println(slice)
	fmt.Println(len(slice))

	// Output:
	// [1 2 3 4 5]
	// 5
}

func ExampleTreeSet_String() {
	s := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])

	fmt.Println(s.String() == "[1 2 3 4 5]")

	// Output:
	// true
}

func ExampleTreeSet_Min() {
	s := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])
	r := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, func(a int, b int) int {
		return b - a
	})

	fmt.Println("asc:", s.Min())
	fmt.Println("desc:", r.Min())

	// Output:
	// asc: 1
	// desc: 5
}

func ExampleTreeSet_Max() {
	s := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])
	r := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, func(a int, b int) int {
		return b - a
	})

	fmt.Println("asc:", s.Max())
	fmt.Println("desc:", r.Max())

	// Output:
	// asc: 5
	// desc: 1
}

func ExampleTreeSet_TopK() {
	s := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])

	fmt.Println(s.TopK(0))
	fmt.Println(s.TopK(1))
	fmt.Println(s.TopK(3))
	fmt.Println(s.TopK(5))

	// Output:
	// []
	// [1]
	// [1 2 3]
	// [1 2 3 4 5]
}

func ExampleTreeSet_BottomK() {
	s := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])

	fmt.Println(s.BottomK(0))
	fmt.Println(s.BottomK(1))
	fmt.Println(s.BottomK(3))
	fmt.Println(s.BottomK(5))

	// Output:
	// []
	// [5]
	// [5 4 3]
	// [5 4 3 2 1]
}

func ExampleTreeSet_FirstAbove() {
	s := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])

	fmt.Println(s.FirstAbove(3))
	fmt.Println(s.FirstAbove(5))
	fmt.Println(s.FirstAbove(10))

	// Output:
	// 4 true
	// 0 false
	// 0 false
}

func ExampleTreeSet_FirstAboveEqual() {
	s := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])

	fmt.Println(s.FirstAboveEqual(3))
	fmt.Println(s.FirstAboveEqual(5))
	fmt.Println(s.FirstAboveEqual(10))

	// Output:
	// 3 true
	// 5 true
	// 0 false
}

func ExampleTreeSet_Above() {
	s := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])

	fmt.Println(s.Above(3))
	fmt.Println(s.Above(5))
	fmt.Println(s.Above(10))

	// Output:
	// [4 5]
	// []
	// []
}

func ExampleTreeSet_AboveEqual() {
	s := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])

	fmt.Println(s.AboveEqual(3))
	fmt.Println(s.AboveEqual(5))
	fmt.Println(s.AboveEqual(10))

	// Output:
	// [3 4 5]
	// [5]
	// []
}

func ExampleTreeSet_FirstBelow() {
	s := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])

	fmt.Println(s.FirstBelow(1))
	fmt.Println(s.FirstBelow(3))
	fmt.Println(s.FirstBelow(10))

	// Output:
	// 0 false
	// 2 true
	// 5 true
}

func ExampleTreeSet_FirstBelowEqual() {
	s := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])

	fmt.Println(s.FirstBelowEqual(1))
	fmt.Println(s.FirstBelowEqual(3))
	fmt.Println(s.FirstBelowEqual(10))

	// Output:
	// 1 true
	// 3 true
	// 5 true
}

func ExampleTreeSet_Below() {
	s := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])

	fmt.Println(s.Below(1))
	fmt.Println(s.Below(3))
	fmt.Println(s.Below(10))

	// Output:
	// []
	// [1 2]
	// [1 2 3 4 5]
}

func ExampleTreeSet_BelowEqual() {
	s := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])

	fmt.Println(s.BelowEqual(1))
	fmt.Println(s.BelowEqual(3))
	fmt.Println(s.BelowEqual(10))

	// Output:
	// [1]
	// [1 2 3]
	// [1 2 3 4 5]
}
