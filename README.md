# go-set

[![Run CI Tests](https://github.com/hashicorp/go-set/actions/workflows/ci.yaml/badge.svg)](https://github.com/hashicorp/go-set/actions/workflows/ci.yaml)

Provides the `set` package that implements a generic mathematical [set](https://en.wikipedia.org/wiki/Set) for Go. The package only provides a basic implementation that is optimized for correctness and convenience. This package is not thread-safe.

# Documentation

The full documentation is available on [pkg.dev](https://pkg.go.dev/github.com/hashicorp/go-set).

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
list := set.From[string](items).List()
```

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
- Union
- Difference
- Intersect

Provides helper methods

- Equal
- Copy
- List
- String

# Install

```
go get github.com/hashicorp/go-set@latest
```

```
import "github.com/hashicorp/go-set"
```

# Example

Below are simple example of usages

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
