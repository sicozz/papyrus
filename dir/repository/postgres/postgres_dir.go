package postgres

import (
	"context"
	"database/sql"

	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/utils"
	"github.com/sicozz/papyrus/utils/constants"
)

type postgresDirRepository struct {
	Conn *sql.DB
	log  utils.AggregatedLogger
}

// NewPostgresDirRepository will create an object that represent the DirRepository interface
func NewPostgresDirRepository(conn *sql.DB) domain.DirRepository {
	logger := utils.NewAggregatedLogger(constants.Repository, constants.Dir)
	return &postgresDirRepository{conn, logger}
}

// Get dir by uuid
func (r *postgresDirRepository) GetAll(ctx context.Context) (res []domain.Dir, err error) {
	query := `SELECT * FROM dir`
	rows, err := r.Conn.QueryContext(ctx, query)
	if err != nil {
		res = nil
		return
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			r.log.Err("IN [GetAll]: could not close *rows ->", err)
		}
	}()

	res = make([]domain.Dir, 0)
	for rows.Next() {
		t := domain.Dir{}
		err = rows.Scan(
			&t.Uuid,
			&t.Name,
			&t.ParentDir,
			&t.Path,
			&t.Nchild,
		)

		if err != nil {
			r.log.Err("IN [GetAll]: could not scan dir ->", err)
		}
		res = append(res, t)
	}

	return
}
