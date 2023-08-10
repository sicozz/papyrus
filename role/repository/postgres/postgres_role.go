package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/utils"
	"github.com/sicozz/papyrus/utils/constants"
)

type postgresRoleRepository struct {
	Conn *sql.DB
	log  utils.AggregatedLogger
}

/*
* NewPostgresRoleRepository will create an object that represent the
* RoleRepository interface
 */
func NewPostgresRoleRepository(conn *sql.DB) domain.RoleRepository {
	// TODO: Add layer enum and domain enum in utils package
	logger := utils.NewAggregatedLogger(constants.Repository, constants.Role)
	return &postgresRoleRepository{conn, logger}
}

func (r *postgresRoleRepository) fetch(ctx context.Context, query string, args ...interface{}) (res []domain.Role, err error) {
	rows, err := r.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			r.log.Error(errRow)
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
			r.log.Error(err)
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
		r.log.Error(err)
		return domain.Role{}, err
	}

	if l := len(roles); l != 1 {
		r.log.Error("Could not find role with code:", code)
		return domain.Role{}, err
	}

	res = roles[0]
	return
}

func (r *postgresRoleRepository) GetByDescription(ctx context.Context, desc string) (res domain.Role, err error) {
	query := `SELECT code, description FROM role WHERE description=$1`
	roles, err := r.fetch(ctx, query, desc)
	if err != nil {
		r.log.Error(err)
		return domain.Role{}, err
	}

	if l := len(roles); l != 1 {
		r.log.Error("Could not find role with description:", desc)
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
