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

// Insert

// InsertSlice

// InsertSet

// Remove

// RemoveSlice

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

// ContainsSlice

// Subset

// Size

// Empty

// Union

// Difference

// Intersect

// Equal

// Copy

// Slice

// String

// Min

// Max

// TopK

// BottomK

// FirstAbove

// FirstAboveEqual

// Above

// AboveEqual

// FirstBelow

// FirstBelowEqual

// Below

// BelowEqual
