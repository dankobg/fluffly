package ptr

// Of returns the address of any given value
func Of[T any](v T) *T {
	return &v
}
