// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package set

import (
	"fmt"
	"sort"
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
	s.Insert(&person{name: "armon", id: 101})

	fmt.Println(s)

	// Output:
	// [armon dave mitchell]
}

func ExampleHashSet_InsertSlice() {
	s := NewHashSet[*person, string](10)
	s.InsertSlice([]*person{
		{name: "dave", id: 108},
		{name: "mitchell", id: 100},
		{name: "dave", id: 108},
		{name: "armon", id: 101},
	})

	fmt.Println(s)

	// Output:
	// [armon dave mitchell]
}

func ExampleHashSet_InsertSet() {
	anna := &person{name: "anna", id: 94}
	bill := &person{name: "bill", id: 50}
	carl := &person{name: "carl", id: 10}
	dave := &person{name: "dave", id: 32}
	s1 := HashSetFrom[*person, string]([]*person{anna, carl})
	s2 := HashSetFrom[*person, string]([]*person{carl, dave, bill})
	s2.InsertSet(s1)

	fmt.Println(s1)
	fmt.Println(s2)

	// Output:
	// [anna carl]
	// [anna bill carl dave]
}

func ExampleHashSet_Remove() {
	anna := &person{name: "anna", id: 94}
	bill := &person{name: "bill", id: 50}
	carl := &person{name: "carl", id: 10}
	dave := &person{name: "dave", id: 32}
	s := HashSetFrom[*person, string]([]*person{anna, carl, dave, bill})

	fmt.Println(s)

	s.Remove(carl)

	fmt.Println(s)

	// Output:
	// [anna bill carl dave]
	// [anna bill dave]
}

func ExampleHashSet_RemoveSlice() {
	anna := &person{name: "anna", id: 94}
	bill := &person{name: "bill", id: 50}
	carl := &person{name: "carl", id: 10}
	dave := &person{name: "dave", id: 32}
	s := HashSetFrom[*person, string]([]*person{anna, carl, dave, bill})

	fmt.Println(s)

	s.RemoveSlice([]*person{anna, carl})

	fmt.Println(s)

	// Output:
	// [anna bill carl dave]
	// [bill dave]
}

func ExampleHashSet_RemoveSet() {
	anna := &person{name: "anna", id: 94}
	bill := &person{name: "bill", id: 50}
	carl := &person{name: "carl", id: 10}
	dave := &person{name: "dave", id: 32}
	s := HashSetFrom[*person, string]([]*person{carl, dave, bill})
	r := HashSetFrom[*person, string]([]*person{anna, carl})

	fmt.Println(s)

	s.RemoveSet(r)

	fmt.Println(s)

	// Output:
	// [bill carl dave]
	// [bill dave]
}

func ExampleHashSet_RemoveFunc() {
	anna := &person{name: "anna", id: 94}
	bill := &person{name: "bill", id: 50}
	carl := &person{name: "carl", id: 10}
	dave := &person{name: "dave", id: 32}
	s := HashSetFrom[*person, string]([]*person{anna, bill, carl, dave})

	idAbove50 := func(p *person) bool {
		return p.id >= 50
	}

	fmt.Println(s)

	s.RemoveFunc(idAbove50)

	fmt.Println(s)

	// Output:
	// [anna bill carl dave]
	// [carl dave]
}

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

func ExampleHashSet_ContainsSlice() {
	anna := &person{name: "anna", id: 94}
	bill := &person{name: "bill", id: 50}
	carl := &person{name: "carl", id: 10}
	dave := &person{name: "dave", id: 32}
	s := HashSetFrom[*person, string]([]*person{anna, bill, carl})

	fmt.Println(s.ContainsSlice([]*person{anna, bill}))
	fmt.Println(s.ContainsSlice([]*person{anna, bill, carl}))
	fmt.Println(s.ContainsSlice([]*person{carl, dave}))

	// Output:
	// false
	// true
	// false
}

func ExampleHashSet_Subset() {
	anna := &person{name: "anna", id: 94}
	bill := &person{name: "bill", id: 50}
	carl := &person{name: "carl", id: 10}
	dave := &person{name: "dave", id: 32}
	s1 := HashSetFrom[*person, string]([]*person{anna, bill, carl})
	s2 := HashSetFrom[*person, string]([]*person{anna, bill})
	s3 := HashSetFrom[*person, string]([]*person{bill, carl, dave})

	fmt.Println(s1.Subset(s2))
	fmt.Println(s1.Subset(s3))

	// Output:
	// true
	// false
}

func ExampleHashSet_Size() {
	anna := &person{name: "anna", id: 94}
	bill := &person{name: "bill", id: 50}
	carl := &person{name: "carl", id: 10}
	s := HashSetFrom[*person, string]([]*person{anna, bill, carl})

	fmt.Println(s.Size())

	// Output:
	// 3
}

func ExampleHashSet_Empty() {
	s := NewHashSet[*person, string](0)

	fmt.Println(s.Empty())

	// Output:
	// true
}

func ExampleHashSet_Union() {
	anna := &person{name: "anna", id: 94}
	bill := &person{name: "bill", id: 50}
	carl := &person{name: "carl", id: 10}
	dave := &person{name: "dave", id: 32}
	s1 := HashSetFrom[*person, string]([]*person{anna, bill, carl})
	s2 := HashSetFrom[*person, string]([]*person{anna, bill, dave})
	union := s1.Union(s2)

	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(union)

	// Output:
	// [anna bill carl]
	// [anna bill dave]
	// [anna bill carl dave]
}

func ExampleHashSet_Difference() {
	anna := &person{name: "anna", id: 94}
	bill := &person{name: "bill", id: 50}
	carl := &person{name: "carl", id: 10}
	dave := &person{name: "dave", id: 32}
	s1 := HashSetFrom[*person, string]([]*person{anna, bill, carl})
	s2 := HashSetFrom[*person, string]([]*person{anna, bill, dave})
	difference := s1.Difference(s2)

	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(difference)

	// Output:
	// [anna bill carl]
	// [anna bill dave]
	// [carl]
}

func ExampleHashSet_Intersect() {
	anna := &person{name: "anna", id: 94}
	bill := &person{name: "bill", id: 50}
	carl := &person{name: "carl", id: 10}
	dave := &person{name: "dave", id: 32}
	s1 := HashSetFrom[*person, string]([]*person{anna, bill, carl})
	s2 := HashSetFrom[*person, string]([]*person{anna, bill, dave})
	intersect := s1.Intersect(s2)

	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(intersect)

	// Output:
	// [anna bill carl]
	// [anna bill dave]
	// [anna bill]
}

func ExampleHashSet_Equal() {
	anna := &person{name: "anna", id: 94}
	bill := &person{name: "bill", id: 50}
	carl := &person{name: "carl", id: 10}
	dave := &person{name: "dave", id: 32}

	s1 := HashSetFrom[*person, string]([]*person{anna, bill, carl})
	s2 := HashSetFrom[*person, string]([]*person{anna, bill, carl})
	s3 := HashSetFrom[*person, string]([]*person{anna, bill, dave})

	fmt.Println(s1.Equal(s2))
	fmt.Println(s1.Equal(s3))

	// Output:
	// true
	// false
}

func ExampleHashSet_Copy() {
	anna := &person{name: "anna", id: 94}
	bill := &person{name: "bill", id: 50}
	carl := &person{name: "carl", id: 10}
	s := HashSetFrom[*person, string]([]*person{anna, bill, carl})
	c := s.Copy()

	fmt.Println(c)

	// Output:
	// [anna bill carl]
}

func ExampleHashSet_Slice() {
	anna := &person{name: "anna", id: 94}
	bill := &person{name: "bill", id: 50}
	carl := &person{name: "carl", id: 10}
	s := HashSetFrom[*person, string]([]*person{anna, bill, carl})

	slice := s.Slice()
	sort.Slice(slice, func(a, b int) bool {
		return slice[a].id < slice[b].id
	})

	fmt.Println(slice)

	// Output:
	// [carl bill anna]
}

func ExampleHashSet_String() {
	anna := &person{name: "anna", id: 94}
	bill := &person{name: "bill", id: 50}
	carl := &person{name: "carl", id: 10}
	s := HashSetFrom[*person, string]([]*person{anna, bill, carl})

	fmt.Println(s.String())

	// Output:
	// [anna bill carl]
}
