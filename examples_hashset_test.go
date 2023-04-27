package set

import (
	"fmt"
)

type person struct {
	name string
	id   int
}

func (p *person) Hash() string {
	return fmt.Sprintf("%s:%d", p.name, p.id)
}

func (p *person) String() string {
	return p.name
}

func ExampleHashSet_Insert() {
	s := NewHashSet[*person, string](10)
	s.Insert(&person{name: "dave", id: 108})
	s.Insert(&person{name: "armon", id: 101})
	s.Insert(&person{name: "mitchell", id: 100})

	fmt.Println(s)
	// Output:
	// [armon dave mitchell]
}

// InsertSlice

// InsertSet

// Remove

// RemoveSlice

// RemoveSet

// RemoveFunc

func ExampleHashSet_Contains() {
	anna := &person{name: "anna", id: 94}
	bill := &person{name: "bill", id: 50}
	carl := &person{name: "carl", id: 10}
	dave := &person{name: "dave", id: 32}
	s := HashSetFrom[*person, string]([]*person{anna, bill, carl})

	fmt.Println(s.Contains(anna))
	fmt.Println(s.Contains(dave))

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
