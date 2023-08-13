package postgres

import (
	"context"
	"database/sql"
	"errors"

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

// Retrieve all dirs
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

// Know if a dir exists by uuid
func (r *postgresDirRepository) ExistByUuid(ctx context.Context, uuid string) (res bool) {
	query := `SELECT COUNT(*) > 0 FROM dir WHERE uuid = $1`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ExistByUuid]: could not prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, uuid).Scan(&res)

	return
}

// Get dir by uuid
func (r *postgresDirRepository) GetByUuid(ctx context.Context, uuid string) (res domain.Dir, err error) {
	query := `SELECT * FROM dir WHERE uuid = $1`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [GetByUuid]: could not prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, uuid).Scan(
		&res.Uuid,
		&res.Name,
		&res.ParentDir,
		&res.Path,
		&res.Nchild,
	)

	if err != nil {
		r.log.Err("IN [GetByUuid]: could not scan rows ->", err)
		return
	}

	return
}

// Store a new dir
func (r *postgresDirRepository) Store(ctx context.Context, d *domain.Dir) (err error) {
	query :=
		`INSERT INTO dir (name, parent_dir, path, nchild)
		VALUES ($1, $2, $3, $4)
		RETURNING uuid`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [Store]: could not prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(
		ctx,
		d.Name,
		d.ParentDir,
		d.Path,
		d.Nchild,
	).Scan(&d.Uuid)

	if err != nil {
		r.log.Err("IN [Store]: could not scan rows ->", err)
		return
	}

	return
}

// Delete dir by uuid
func (r *postgresDirRepository) Delete(ctx context.Context, uuid string) (err error) {
	query := `DELETE FROM dir WHERE uuid=$1 RETURNING uuid`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [Delete]: could not prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, uuid).Scan(&uuid)
	if uuid == "" && err == nil {
		err = errors.New("Could not delete user")
	}

	return
}

func (r *postgresDirRepository) ChgName(ctx context.Context, uuid string, nName string) (err error) {
	query := `UPDATE dir SET name=$1 WHERE uuid=$2`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ChgName]: could not prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.QueryContext(ctx, nName, uuid)

	return
}

func (r *postgresDirRepository) ChgParentDir(ctx context.Context, uuid string, nPUuid string) (err error) {
	query := `UPDATE dir SET parent_dir=$1 WHERE uuid=$2`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ChgParentDir]: could not prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.QueryContext(ctx, nPUuid, uuid)

	return
}

func (r *postgresDirRepository) IncNchild(ctx context.Context, uuid string, nNchild int) (err error) {
	query := `UPDATE dir SET nchild=nchild + $1 WHERE uuid=$2`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [IncNchild]: could not prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.QueryContext(ctx, nNchild, uuid)

	return
}

func (r *postgresDirRepository) DecNchild(ctx context.Context, uuid string, nNchild int) (err error) {
	query := `UPDATE dir SET nchild=nchild - $1 WHERE uuid=$2`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [DecNchild]: could not prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.QueryContext(ctx, nNchild, uuid)

	return
}
