# go-set

Provides the `set` package that implements a generic mathematical [set](https://en.wikipedia.org/wiki/Set) for Go. The package only provides a basic implementation that is optimized for correctness and convenience. This package is not thread-safe.

# Documentation

The full documentation is available on GoDoc.

# Example

Below are simple example of usages

```go
s := Set[int](10)
s.Insert(1)
s.InsertAll([]int{2, 3, 4})
s.Size() # 3
```

```go
s := From[string]([]string{"one", "two", "three"})
s.Contains("three") # true
s.Remove("one") # true
```


```go
a := From[int]([]int{2, 4, 6, 8})
b := From[int]([]int{4, 5, 6})
a.Intersect(b) # {4, 6}
```
