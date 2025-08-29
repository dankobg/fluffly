package postgres

import (
	"context"
	"fmt"

	"github.com/dankobg/fluffly/db/gen/test/public/model"
	t "github.com/dankobg/fluffly/db/gen/test/public/table"
	"github.com/dankobg/fluffly/persistence"
	p "github.com/go-jet/jet/v2/postgres"
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

func (pu *PgUserPersistor) ListUsers(ctx context.Context) ([]model.User, error) {
	q := p.SELECT(t.User.AllColumns).
		FROM(t.User)
	var dest []model.User
	if err := q.QueryContext(ctx, pu.db, &dest); err != nil {
		return nil, err
	}
	return dest, nil
}

func (pu *PgUserPersistor) GetUserByID(ctx context.Context, userID uuid.UUID) (model.User, error) {
	q := p.SELECT(t.User.AllColumns).
		FROM(t.User).
		WHERE(t.User.ID.EQ(p.UUID(userID)))
	var dest model.User
	if err := q.QueryContext(ctx, pu.db, &dest); err != nil {
		return model.User{}, err
	}
	return dest, nil
}

func (pu *PgUserPersistor) DeleteUserByID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	q := t.User.DELETE().WHERE(t.User.ID.EQ(p.UUID(userID)))
	if _, err := q.ExecContext(ctx, pu.db); err != nil {
		return uuid.Nil, fmt.Errorf("failed to delete an user: %w", err)
	}
	return userID, nil
}

func (pu *PgUserPersistor) CreateUser(ctx context.Context, in persistence.UserSetter) (model.User, error) {
	cols, m := in.ToModel()
	q := t.User.INSERT(cols).
		MODEL(m).
		RETURNING(t.User.AllColumns)

	var dest model.User
	if err := q.QueryContext(ctx, pu.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to create an user: %w", err)
	}
	return dest, nil
}

func (pu *PgUserPersistor) UpdateUser(ctx context.Context, userID uuid.UUID, in persistence.UserSetter) (model.User, error) {
	cols, m := in.ToModel(true)

	q := t.User.UPDATE(cols).
		MODEL(m).
		WHERE(t.User.ID.EQ(p.UUID(userID))).
		RETURNING(t.User.AllColumns)

	var dest model.User
	if err := q.QueryContext(ctx, pu.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to update an user: %w", err)
	}
	return dest, nil
}
