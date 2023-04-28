# go-set

[![Go Reference](https://pkg.go.dev/badge/github.com/hashicorp/go-set.svg)](https://pkg.go.dev/github.com/hashicorp/go-set)
[![Run CI Tests](https://github.com/hashicorp/go-set/actions/workflows/ci.yaml/badge.svg)](https://github.com/hashicorp/go-set/actions/workflows/ci.yaml)
[![GitHub](https://img.shields.io/github/license/hashicorp/go-set)](LICENSE)

The `go-set` repository provides a `set` package containing a few
generic [Set](https://en.wikipedia.org/wiki/Set) implementations for Go.

Each implementation is optimal for a particular use case.

`Set` is ideal for `comparable` types.
  - backed by `map` builtin
  - commonly used with `string`, `int`, simple `struct` types, etc.

`HashSet` is useful for types that implement a `Hash()` function.
  - backed by `map` builtin
  - commonly used with complex structs

`TreeSet` is useful for comparable data (via `Compare[T]`)
  - backed by Red-Black Binary Search Tree
  - commonly used with complex structs with extrinsic order
  - efficient iteration in sort order
  - additional methods `Min` / `Max` / `TopK` / `BottomK`

This package is not thread-safe.

# Documentation

The full documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/hashicorp/go-set).

# Motivation

Package `set` helps reduce the boiler plate of using a `map[<type>]struct{}` as a set.

Say we want to de-duplicate a slice of strings
```go
items := []string{"mitchell", "armon", "jack", "dave", "armon", "dave"}
```

A typical example of the classic way using `map` built-in:
```go
m := make(map[string]struct{})
for _, item := range items {
  m[item] = struct{}{}
}
list := make([]string, 0, len(items))
for k := range m {
  list = append(list, k)
}
```

The same result, but in one line using package `go-set`.
```go
list := set.From[string](items).Slice()
```

# Set

The `go-set` package includes `Set` for types that satisfy the `comparable` constraint.
Note that complex structs (i.e. is a pointer or contains pointer fields) will not be
"comparable" in the sense of deep equality, but rather in the sense of pointer addresses.
The `Set` type should be used with builtin types like `string`, `int`, or simple struct
types with no pointers.

# HashSet

The `go-set` package includes `HashSet` for types that implement a `Hash()` function.
The custom type must satisfy `HashFunc[H Hash]` - essentially any `Hash()` function
that returns a `string` or `integer`. This enables types to use string-y hash
functions like `md5`, `sha1`, or even `GoString()`, but also enables types to
implement an efficient hash function using a hash code based on prime multiples.

# TreeSet

The `go-set` package includes `TreeSet` for creating sorted sets. A `TreeSet` may
be used with any type `T` as the comparison between elements is provided by implementing
`Compare[T]`. The `Cmp[builtin]` helper provides a convenient implementation of
`Compare` for `builtin` types like `string` or `int`. A `TreeSet` is backed by
an underlying balanced binary search tree, making operations like in-order traversal
efficient, in addition to enabling functions like `Min()`, `Max()`, `TopK()`, and
`BottomK()`.


### Methods

Implements the following set operations

- Insert
- InsertSlice
- InsertSet
- Remove
- RemoveSlice
- RemoveSet
- RemoveFunc
- Contains
- ContainsAll
- ContainsSlice
- Subset
- Size
- Empty
- Union
- Difference
- Intersect

Provides helper methods

- Equal
- Copy
- Slice
- String

TreeSet helper methods
- Min
- Max
- TopK
- BottomK
- FirstAbove
- FirstAboveEqual
- Above
- AboveEqual
- FirstBelow
- FirstBelowEqual
- Below
- BelowEqual

# Install

```
go get github.com/hashicorp/go-set@latest
```

```
import "github.com/hashicorp/go-set"
```

# Set Examples

Below are simple example usages of `Set`

```go
s := set.New[int](10)
s.Insert(1)
s.InsertSlice([]int{2, 3, 4})
s.Size()
```

```go
s := set.From[string]([]string{"one", "two", "three"})
s.Contains("three")
s.Remove("one")
```


```go
a := set.From[int]([]int{2, 4, 6, 8})
b := set.From[int]([]int{4, 5, 6})
a.Intersect(b)
```

# HashSet Examples

Below are simple example usages of `HashSet`

(using a hash code)
```go
type inventory struct {
    item   int
    serial int
}

func (i *inventory) Hash() int {
    code := 3 * item * 5 * serial
    return code
}

i1 := &inventory{item: 42, serial: 101}

s := set.NewHashSet[*inventory, int](10)
s.Insert(i1)
```

(using a string hash)
```go
type employee struct {
    name string
    id   int
}

func (e *employee) Hash() string {
    return fmt.Sprintf("%s:%d", e.name, e.id)
}

e1 := &employee{name: "armon", id: 2}

s := set.NewHashSet[*employee, string](10)
s.Insert(e1)
```

# TreeSet Examples

Below are simple example usages of `TreeSet`

(using `Cmp` as `Compare`)

```go
ts := NewTreeSet[int, Compare[int]](Cmp[int])
ts.Insert(5)
```

(using custom `Compare`)

```go
type waypoint struct {
    distance int
    name     string
}

cmp := func(w1, w2 *waypoint) int {
    return w1.distance - w2.distance
}

ts := NewTreeSet[*waypoint, Compare[*waypoint]](cmp)
ts.Insert(&waypoint{distance: 42, name: "tango"})
ts.Insert(&waypoint{distance: 13, name: "alpha"})
ts.Insert(&waypoint{distance: 71, name: "xray"})
```

