package ptr

// Of returns the address of any given value
func Of[T any](v T) *T {
	return &v
}

// Value will return the value of the pointer or the zero value if the pointer is nil
func Value[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}
