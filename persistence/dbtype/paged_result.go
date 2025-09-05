package dbtype

type PagedResult[T any] struct {
	Data       []T
	TotalCount int64
}
