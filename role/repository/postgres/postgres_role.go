package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sicozz/papyrus/domain"
)

type postgresRoleRepository struct {
	Conn *sql.DB
}

/*
* NewPostgresRoleRepository will create an object that represent the
* RoleRepository interface
 */
func NewPostgresRoleRepository(conn *sql.DB) domain.RoleRepository {
	return &postgresRoleRepository{conn}
}

func (r *postgresRoleRepository) fetch(ctx context.Context, query string, args ...interface{}) (res []domain.Role, err error) {
	rows, err := r.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		domain.AgLog.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			domain.AgLog.Error(errRow)
		}
	}()

	res = make([]domain.Role, 0)
	for rows.Next() {
		t := domain.Role{}
		// Get from db
		err = rows.Scan(
			&t.Code,
			&t.Description,
		)

		if err != nil {
			domain.AgLog.Error(err)
			return nil, err
		}
		res = append(res, t)
	}

	return res, nil
}

func (r *postgresRoleRepository) GetByCode(ctx context.Context, code int64) (res domain.Role, err error) {
	query := `SELECT code, description FROM role WHERE code=$1`
	roles, err := r.fetch(ctx, query, code)
	if err != nil {
		domain.AgLog.Error(err)
		return domain.Role{}, err
	}

	if l := len(roles); l != 1 {
		domain.AgLog.Error("Could not find role with code:", code)
		return domain.Role{}, err
	}

	res = roles[0]
	return
}

func (r *postgresRoleRepository) GetByDescription(ctx context.Context, desc string) (res domain.Role, err error) {
	query := `SELECT code, description FROM role WHERE description=$1`
	roles, err := r.fetch(ctx, query, desc)
	if err != nil {
		domain.AgLog.Error(err)
		return domain.Role{}, err
	}

	if l := len(roles); l != 1 {
		domain.AgLog.Error("Could not find role with description:", desc)
		err = errors.New(fmt.Sprint("No role with description: ", desc))
		return domain.Role{}, err
	}

	res = roles[0]
	return
}

func (r *postgresRoleRepository) GetAll(ctx context.Context) ([]domain.Role, error) {
	query := `SELECT code, description FROM role`
	return r.fetch(ctx, query)
}
