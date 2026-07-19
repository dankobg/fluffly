package dbtype

import (
	api "github.com/dankobg/fluffly/api/gen"
)

type ListUsersFilters struct {
	Page     *api.PaginationPage
	PageSize *api.PaginationPageSize
	Sort     *[]string
}
