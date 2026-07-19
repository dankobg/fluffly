package dbtype

import (
	api "github.com/dankobg/fluffly/api/gen"
)

type MicrochipFilters struct {
	Page     *api.PaginationPage
	PageSize *api.PaginationPageSize
}
