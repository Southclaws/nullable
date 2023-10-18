package nullable

import (
	"bytes"
	"encoding/json"
)

type Nullable[T any] struct {
	v      *T
	isNull bool
}

func New[T any](v T) Nullable[T] {
	return Nullable[T]{v: &v}
}

func NewNull[T any]() Nullable[T] {
	return Nullable[T]{isNull: true}
}

func (n Nullable[T]) IsNull() bool {
	return n.isNull
}

func (n Nullable[T]) IsUnset() bool {
	return !n.isNull && n.v == nil
}

func (n Nullable[T]) Value() (T, bool) {
	if n.v != nil {
		return *n.v, true
	}

	var v T
	return v, false
}

func (u Nullable[T]) MarshalJSON() ([]byte, error) {
	if u.isNull {
		return []byte(`null`), nil
	}

	return json.Marshal(u.v)
}

func (u *Nullable[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		u.isNull = true
		return nil
	}

	type Alias = T

	var v Alias

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	u.v = &v

	return nil
}
