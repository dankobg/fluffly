package postgres

import (
	"context"
	"fmt"

	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/um"
	"github.com/stephenafamo/scan"
)

func (pst *PgOrganizationPersistor) CreateMember(ctx context.Context, organizationID int64, userID uuid.UUID, in models.OrganizationMembershipSetter) (models.OrganizationMembership, error) {
	q := models.OrganizationMemberships.Insert(&in, im.Returning(models.OrganizationMemberships.Columns))

	membership, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.OrganizationMembership]())
	if err != nil {
		return models.OrganizationMembership{}, fmt.Errorf("insert organization membership")
	}

	return membership, nil
}

func (pst *PgOrganizationPersistor) DeleteMember(ctx context.Context, organizationID int64, userID uuid.UUID) error {
	q := models.OrganizationMemberships.Delete(
		dm.Where(models.OrganizationMemberships.Columns.OrganizationID.EQ(psql.Arg(organizationID)).
			And(models.OrganizationMemberships.Columns.UserID.EQ(psql.Arg(userID))),
		),
	)
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete organization membership: %w", err)
	}

	return nil
}

func (pst *PgOrganizationPersistor) UpdateMembership(ctx context.Context, organizationID int64, userID uuid.UUID, in models.OrganizationMembershipSetter) (models.OrganizationMembership, error) {
	q := models.OrganizationMemberships.Update(
		in.UpdateMod(),
		um.Where(models.OrganizationMemberships.Columns.OrganizationID.EQ(psql.Arg(organizationID)).
			And(models.OrganizationMemberships.Columns.UserID.EQ(psql.Arg(userID))),
		),
		um.Returning(models.OrganizationMemberships.Columns),
	)

	membership, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.OrganizationMembership]())
	if err != nil {
		return models.OrganizationMembership{}, fmt.Errorf("update organization membership")
	}

	return membership, nil
}
