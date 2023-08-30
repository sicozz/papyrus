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
			name,
			fs_path,
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
			&t.Name,
			&t.FsPath,
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

func (r *postgresPFileRepository) GetByUuid(ctx context.Context, uuid string) (res domain.PFile, err error) {
	query :=
		`SELECT
			uuid,
			pf.code,
			name,
			fs_path,
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
			INNER JOIN pfile_stage pfsg ON(pf.stage = pfsg.code)
		WHERE
			uuid = $1`

	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [GetByUuid] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, uuid).Scan(
		&res.Uuid,
		&res.Code,
		&res.Name,
		&res.FsPath,
		&res.DateCreation,
		&res.DateInput,
		&res.Type,
		&res.State,
		&res.Stage,
		&res.Dir,
		&res.RevUser,
		&res.AppUser,
	)

	if err != nil {
		r.log.Err("IN [GetByUuid] failed to scan row ->", err)
		return
	}

	return
}

// Store a new dir
func (r *postgresPFileRepository) Store(ctx context.Context, pf domain.PFile) (uuid string, err error) {
	// TODO: Add fs_path column
	query :=
		`INSERT INTO pfile (
			code,
			name,
			fs_path,
			date_creation,
			date_input,
			type,
			state,
			stage,
			dir,
			user_revision,
			user_approval
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11
		)
		RETURNING uuid`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [Store] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(
		ctx,
		pf.Code,
		pf.Name,
		pf.FsPath,
		pf.DateCreation,
		pf.DateInput,
		1,
		1,
		1,
		pf.Dir,
		pf.RevUser,
		pf.AppUser,
	).Scan(&uuid)

	if err != nil {
		r.log.Err("IN [Store] failed to scan rows ->", err)
		return
	}

	return
}

func (r *postgresPFileRepository) StoreUuid(ctx context.Context, pf domain.PFile) (uuid string, err error) {
	// TODO: Add fs_path column
	query :=
		`INSERT INTO pfile (
			uuid,
			code,
			name,
			fs_path,
			date_creation,
			date_input,
			type,
			state,
			stage,
			dir,
			user_revision,
			user_approval
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12
		)
		RETURNING uuid`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [Store] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(
		ctx,
		pf.Uuid,
		pf.Code,
		pf.Name,
		pf.FsPath,
		pf.DateCreation,
		pf.DateInput,
		1,
		1,
		1,
		pf.Dir,
		pf.RevUser,
		pf.AppUser,
	).Scan(&uuid)

	if err != nil {
		r.log.Err("IN [Store] failed to scan rows ->", err)
		return
	}

	return
}

func (r *postgresPFileRepository) Delete(ctx context.Context, uuid string) (err error) {
	query := `DELETE FROM pfile WHERE uuid = $1`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [Delete] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, uuid)
	if uuid == "" && err == nil {
		r.log.Err("IN [Delete] failed to exec statement ->", err)
		return
	}

	return
}

func (r *postgresPFileRepository) ExistsByCode(ctx context.Context, code string) (res bool) {
	query := `SELECT COUNT(*) > 0 FROM pfile WHERE code = $1`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ExistsByCode] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, code).Scan(&res)
	if err != nil {
		r.log.Err("IN [ExistsByCode] failed to exec statement ->", err)
	}

	return
}

func (r *postgresPFileRepository) IsNameTaken(ctx context.Context, name string, dirUuid string) (res bool) {
	query := `SELECT COUNT(*) > 0 FROM pfile WHERE name = $1 AND dir = $2`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [IsNameTaken] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, name, dirUuid).Scan(&res)
	if err != nil {
		r.log.Err("IN [IsNameTaken] failed to exec statement ->", err)
	}

	return
}

func (r *postgresPFileRepository) ExistsTypeByDesc(ctx context.Context, desc string) (res bool) {
	query := `SELECT COUNT(*) > 0 FROM pfile_type WHERE description = $1`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ExistsTypeByDesc] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, desc).Scan(&res)
	if err != nil {
		r.log.Err("IN [ExistsTypeByDesc] failed to exec statement ->", err)
	}

	return
}

func (r *postgresPFileRepository) ExistsStateByDesc(ctx context.Context, desc string) (res bool) {
	query := `SELECT COUNT(*) > 0 FROM pfile_state WHERE description = $1`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ExistsStateByDesc] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, desc).Scan(&res)
	if err != nil {
		r.log.Err("IN [ExistsStateByDesc] failed to exec statement ->", err)
	}

	return
}

func (r *postgresPFileRepository) ExistsStageByDesc(ctx context.Context, desc string) (res bool) {
	query := `SELECT COUNT(*) > 0 FROM pfile_stage WHERE description = $1`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ExistsStageByDesc] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, desc).Scan(&res)
	if err != nil {
		r.log.Err("IN [ExistsStageByDesc] failed to exec statement ->", err)
	}

	return
}
