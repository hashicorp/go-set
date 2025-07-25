// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package set

import (
	"fmt"
	"sort"
)

type person struct {
	Name string
	ID   int
}

func (p *person) Hash() string {
	return fmt.Sprintf("%s:%d", p.Name, p.ID)
}

func (p *person) String() string {
	return p.Name
}

func ExampleHashSet_Insert() {
	s := NewHashSet[*person, string](10)
	s.Insert(&person{Name: "dave", ID: 108})
	s.Insert(&person{Name: "armon", ID: 101})
	s.Insert(&person{Name: "mitchell", ID: 100})
	s.Insert(&person{Name: "armon", ID: 101})

	fmt.Println(s)

	// Output:
	// [armon dave mitchell]
}

func ExampleHashSet_InsertSlice() {
	s := NewHashSet[*person, string](10)
	s.InsertSlice([]*person{
		{Name: "dave", ID: 108},
		{Name: "mitchell", ID: 100},
		{Name: "dave", ID: 108},
		{Name: "armon", ID: 101},
	})

	fmt.Println(s)

	// Output:
	// [armon dave mitchell]
}

func ExampleHashSet_InsertSet() {
	anna := &person{Name: "anna", ID: 94}
	bill := &person{Name: "bill", ID: 50}
	carl := &person{Name: "carl", ID: 10}
	dave := &person{Name: "dave", ID: 32}
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
	anna := &person{Name: "anna", ID: 94}
	bill := &person{Name: "bill", ID: 50}
	carl := &person{Name: "carl", ID: 10}
	dave := &person{Name: "dave", ID: 32}
	s := HashSetFrom[*person, string]([]*person{anna, carl, dave, bill})

	fmt.Println(s)

	s.Remove(carl)

	fmt.Println(s)

	// Output:
	// [anna bill carl dave]
	// [anna bill dave]
}

func ExampleHashSet_RemoveSlice() {
	anna := &person{Name: "anna", ID: 94}
	bill := &person{Name: "bill", ID: 50}
	carl := &person{Name: "carl", ID: 10}
	dave := &person{Name: "dave", ID: 32}
	s := HashSetFrom[*person, string]([]*person{anna, carl, dave, bill})

	fmt.Println(s)

	s.RemoveSlice([]*person{anna, carl})

	fmt.Println(s)

	// Output:
	// [anna bill carl dave]
	// [bill dave]
}

func ExampleHashSet_RemoveSet() {
	anna := &person{Name: "anna", ID: 94}
	bill := &person{Name: "bill", ID: 50}
	carl := &person{Name: "carl", ID: 10}
	dave := &person{Name: "dave", ID: 32}
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
	anna := &person{Name: "anna", ID: 94}
	bill := &person{Name: "bill", ID: 50}
	carl := &person{Name: "carl", ID: 10}
	dave := &person{Name: "dave", ID: 32}
	s := HashSetFrom[*person, string]([]*person{anna, bill, carl, dave})

	idAbove50 := func(p *person) bool {
		return p.ID >= 50
	}

	fmt.Println(s)

	s.RemoveFunc(idAbove50)

	fmt.Println(s)

	// Output:
	// [anna bill carl dave]
	// [carl dave]
}

func ExampleHashSet_Contains() {
	anna := &person{Name: "anna", ID: 94}
	bill := &person{Name: "bill", ID: 50}
	carl := &person{Name: "carl", ID: 10}
	dave := &person{Name: "dave", ID: 32}
	s := HashSetFrom[*person, string]([]*person{anna, bill, carl})

	fmt.Println(s.Contains(anna))
	fmt.Println(s.Contains(dave))

	// Output:
	// true
	// false
}

func ExampleHashSet_ContainsSlice() {
	anna := &person{Name: "anna", ID: 94}
	bill := &person{Name: "bill", ID: 50}
	carl := &person{Name: "carl", ID: 10}
	dave := &person{Name: "dave", ID: 32}
	s := HashSetFrom[*person, string]([]*person{anna, bill, carl})

	fmt.Println(s.ContainsSlice([]*person{anna, bill}))
	fmt.Println(s.ContainsSlice([]*person{anna, bill, carl}))
	fmt.Println(s.ContainsSlice([]*person{carl, dave}))

	// Output:
	// true
	// true
	// false
}

func ExampleHashSet_Subset() {
	anna := &person{Name: "anna", ID: 94}
	bill := &person{Name: "bill", ID: 50}
	carl := &person{Name: "carl", ID: 10}
	dave := &person{Name: "dave", ID: 32}
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
	anna := &person{Name: "anna", ID: 94}
	bill := &person{Name: "bill", ID: 50}
	carl := &person{Name: "carl", ID: 10}
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
	anna := &person{Name: "anna", ID: 94}
	bill := &person{Name: "bill", ID: 50}
	carl := &person{Name: "carl", ID: 10}
	dave := &person{Name: "dave", ID: 32}
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
	anna := &person{Name: "anna", ID: 94}
	bill := &person{Name: "bill", ID: 50}
	carl := &person{Name: "carl", ID: 10}
	dave := &person{Name: "dave", ID: 32}
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
	anna := &person{Name: "anna", ID: 94}
	bill := &person{Name: "bill", ID: 50}
	carl := &person{Name: "carl", ID: 10}
	dave := &person{Name: "dave", ID: 32}
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
	anna := &person{Name: "anna", ID: 94}
	bill := &person{Name: "bill", ID: 50}
	carl := &person{Name: "carl", ID: 10}
	dave := &person{Name: "dave", ID: 32}

	s1 := HashSetFrom[*person, string]([]*person{anna, bill, carl})
	s2 := HashSetFrom[*person, string]([]*person{anna, bill, carl})
	s3 := HashSetFrom[*person, string]([]*person{anna, bill, dave})

	fmt.Println(s1.Equal(s2))
	fmt.Println(s1.Equal(s3))

	// Output:
	// true
	// false
}

func ExampleHashSet_EqualSlice() {
	anna := &person{Name: "anna", ID: 94}
	bill := &person{Name: "bill", ID: 50}
	carl := &person{Name: "carl", ID: 10}
	dave := &person{Name: "dave", ID: 32}

	s := HashSetFrom[*person, string]([]*person{anna, bill, carl})

	fmt.Println(s.EqualSlice([]*person{bill, anna, carl}))
	fmt.Println(s.EqualSlice([]*person{anna, anna, bill, carl}))
	fmt.Println(s.EqualSlice([]*person{dave, bill, carl}))

	// Output:
	// true
	// true
	// false
}

func ExampleHashSet_EqualSliceSet() {
	anna := &person{Name: "anna", ID: 94}
	bill := &person{Name: "bill", ID: 50}
	carl := &person{Name: "carl", ID: 10}
	dave := &person{Name: "dave", ID: 32}

	s := HashSetFrom[*person, string]([]*person{anna, bill, carl})

	fmt.Println(s.EqualSliceSet([]*person{bill, anna, carl}))
	fmt.Println(s.EqualSliceSet([]*person{anna, anna, bill, carl}))
	fmt.Println(s.EqualSliceSet([]*person{dave, bill, carl}))

	// Output:
	// true
	// false
	// false
}

func ExampleHashSet_Copy() {
	anna := &person{Name: "anna", ID: 94}
	bill := &person{Name: "bill", ID: 50}
	carl := &person{Name: "carl", ID: 10}
	s := HashSetFrom[*person, string]([]*person{anna, bill, carl})
	c := s.Copy()

	fmt.Println(c)

	// Output:
	// [anna bill carl]
}

func ExampleHashSet_Slice() {
	anna := &person{Name: "anna", ID: 94}
	bill := &person{Name: "bill", ID: 50}
	carl := &person{Name: "carl", ID: 10}
	s := HashSetFrom[*person, string]([]*person{anna, bill, carl})

	slice := s.Slice()
	sort.Slice(slice, func(a, b int) bool {
		return slice[a].ID < slice[b].ID
	})

	fmt.Println(slice)

	// Output:
	// [carl bill anna]
}

func ExampleHashSet_String() {
	anna := &person{Name: "anna", ID: 94}
	bill := &person{Name: "bill", ID: 50}
	carl := &person{Name: "carl", ID: 10}
	s := HashSetFrom[*person, string]([]*person{anna, bill, carl})

	fmt.Println(s.String())

	// Output:
	// [anna bill carl]
}

// TODO: will not work as long as [HashFunc] cannot be derived from the type parameters.
func ExampleHashSet_UnmarshalJSON() {
	// type Foo struct {
	// 	Persons *HashSet[*person, string] `json:"persons"`
	// }

	// in := `{"persons":[{"Name":"anna","ID":94},{"Name":"bill","ID":50},{"Name":"bill","ID":50},{"Name":"carl","ID":10}]}`
	// var out Foo

	// _ = json.Unmarshal([]byte(in), &out)

	// fmt.Println(out.Persons)

	// Output:
}
