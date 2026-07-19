package postgres

import (
	"context"
	"fmt"

	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
	"github.com/stephenafamo/scan"
)

func (pst *PgOrganizationPersistor) GetInvitationByID(ctx context.Context, invitationID int64) (models.OrganizationInvitation, error) {
	q := psql.Select(
		sm.Columns(models.OrganizationInvitations.Columns),
		sm.From(models.OrganizationInvitations.Name()),
		sm.Where(models.OrganizationInvitations.Columns.ID.EQ(psql.Arg(invitationID))),
	)

	invitation, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.OrganizationInvitation]())
	if err != nil {
		return models.OrganizationInvitation{}, fmt.Errorf("query organization invitation")
	}

	return invitation, nil
}

func (pst *PgOrganizationPersistor) CreateInvitation(ctx context.Context, in models.OrganizationInvitationSetter) (models.OrganizationInvitation, error) {
	q := models.OrganizationInvitations.Insert(&in, im.Returning(models.OrganizationInvitations.Columns))

	country, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.OrganizationInvitation]())
	if err != nil {
		return models.OrganizationInvitation{}, fmt.Errorf("insert organization invitation")
	}

	return country, nil
}

func (pst *PgOrganizationPersistor) DeleteInvitation(ctx context.Context, invitationID int64) error {
	q := models.OrganizationInvitations.Delete(dm.Where(models.OrganizationInvitations.Columns.ID.EQ(psql.Arg(invitationID))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete organization invitation: %w", err)
	}

	return nil
}

func (pst *PgOrganizationPersistor) UpdateInvitation(ctx context.Context, invitationID int64, in models.OrganizationInvitationSetter) (models.OrganizationInvitation, error) {
	q := models.OrganizationInvitations.Update(
		in.UpdateMod(),
		um.Where(models.OrganizationInvitations.Columns.ID.EQ(psql.Arg(invitationID))),
		um.Returning(models.OrganizationInvitations.Columns),
	)

	country, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.OrganizationInvitation]())
	if err != nil {
		return models.OrganizationInvitation{}, fmt.Errorf("update organization invitation")
	}

	return country, nil
}
