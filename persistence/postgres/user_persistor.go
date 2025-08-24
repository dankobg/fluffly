package postgres

import (
	"context"
	"fmt"

	"github.com/dankobg/fluffly/db/dbmodel"
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

func (p *PgUserPersistor) Create(ctx context.Context, in dbmodel.UserSetter) (dbmodel.User, error) {
	insertedUser, err := dbmodel.Users.Insert(&in).One(ctx, p.db)
	if err != nil {
		return dbmodel.User{}, fmt.Errorf("failed to create a user: %w", err)
	}
	return *insertedUser, nil
}

func (p *PgUserPersistor) Update(ctx context.Context, userID uuid.UUID, in dbmodel.UserSetter) (dbmodel.User, error) {
	updatedUser, err := dbmodel.Users.Update(in.UpdateMod(), dbmodel.UpdateWhere.Users.ID.EQ(userID)).One(ctx, p.db)
	if err != nil {
		return dbmodel.User{}, fmt.Errorf("failed to update a user: %w", err)
	}
	return *updatedUser, nil
}

func (p *PgUserPersistor) Delete(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	_, err := dbmodel.Users.Delete(dbmodel.DeleteWhere.Users.ID.EQ(userID)).Exec(ctx, p.db)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to delete a user: %w", err)
	}
	return userID, nil
}

func (p *PgUserPersistor) Get(ctx context.Context, userID uuid.UUID) (dbmodel.User, error) {
	user, err := dbmodel.FindUser(ctx, p.db, userID)
	if err != nil {
		return dbmodel.User{}, fmt.Errorf("failed to get a user: %w", err)
	}
	return *user, nil
}

func (p *PgUserPersistor) List(ctx context.Context) ([]dbmodel.User, error) {
	userRows, err := queries.ListUsers().All(ctx, p.db)
	users := make([]dbmodel.User, len(userRows))
	if err != nil {
		return users, fmt.Errorf("failed to list users: %w", err)
	}
	for i, ur := range userRows {
		users[i] = dbmodel.User{
			ID:        ur.ID,
			CreatedAt: ur.CreatedAt,
			UpdatedAt: ur.UpdatedAt,
		}
	}
	return nil, nil
}
