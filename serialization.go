// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package set

import "encoding/json"

// marshalJSON will serialize a Serializable[T] into a json byte array
func marshalJSON[T any](s Collection[T]) ([]byte, error) {
	return json.Marshal(s.Slice())
}

// unmarshalJSON replaces the contents of s with the JSON array in data.
func unmarshalJSON[T any](s Collection[T], data []byte) error {
	slice := make([]T, 0)
	err := json.Unmarshal(data, &slice)
	if err != nil {
		return err
	}
	for _, item := range s.Slice() {
		s.Remove(item)
	}
	s.InsertSlice(slice)
	return nil
}
