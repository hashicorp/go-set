package set

import "encoding/json"

// serializable is an interface that allows a set to be serialized
type serializable[T any] interface {
	Slice() []T
	InsertSlice([]T) bool
}

// marshalJSON will serialize a Serializable[T] into a json byte array
func marshalJSON[T any](s serializable[T]) ([]byte, error) {
	return json.Marshal(s.Slice())
}

// unmarshalJSON will deserialize a json byte array into a Serializable[T]
func unmarshalJSON[T any](s serializable[T], data []byte) error {
	slice := make([]T, 0)
	err := json.Unmarshal(data, &slice)
	if err != nil {
		return err
	}
	s.InsertSlice(slice)
	return nil
}
