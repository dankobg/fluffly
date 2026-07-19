package dbtype

import (
	api "github.com/dankobg/fluffly/api/gen"
)

type ListGeocodingResultsFilters struct {
	Page     *api.PaginationPage
	PageSize *api.PaginationPageSize
	Sort     *[]string
}
