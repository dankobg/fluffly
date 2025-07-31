package persistence

import (
	"context"

	"github.com/dankobg/fluffly/db/model"
	"github.com/google/uuid"
)

type UserPersistor interface {
	Create(ctx context.Context, in model.UserSetter) (model.User, error)
	Update(ctx context.Context, userID uuid.UUID, in model.UserSetter) (model.User, error)
	Delete(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	Get(ctx context.Context, userID uuid.UUID) (model.User, error)
	List(ctx context.Context) ([]model.User, error)
}

type OrganizationPersistor interface {
	Create(ctx context.Context, in model.OrganizationSetter) (model.Organization, error)
	Update(ctx context.Context, organizationID int64, in model.OrganizationSetter) (model.Organization, error)
	Delete(ctx context.Context, organizationID int64) (int64, error)
	Get(ctx context.Context, organizationID int64) (model.Organization, error)
	List(ctx context.Context) ([]model.Organization, error)
}

type Persistor interface {
	User() UserPersistor
	Organization() OrganizationPersistor
}
