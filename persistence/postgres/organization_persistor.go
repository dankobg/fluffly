package postgres

import (
	"context"
	"fmt"

	"github.com/dankobg/fluffly/db/model"
	"github.com/dankobg/fluffly/db/queries"
	"github.com/dankobg/fluffly/persistence"
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

func (p *PgOrganizationPersistor) Create(ctx context.Context, in model.OrganizationSetter) (model.Organization, error) {
	insertedOrganization, err := model.Organizations.Insert(&in).One(ctx, p.db)
	if err != nil {
		return model.Organization{}, fmt.Errorf("failed to create an organization: %w", err)
	}
	return *insertedOrganization, nil
}

func (p *PgOrganizationPersistor) Update(ctx context.Context, organizationID int64, in model.OrganizationSetter) (model.Organization, error) {
	updatedOrganization, err := model.Organizations.Update(in.UpdateMod(), model.UpdateWhere.Organizations.ID.EQ(organizationID)).One(ctx, p.db)
	if err != nil {
		return model.Organization{}, fmt.Errorf("failed to update an organization: %w", err)
	}
	return *updatedOrganization, nil
}

func (p *PgOrganizationPersistor) Delete(ctx context.Context, organizationID int64) (int64, error) {
	_, err := model.Organizations.Delete(model.DeleteWhere.Organizations.ID.EQ(organizationID)).Exec(ctx, p.db)
	if err != nil {
		return 0, fmt.Errorf("failed to delete an organization: %w", err)
	}
	return organizationID, nil
}

func (p *PgOrganizationPersistor) Get(ctx context.Context, organizationID int64) (model.Organization, error) {
	organization, err := model.FindOrganization(ctx, p.db, organizationID)
	if err != nil {
		return model.Organization{}, fmt.Errorf("failed to get an organization: %w", err)
	}
	return *organization, nil
}

func (p *PgOrganizationPersistor) List(ctx context.Context) ([]model.Organization, error) {
	organizationRows, err := queries.ListOrganizations().All(ctx, p.db)
	organizations := make([]model.Organization, len(organizationRows))
	if err != nil {
		return organizations, fmt.Errorf("failed to list organizations: %w", err)
	}
	for i, org := range organizationRows {
		organizations[i] = model.Organization{
			ID:               org.ID,
			ContactID:        org.ContactID,
			Name:             org.Name,
			Website:          org.Website,
			MissionStatement: org.MissionStatement,
			AdoptionPolicy:   org.AdoptionPolicy,
			AdoptionURL:      org.AdoptionURL,
			Distance:         org.Distance,
			Facebook:         org.Facebook,
			Twitter:          org.Twitter,
			Youtube:          org.Youtube,
			Instagram:        org.Instagram,
			Pinterest:        org.Pinterest,
			CreatedAt:        org.CreatedAt,
			UpdatedAt:        org.UpdatedAt,
		}
	}
	return nil, nil
}
