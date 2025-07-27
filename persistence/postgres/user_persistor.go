package postgres

import (
	"context"
	"fmt"

	"github.com/dankobg/fluffly/db/model"
	"github.com/dankobg/fluffly/db/queries"
	"github.com/dankobg/fluffly/persistence"
	"github.com/google/uuid"
)

var _ persistence.UserPersistor = (*PgUserPersistor)(nil)

type PgUserPersistor struct {
	*PgPersistor
}

func NewPgUserPersistor(ps *PgPersistor) *PgUserPersistor {
	return &PgUserPersistor{
		PgPersistor: ps,
	}
}

func (p *PgUserPersistor) Create(ctx context.Context, in model.UserSetter) (model.User, error) {
	insertedUser, err := model.Users.Insert(&in).One(ctx, p.db)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to create a user: %w", err)
	}
	return *insertedUser, nil
}

func (p *PgUserPersistor) Update(ctx context.Context, in model.UserSetter) (model.User, error) {
	updatedUser, err := model.Users.Update(in.UpdateMod(), model.UpdateWhere.Users.ID.EQ(in.ID.GetOrZero())).One(ctx, p.db)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to update a user: %w", err)
	}
	return *updatedUser, nil
}

func (p *PgUserPersistor) Delete(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	_, err := model.Users.Delete(model.DeleteWhere.Users.ID.EQ(userID)).Exec(ctx, p.db)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to delete a user: %w", err)
	}
	return userID, nil
}

func (p *PgUserPersistor) Get(ctx context.Context, userID uuid.UUID) (model.User, error) {
	user, err := model.FindUser(ctx, p.db, userID)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to get a user: %w", err)
	}
	return *user, nil
}

func (p *PgUserPersistor) List(ctx context.Context) ([]model.User, error) {
	userRows, err := queries.ListUsers().All(ctx, p.db)
	users := make([]model.User, len(userRows))
	if err != nil {
		return users, fmt.Errorf("failed to list users: %w", err)
	}
	for i, ur := range userRows {
		users[i] = model.User{
			ID:        ur.ID,
			CreatedAt: ur.CreatedAt,
			UpdatedAt: ur.UpdatedAt,
		}
	}
	return nil, nil
}
