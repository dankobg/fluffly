package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/persistence"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
	"github.com/stephenafamo/scan"
)

var _ persistence.GeocodingResultPersistor = (*PgGeocodingResultPersistor)(nil)

type PgGeocodingResultPersistor struct {
	*PgPersistor
}

func NewPgGeocodingResultPersistor(ps *PgPersistor) *PgGeocodingResultPersistor {
	return &PgGeocodingResultPersistor{
		PgPersistor: ps,
	}
}

type (
	ErrGeocodingResultIntegrityViolation struct{ errIntegrityViolation }
	ErrGeocodingResultUniqueViolation    struct{ errUniqueViolation }
)

var (
	ErrGeocodingResultNotFound  = errors.New("geocoding result not found")
	errGeocodingResultIntegrity = ErrGeocodingResultIntegrityViolation{}
)

func convertGeocodingResultPgError(pgErr *pgconn.PgError) error {
	if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		return errGeocodingResultIntegrity
	}

	return pgErr
}

func (pst *PgGeocodingResultPersistor) ListGeocodingResults(ctx context.Context, filters dbtype.ListGeocodingResultsFilters) (dbtype.PagedResult[models.GeocodingResult], error) {
	q := psql.Select(
		sm.Columns(
			models.GeocodingResults.Columns.Except("coords"),
			psql.Raw(`ST_AsBinary("geocoding_result"."coords") AS "coords"`),
		),
		sm.From(models.GeocodingResults.Name()),
		sm.GroupBy(models.GeocodingResults.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.GeocodingResults.Columns.Except(
		models.GeocodingResults.Columns.Coords.String(),
	).Names())
	addPagination(&q, filters.Page, filters.PageSize)

	type ListGeocodingResultsRow struct {
		models.GeocodingResult
		TotalCount int64
	}

	geocodingResults, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListGeocodingResultsRow]())
	if err != nil {
		return dbtype.PagedResult[models.GeocodingResult]{}, fmt.Errorf("query geocoding results")
	}

	result := dbtype.PagedResult[models.GeocodingResult]{
		Data: make([]models.GeocodingResult, len(geocodingResults)),
	}
	for i, row := range geocodingResults {
		result.Data[i] = row.GeocodingResult
	}

	if len(geocodingResults) > 0 {
		result.TotalCount = geocodingResults[0].TotalCount
	}

	return result, nil
}

func (pst *PgGeocodingResultPersistor) GetGeocodingResultByID(ctx context.Context, geocodingresultID int64) (models.GeocodingResult, error) {
	q := psql.Select(
		sm.Columns(
			models.GeocodingResults.Columns.Except("coords"),
			psql.Raw(`ST_AsBinary("geocoding_result"."coords") AS "coords"`),
		),
		sm.From(models.GeocodingResults.Name()),
		sm.Where(models.GeocodingResults.Columns.ID.EQ(psql.Arg(geocodingresultID))),
	)

	geocodingResult, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.GeocodingResult]())
	if err != nil {
		return models.GeocodingResult{}, fmt.Errorf("query geocoding result")
	}

	return geocodingResult, nil
}

func (pst *PgGeocodingResultPersistor) GetGeocodingResultByQuery(ctx context.Context, query string) (models.GeocodingResult, error) {
	q := psql.Select(
		sm.Columns(
			models.GeocodingResults.Columns.Except("coords"),
			psql.Raw(`ST_AsBinary("geocoding_result"."coords") AS "coords"`),
		),
		sm.From(models.GeocodingResults.Name()),
		sm.Where(models.GeocodingResults.Columns.Query.EQ(psql.Arg(query))),
	)

	geocodingResult, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.GeocodingResult]())
	if err != nil {
		return models.GeocodingResult{}, fmt.Errorf("query geocoding result by query")
	}

	return geocodingResult, nil
}

func (pst *PgGeocodingResultPersistor) DeleteGeocodingResultByID(ctx context.Context, geocodingresultID int64) (int64, error) {
	q := models.GeocodingResults.Delete(dm.Where(models.GeocodingResults.Columns.ID.EQ(psql.Arg(geocodingresultID))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return 0, fmt.Errorf("delete geocoding result: %w", err)
	}

	return geocodingresultID, nil
}

func (pst *PgGeocodingResultPersistor) DeleteGeocodingResults(ctx context.Context, ids []int64) error {
	geocodingResultIDs := make([]any, len(ids))
	for i, id := range ids {
		geocodingResultIDs[i] = id
	}

	q := models.GeocodingResults.Delete(dm.Where(models.GeocodingResults.Columns.ID.In(psql.Arg(geocodingResultIDs...))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete geocoding results: %w", err)
	}

	return nil
}

func (pst *PgGeocodingResultPersistor) CreateGeocodingResult(ctx context.Context, in models.GeocodingResultSetter) (models.GeocodingResult, error) {
	q := psql.Insert(
		im.Into(models.GeocodingResults.Name(), "query", "coords"),
		im.Values(psql.Arg(in.Query.MustGet()), psql.Raw("ST_GeomFromEWKB(?)", in.Coords.MustGet())),
		im.Returning(
			models.GeocodingResults.Columns.Except("coords"),
			psql.Raw(`ST_AsBinary("geocoding_result"."coords") AS "coords"`),
		),
	)

	geocodingResult, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.GeocodingResult]())
	if err != nil {
		return models.GeocodingResult{}, fmt.Errorf("insert geocoding result")
	}

	return geocodingResult, nil
}

func (pst *PgGeocodingResultPersistor) UpdateGeocodingResult(ctx context.Context, geocodingresultID int64, in models.GeocodingResultSetter) (models.GeocodingResult, error) {
	// q := models.GeocodingResults.Update(
	// 	in.UpdateMod(),
	// 	um.Where(models.GeocodingResults.Columns.ID.EQ(psql.Arg(geocodingresultID))),
	// 	um.Returning(models.GeocodingResults.Columns),
	// )

	// @TODO: hack untill i can directly pass the postgis `geography` type
	q := psql.Update(
		um.Table(models.GeocodingResults.Name()),
		um.Where(models.GeocodingResults.Columns.ID.EQ(psql.Arg(geocodingresultID))),
		um.Returning(models.GeocodingResults.Columns.Except("coords"),
			psql.Raw(`ST_AsBinary("coords") AS "coords"`)),
	)
	if in.Query.IsValue() {
		q.Apply(um.SetCol("query").ToArg(in.Query.MustGet()))
	}

	if in.Coords.IsValue() {
		q.Apply(um.SetCol("coords").To(psql.Raw("ST_GeomFromEWKB(?)", in.Coords.MustGet())))
	}

	geocodingResult, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.GeocodingResult]())
	if err != nil {
		return models.GeocodingResult{}, fmt.Errorf("update geocoding result")
	}

	return geocodingResult, nil
}
