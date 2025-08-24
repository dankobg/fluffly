package postgres

import (
	"context"
	"fmt"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/dankobg/fluffly/db/dbmodel"
	"github.com/dankobg/fluffly/db/queries"
	"github.com/dankobg/fluffly/persistence"
	"github.com/stephenafamo/bob"
)

var _ persistence.OrganizationPersistor = (*PgOrganizationPersistor)(nil)

type PgOrganizationPersistor struct {
	*PgPersistor
}

func NewPgOrganizationPersistor(ps *PgPersistor) *PgOrganizationPersistor {
	return &PgOrganizationPersistor{
		PgPersistor: ps,
	}
}

func (p *PgOrganizationPersistor) Create(ctx context.Context, in persistence.OrganizationCreate) (*dbmodel.Organization, error) {
	var org *dbmodel.Organization
	txErr := p.db.RunInTx(ctx, nil, func(ctx context.Context, tx bob.Executor) error {
		insertedOrg, err := dbmodel.Organizations.Insert(in.Org).One(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to insert an organization: %w", err)
		}
		insertedAddress, err := dbmodel.Addresses.Insert(in.Address).One(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to insert organization address: %w", err)
		}
		in.Contact.OrganizationID = omit.From(insertedOrg.ID)
		in.Contact.AddressID = omit.From(insertedAddress.ID)
		in.WorkHour.OrganizationID = omitnull.From(insertedOrg.ID)
		for _, photo := range in.Photos {
			photo.OrganizationID = omitnull.From(insertedOrg.ID)
		}
		for _, social := range in.Socials {
			social.OrganizationID = omit.From(insertedOrg.ID)
		}
		insertedContact, err := dbmodel.OrganizationContacts.Insert(in.Contact).One(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to insert organization contact: %w", err)
		}
		insertedWorkHour, err := dbmodel.OrganizationWorkHours.Insert(in.WorkHour).One(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to insert organization work hour: %w", err)
		}
		var insertedPhotos dbmodel.OrganizationPhotoSlice
		if len(in.Photos) > 0 {
			if insertedPhotos, err = dbmodel.OrganizationPhotos.Insert(bob.ToMods(in.Photos...)).All(ctx, tx); err != nil {
				return fmt.Errorf("failed to insert organization photos: %w", err)
			}
		}
		var insertedSocials dbmodel.OrganizationSocialSlice
		if len(in.Socials) > 0 {
			if insertedSocials, err = dbmodel.OrganizationSocials.Insert(bob.ToMods(in.Socials...)).All(ctx, tx); err != nil {
				return fmt.Errorf("failed to insert organization social platforms: %w", err)
			}
		}
		insertedOrg.R.OrganizationWorkHour = insertedWorkHour
		insertedOrg.R.OrganizationContact = insertedContact
		insertedOrg.R.OrganizationContact.R.Address = insertedAddress
		insertedOrg.R.OrganizationPhotos = insertedPhotos
		insertedOrg.R.OrganizationSocials = insertedSocials
		org = insertedOrg
		return nil
	})
	if txErr != nil {
		return nil, txErr
	}
	return org, nil
}

func (p *PgOrganizationPersistor) Update(ctx context.Context, organizationID int64, in dbmodel.OrganizationSetter) (*dbmodel.Organization, error) {
	updatedOrganization, err := dbmodel.Organizations.Update(in.UpdateMod(), dbmodel.UpdateWhere.Organizations.ID.EQ(organizationID)).One(ctx, p.db)
	if err != nil {
		return nil, fmt.Errorf("failed to update an organization: %w", err)
	}
	return updatedOrganization, nil
}

func (p *PgOrganizationPersistor) Delete(ctx context.Context, organizationID int64) (int64, error) {
	_, err := dbmodel.Organizations.Delete(dbmodel.DeleteWhere.Organizations.ID.EQ(organizationID)).Exec(ctx, p.db)
	if err != nil {
		return 0, fmt.Errorf("failed to delete an organization: %w", err)
	}
	return organizationID, nil
}

func (p *PgOrganizationPersistor) Get(ctx context.Context, organizationID int64) (queries.GetOrganizationByIdRow, error) {
	organizationRow, err := queries.GetOrganizationById(organizationID).One(ctx, p.db)
	if err != nil {
		return queries.GetOrganizationByIdRow{}, fmt.Errorf("failed to list organizations: %w", err)
	}
	return organizationRow, nil
}

func (p *PgOrganizationPersistor) List(ctx context.Context) (queries.AllListOrganizationsRow, error) {
	organizationRows, err := queries.ListOrganizations().All(ctx, p.db)
	if err != nil {
		return queries.AllListOrganizationsRow{}, fmt.Errorf("failed to list organizations: %w", err)
	}
	return organizationRows, nil
}
