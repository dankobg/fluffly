package convert

import (
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/oapi-codegen/nullable"
)

// PtrToNullable will return the [nullable.Nullable] type if the pointer is not nil
func PtrToNullable[T any](v *T) nullable.Nullable[T] {
	var out nullable.Nullable[T]
	if v != nil {
		out.Set(*v)
	}

	return out
}

// PtrToNullableNull will return the [nullable.Nullable] with pointer value or sets it to Null explicitly
func PtrToNullableNull[T any](v *T) nullable.Nullable[T] {
	var out nullable.Nullable[T]
	if v == nil {
		out.SetNull()
	} else {
		out.Set(*v)
	}

	return out
}

func NullableToOmit[T any](v nullable.Nullable[T]) omit.Val[T] {
	var out omit.Val[T]
	if v.IsSpecified() && !v.IsNull() {
		out.Set(v.MustGet())
	}

	return out
}

func NullableToOmitNull[T any](v nullable.Nullable[T]) omitnull.Val[T] {
	var out omitnull.Val[T]

	if v.IsSpecified() {
		if v.IsNull() {
			out.SetPtr(nil)
		} else {
			out.Set(v.MustGet())
		}
	}

	return out
}
