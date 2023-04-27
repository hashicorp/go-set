package set

import (
	"fmt"
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

// InsertSlice

func ExampleSet_InsertSlice() {
    
    s := set.New()
    s.InsertSlice([]int{1, 2, 3})
	
    fmt.Println(s)
    
    // Output:
   // [1 2 3]
	
}

// InsertSet

// Remove

// RemoveSlice

// RemoveSet

// RemoveFunc

func ExampleSet_Contains() {
	s := From([]string{"red", "green", "blue"})

	fmt.Println(s.Contains("red"))
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
