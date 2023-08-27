package postgres

import (
	"context"
	"database/sql"

	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/utils"
	"github.com/sicozz/papyrus/utils/constants"
)

type postgresPFileRepository struct {
	Conn *sql.DB
	log  utils.AggregatedLogger
}

// NewPostgresPFileRepository will create an object that represent the PFileRepository interface
func NewPostgresPFileRepository(conn *sql.DB) domain.PFileRepository {
	logger := utils.NewAggregatedLogger(constants.Repository, constants.PFile)
	return &postgresPFileRepository{conn, logger}
}

func (r *postgresPFileRepository) GetAll(ctx context.Context) (res []domain.PFile, err error) {
	query :=
		`SELECT
			uuid,
			pf.code,
			date_creation,
			date_input,
			pft.description as pfile_type,
			pfst.description as pfile_state,
			pfsg.description as pfile_stage,
			dir,
			user_revision,
			user_approval
		FROM
			pfile pf
			INNER JOIN pfile_type pft ON(pf.type = pft.code)
			INNER JOIN pfile_state pfst ON(pf.state = pfst.code)
			INNER JOIN pfile_stage pfsg ON(pf.stage = pfsg.code)`

	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [GetAll] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		r.log.Err("In [GetAll] failed to exec statement ->", err)
	}

	res = make([]domain.PFile, 0)
	for rows.Next() {
		t := domain.PFile{}
		err = rows.Scan(
			&t.Uuid,
			&t.Code,
			&t.DateCreation,
			&t.DateInput,
			&t.Type,
			&t.State,
			&t.Stage,
			&t.Dir,
			&t.RevUser,
			&t.AppUser,
		)

		if err != nil {
			r.log.Err("IN [GetAll]: failed to scan pfile ->", err)
			return nil, err
		}

		res = append(res, t)
	}

	return
}
