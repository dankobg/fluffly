package server

import (
	api "github.com/dankobg/fluffly/api/gen"
)

var (
	defaultPage     = 1
	defaultPageSize = 50
	minPage         = 1
	minPageSize     = 1
	maxPageSize     = 500
)

func getPaginationParams(page *api.PaginationPage, pageSize *api.PaginationPageSize) api.PaginationParams {
	params := api.PaginationParams{Page: defaultPage, PageSize: defaultPageSize}
	if page != nil {
		params.Page = max(*page, minPage)
	}

	if pageSize != nil {
		params.PageSize = min(max(*pageSize, minPageSize), maxPageSize)
	}

	return params
}

func getPaginationMeta(page *api.PaginationPage, pageSize *api.PaginationPageSize, total int64) api.PaginationMeta {
	meta := api.PaginationMeta{Page: defaultPage, PageSize: defaultPageSize, Total: total}
	if page != nil {
		meta.Page = max(*page, minPage)
	}

	if pageSize != nil {
		meta.PageSize = min(max(*pageSize, minPageSize), maxPageSize)
	}

	return meta
}
