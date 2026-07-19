package postgres

import (
	"context"
	"fmt"

	"github.com/dankobg/fluffly/persistence"
	"github.com/dankobg/fluffly/persistence/dbcustom"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/types"
	"github.com/stephenafamo/scan"
)

var _ persistence.AnalyticsPersistor = (*PgAnalyticsPersistor)(nil)

type PgAnalyticsPersistor struct {
	*PgPersistor
}

func NewPgAnalyticsPersistor(ps *PgPersistor) *PgAnalyticsPersistor {
	return &PgAnalyticsPersistor{
		PgPersistor: ps,
	}
}

func (pst *PgAnalyticsPersistor) GetAnalyticsStats(ctx context.Context) (dbcustom.AnalyticsStats, error) {
	raw := `select json_build_object(
  'users', json_build_object(
    'total', (select count(*) from "user"),
    'by_state', (
      select json_build_object(
        'active', count(*) filter (where state = 'active'),
        'inactive', count(*) filter (where state = 'inactive')
      )
      from "identities"
    )
  ),
  'organizations', json_build_object(
    'total', (select count(*) from "organization"),
    'by_status', (
      select json_build_object(
        'approved', count(*) filter (where status = 'approved'),
        'pending', count(*) filter (where status = 'pending'),
        'rejected', count(*) filter (where status = 'rejected')
      )
      from "organization"
    )
  ),
  'adoptions', json_build_object(
    'total', (select count(*) from "adoption"),
    'by_status', (
      select json_build_object(
        'approved', count(*) filter (where status = 'approved'),
        'pending', count(*) filter (where status = 'pending'),
        'rejected', count(*) filter (where status = 'rejected')
      )
      from "adoption"
    )
  ),
  'animals', json_build_object(
    'total', (select count(*) from "animal"),
    'by_status', (
      select json_build_object(
        'adoptable', count(*) filter (where status = 'adoptable'),
        'pending', count(*) filter (where status = 'pending'),
        'adopted', count(*) filter (where status = 'adopted'),
        'reserved', count(*) filter (where status = 'reserved'),
        'rejected', count(*) filter (where status = 'rejected')
      )
      from "animal"
    )
  )
) as "stats";`

	q := psql.RawQuery(raw)

	type AnalyticsStatsDest struct {
		Stats types.JSON[dbcustom.AnalyticsStats]
	}

	analyticsStatsResult, err := bob.One(ctx, pst.exec, q, scan.StructMapper[AnalyticsStatsDest]())
	if err != nil {
		return dbcustom.AnalyticsStats{}, fmt.Errorf("query analytics stats")
	}

	return analyticsStatsResult.Stats.Val, nil
}

func (pst *PgAnalyticsPersistor) GetMyAnalyticsStats(ctx context.Context, userID uuid.UUID) (dbcustom.MyAnalyticsStats, error) {
	raw := `with my_orgs as (
  select o.*
  from organization o
  join organization_membership m
    on m.organization_id = o.id
  where m.user_id = ?
),
my_adoptions as (
 select a.*
 from adoption a
 where a.user_id = ?
),
my_animals as (
  select *
  from animal
  where user_id = ?
)
select json_build_object(
  'organizations', json_build_object(
    'total', (select count(*) from my_orgs),
    'by_status', (
      select json_build_object(
        'approved', count(*) filter (where status = 'approved'),
        'pending', count(*) filter (where status = 'pending'),
        'rejected', count(*) filter (where status = 'rejected')
      )
      from my_orgs
    )
  ),
  'adoptions', json_build_object(
    'total', (select count(*) from my_adoptions),
    'by_status', (
      select json_build_object(
        'approved', count(*) filter (where status = 'approved'),
        'pending', count(*) filter (where status = 'pending'),
        'rejected', count(*) filter (where status = 'rejected')
      )
      from my_adoptions
    )
  ),
  'animals', json_build_object(
    'posted', json_build_object(
      'total', (select count(*) from my_animals),
      'by_status', (
        select json_build_object(
          'adoptable', count(*) filter (where status = 'adoptable'),
          'pending', count(*) filter (where status = 'pending'),
          'adopted', count(*) filter (where status = 'adopted'),
          'reserved', count(*) filter (where status = 'reserved'),
          'rejected', count(*) filter (where status = 'rejected')
        )
        from my_animals
      )
    ),
    'favorites', (
      select count(*)
      from user_animal_like
      where user_id = ?
    ),
    'adoptions', (
      select count(*)
      from adoption
      where user_id = ?
    )
  )
) as stats;`

	q := psql.RawQuery(raw, userID, userID, userID, userID, userID)

	type MyAnalyticsStatsDest struct {
		Stats types.JSON[dbcustom.MyAnalyticsStats]
	}

	myAnalyticsStatsResult, err := bob.One(ctx, pst.exec, q, scan.StructMapper[MyAnalyticsStatsDest]())
	if err != nil {
		return dbcustom.MyAnalyticsStats{}, fmt.Errorf("query my analytics stats")
	}

	return myAnalyticsStatsResult.Stats.Val, nil
}
