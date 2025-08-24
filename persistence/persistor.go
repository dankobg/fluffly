package persistence

import (
	"context"

	"github.com/dankobg/fluffly/db/dbmodel"
	"github.com/dankobg/fluffly/db/queries"
	"github.com/google/uuid"
)

type UserPersistor interface {
	Create(ctx context.Context, in dbmodel.UserSetter) (dbmodel.User, error)
	Update(ctx context.Context, userID uuid.UUID, in dbmodel.UserSetter) (dbmodel.User, error)
	Delete(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	Get(ctx context.Context, userID uuid.UUID) (dbmodel.User, error)
	List(ctx context.Context) ([]dbmodel.User, error)
}

type OrganizationPersistor interface {
	Create(ctx context.Context, in OrganizationCreate) (*dbmodel.Organization, error)
	Update(ctx context.Context, organizationID int64, in dbmodel.OrganizationSetter) (*dbmodel.Organization, error)
	Delete(ctx context.Context, organizationID int64) (int64, error)
	Get(ctx context.Context, organizationID int64) (queries.GetOrganizationByIdRow, error)
	List(ctx context.Context) (queries.AllListOrganizationsRow, error)
}

type CountryPersistor interface {
	Create(ctx context.Context, in dbmodel.CountrySetter) (*dbmodel.Country, error)
	Update(ctx context.Context, country_id int64, in dbmodel.CountrySetter) (*dbmodel.Country, error)
	Delete(ctx context.Context, country_id int64) (int64, error)
	Get(ctx context.Context, country_id int64) (*dbmodel.Country, error)
	List(ctx context.Context) (dbmodel.CountrySlice, error)
}

type Persistor interface {
	User() UserPersistor
	Organization() OrganizationPersistor
	Country() CountryPersistor
}

type OrganizationCreate struct {
	Org      *dbmodel.OrganizationSetter
	Contact  *dbmodel.OrganizationContactSetter
	Address  *dbmodel.AddressSetter
	WorkHour *dbmodel.OrganizationWorkHourSetter
	Photos   []*dbmodel.OrganizationPhotoSetter
	Socials  []*dbmodel.OrganizationSocialSetter
}
