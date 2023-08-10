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

type postgresUserStateRepository struct {
	Conn *sql.DB
	log  utils.AggregatedLogger
}

/*
* NewPostgresUserStateRepository will create an object that represent the
* UserStateRepository interface
 */
func NewPostgresUserStateRepository(conn *sql.DB) domain.UserStateRepository {
	logger := utils.NewAggregatedLogger(constants.Repository, constants.UserState)
	return &postgresUserStateRepository{conn, logger}
}

func (r *postgresUserStateRepository) fetch(ctx context.Context, query string, args ...interface{}) (res []domain.UserState, err error) {
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

	res = make([]domain.UserState, 0)
	for rows.Next() {
		t := domain.UserState{}
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

func (r *postgresUserStateRepository) GetByCode(ctx context.Context, code int64) (res domain.UserState, err error) {
	query := `SELECT code, description FROM user_state WHERE code=$1`
	states, err := r.fetch(ctx, query, code)
	if err != nil {
		r.log.Error(err)
		return domain.UserState{}, err
	}

	if l := len(states); l != 1 {
		r.log.Error("Could not find user_state with code: ", code)
		return domain.UserState{}, err
	}

	res = states[0]
	return
}

func (r *postgresUserStateRepository) GetByDescription(ctx context.Context, desc string) (res domain.UserState, err error) {
	query := `SELECT code, description FROM user_state WHERE description=$1`
	states, err := r.fetch(ctx, query, desc)
	if err != nil {
		r.log.Error(err)
		return domain.UserState{}, err
	}

	if l := len(states); l != 1 {
		r.log.Error("Could not find user_state with description: ", desc)
		err = errors.New(fmt.Sprint("No user__state with description: ", desc))
		return domain.UserState{}, err
	}

	res = states[0]
	return
}

func (r *postgresUserStateRepository) GetAll(ctx context.Context) ([]domain.UserState, error) {
	query := `SELECT code, description FROM user_state`
	return r.fetch(ctx, query)
}
