package postgres

import (
	api "github.com/dankobg/fluffly/api/gen"
	p "github.com/go-jet/jet/v2/postgres"
)

func addSelectTotalCount(selectClause p.ProjectionList) {
	selectClause = append(selectClause, p.COUNT(p.STAR).OVER().AS("total_count"))
}

func getSelectTotalCount(pagination *api.PaginationParams) p.ProjectionList {
	var selectClause p.ProjectionList
	if pagination != nil {
		selectClause = append(selectClause, p.COUNT(p.STAR).OVER().AS("total_count"))
	}
	return selectClause
}

func getLimitOffset(q p.SelectStatement, pagination *api.PaginationParams) p.SelectStatement {
	if pagination != nil {
		limit := int64(pagination.PageSize)
		offset := (int64(pagination.Page) - 1) * limit
		return q.LIMIT(limit).OFFSET(offset)
	}
	return q
}
