# go-set

![GoDoc](https://godoc.org/github.com/hashicorp/go-set?status.svg)
[![Run CI Tests](https://github.com/hashicorp/go-set/actions/workflows/ci.yaml/badge.svg)](https://github.com/hashicorp/go-set/actions/workflows/ci.yaml)
![GitHub](https://img.shields.io/github/license/shoenig/nomad-pledge-driver?style=flat-square)

Provides the `set` package that implements a generic mathematical [set](https://en.wikipedia.org/wiki/Set) for Go. The package only provides a basic implementation that is optimized for correctness and convenience. This package is not thread-safe.

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

# Hash Function

In addition to `Set`, there is `HashSet` for types that implement a `Hash()` function.
The custom type must satisfy `HashFunc[H Hash]` - essentially any `Hash()`
function that returns a `string` or `integer`. This enables types to use string-y
hash functions like `md5`, `sha1`, or even `GoString()`, but also enables types
to implement an efficient hash function using a hash code based on prime multiples.

### Methods

Implements the following set operations

- Insert
- InsertAll
- InsertSet
- Remove
- RemoveAll
- RemoveSet
- Contains
- ContainsAll
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
s.InsertAll([]int{2, 3, 4})
s.Size() # 3
```

```go
s := set.From[string]([]string{"one", "two", "three"})
s.Contains("three") # true
s.Remove("one") # true
```


```go
a := set.From[int]([]int{2, 4, 6, 8})
b := set.From[int]([]int{4, 5, 6})
a.Intersect(b) # {4, 6}
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
