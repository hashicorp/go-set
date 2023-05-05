// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package set

import (
	"fmt"
	"sort"
)

func ExampleSet_Insert() {
	s := New[int](10)
	s.Insert(1)
	s.Insert(1)
	s.Insert(2)
	s.Insert(3)
	s.Insert(2)

	fmt.Println(s)

	// Output:
	// [1 2 3]
}

func ExampleSet_InsertSlice() {
	s := New[int](10)
	s.InsertSlice([]int{1, 1, 2, 3, 2})

	fmt.Println(s)

	// Output:
	// [1 2 3]
}

func ExampleSet_InsertSet() {
	s := New[int](10)
	s.InsertSet(From([]int{1, 1, 2, 3, 2}))

	fmt.Println(s)

	// Output:
	// [1 2 3]
}

func ExampleSet_Remove() {
	s := New[int](10)
	s.InsertSlice([]int{1, 1, 2, 3, 2})
	s.Remove(2)

	fmt.Println(s)

	// Output:
	// [1 3]
}

func ExampleSet_RemoveSlice() {
	s := New[int](10)
	s.InsertSlice([]int{1, 1, 2, 3, 2})
	s.RemoveSlice([]int{2, 3})

	fmt.Println(s)

	// Output:
	// [1]
}

func ExampleSet_RemoveSet() {
	s := New[int](10)
	s.InsertSlice([]int{1, 1, 2, 3, 2})
	s.RemoveSet(From([]int{2, 3}))

	fmt.Println(s)

	// Output:
	// [1]
}

func ExampleSet_RemoveFunc() {
	s := From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0})
	even := func(i int) bool {
		return i%2 == 0
	}
	s.RemoveFunc(even)

	fmt.Println(s)

	// Output:
	// [1 3 5 7 9]
}

func ExampleSet_Contains() {
	s := From([]string{"red", "green", "blue"})

	fmt.Println(s.Contains("red"))
	fmt.Println(s.Contains("orange"))

	// Output:
	// true
	// false
}

func ExampleSet_ContainsAll() {
	s := From([]string{"red", "green", "blue"})

	fmt.Println(s.ContainsAll([]string{"red", "blue"}))
	fmt.Println(s.ContainsAll([]string{"red", "orange"}))

	// Output:
	// true
	// false
}

func ExampleSet_ContainsSlice() {
	s := From([]string{"red", "green", "blue"})

	fmt.Println(s.ContainsSlice([]string{"red", "blue"}))
	fmt.Println(s.ContainsSlice([]string{"red", "blue", "orange"}))
	fmt.Println(s.ContainsSlice([]string{"red", "blue", "green"}))

	// Output:
	// false
	// false
	// true
}

func ExampleSet_Subset() {
	t1 := From([]string{"red", "green", "blue"})
	t2 := From([]string{"red", "blue"})
	t3 := From([]string{"red", "orange"})

	fmt.Println(t1.Subset(t2))
	fmt.Println(t1.Subset(t3))

	// Output:
	// true
	// false
}

func ExampleSet_Size() {
	s := From([]string{"red", "green", "blue"})

	fmt.Println(s.Size())

	// Output:
	// 3
}

func ExampleSet_Empty() {
	s := New[string](10)

	fmt.Println(s.Empty())

	// Output:
	// true
}

func ExampleSet_Union() {
	t1 := From([]string{"red", "green", "blue"})
	t2 := From([]string{"red", "blue"})
	t3 := From([]string{"red", "orange"})

	fmt.Println(t1.Union(t2))
	fmt.Println(t1.Union(t3))

	// Output:
	// [blue green red]
	// [blue green orange red]
}

func ExampleSet_Difference() {
	t1 := From([]string{"red", "green", "blue"})
	t2 := From([]string{"red", "blue"})
	t3 := From([]string{"red", "orange"})

	fmt.Println(t1.Difference(t2))
	fmt.Println(t1.Difference(t3))

	// Output:
	// [green]
	// [blue green]
}

func ExampleSet_Intersect() {
	t1 := From([]string{"red", "green", "blue"})
	t2 := From([]string{"red", "blue"})
	t3 := From([]string{"red", "orange"})
	t4 := From([]string{"yellow"})

	fmt.Println(t1.Intersect(t2))
	fmt.Println(t1.Intersect(t3))
	fmt.Println(t1.Intersect(t4))

	// Output:
	// [blue red]
	// [red]
	// []
}

func ExampleSet_Equal() {
	t1 := From([]string{"red", "green", "blue"})
	t2 := From([]string{"red", "blue"})
	t3 := From([]string{"red", "green", "yellow"})
	t4 := From([]string{"red", "green", "blue"})

	fmt.Println(t1.Equal(t2))
	fmt.Println(t1.Equal(t3))
	fmt.Println(t1.Equal(t4))

	// Output:
	// false
	// false
	// true
}

func ExampleSet_Copy() {
	s := From([]string{"red", "green", "blue"})
	t := s.Copy()

	fmt.Println(t)

	// Output:
	// [blue green red]
}

func ExampleSet_Slice() {
	s := From([]string{"red", "green", "blue"})
	t := s.Slice()

	sort.Strings(t)
	fmt.Println(t)

	// Output:
	// [blue green red]
}

func ExampleSet_String() {
	s := From([]string{"red", "green", "blue"})

	fmt.Println(s.String())

	// Output:
	// [blue green red]
}
