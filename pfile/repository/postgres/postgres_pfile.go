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
			dir,
			version,
			term,
			subtype,
			resp_user
		FROM
			pfile pf
			INNER JOIN pfile_type pft ON(pf.type = pft.code)
			INNER JOIN pfile_state pfst ON(pf.state = pfst.code)`

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
			&t.Dir,
			&t.Version,
			&t.Term,
			&t.Subtype,
			&t.RespUser,
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
			dir,
			version,
			term,
			subtype,
			resp_user
		FROM
			pfile pf
			INNER JOIN pfile_type pft ON(pf.type = pft.code)
			INNER JOIN pfile_state pfst ON(pf.state = pfst.code)
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
		&res.Dir,
		&res.Version,
		&res.Term,
		&res.Subtype,
		&res.RespUser,
	)

	if err != nil {
		r.log.Err("IN [GetByUuid] failed to scan row ->", err)
		return
	}

	return
}

func (r *postgresPFileRepository) StoreUuid(ctx context.Context, pf domain.PFile, apps []domain.Approvation) (uuid string, err error) {
	// WARN: CONTINUE HERE with transaction
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
			dir,
			version,
			term,
			subtype,
			resp_user
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
			$12,
			$13
		)
		RETURNING uuid`

	tx, err := r.Conn.Begin()
	if err != nil {
		r.log.Err("IN [StoreUuid] failed to begin transaction ->", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query)
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
		pf.Dir,
		pf.Version,
		pf.Term,
		pf.Subtype,
		pf.RespUser,
	).Scan(&uuid)

	if err != nil {
		r.log.Err("IN [StoreUuid] failed to scan rows ->", err)
		return "", err
	}

	for _, ap := range apps {
		err = r.storeApprovation(ctx, tx, ap)
		if err != nil {
			return "", err
		}
	}

	err = tx.Commit()
	if err != nil {
		r.log.Err("IN [StoreUuid] failed to commit changes -> ", err)
		return "", err
	}

	return
}

func (r *postgresPFileRepository) GetApprovations(ctx context.Context, pfUuid string) (res []domain.Approvation, err error) {
	query :=
		`SELECT user_uuid, pfile_uuid, is_approved
		FROM approvation
		WHERE pfile_uuid = $1`

	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [GetApprovations] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, pfUuid)
	if err != nil {
		r.log.Err("In [GetApprovations] failed to exec statement ->", err)
	}

	res = []domain.Approvation{}
	for rows.Next() {
		t := domain.Approvation{}
		err = rows.Scan(
			&t.UserUuid,
			&t.PFileUuid,
			&t.IsApproved,
		)

		if err != nil {
			r.log.Err("IN [GetApprovations]: failed to scan approvation ->", err)
			return nil, err
		}

		res = append(res, t)
	}

	return
}

func (r *postgresPFileRepository) storeApprovation(ctx context.Context, tx *sql.Tx, ap domain.Approvation) (err error) {
	query :=
		`INSERT INTO approvation (user_uuid, pfile_uuid, is_approved)
		VALUES ($1, $2, $3)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [storeApprovation] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, ap.UserUuid, ap.PFileUuid, ap.IsApproved)
	if err != nil {
		r.log.Err("IN [storeApprovation] failed to exec statement ->", err)
		return
	}

	return
}

func (r *postgresPFileRepository) Delete(ctx context.Context, uuid string) (err error) {
	// BUG: Make it delete from file system
	tx, err := r.Conn.Begin()
	if err != nil {
		r.log.Err("IN [Delete] failed to begin transaction ->", err)
	}
	defer tx.Rollback()

	appQuery := `DELETE FROM approvation WHERE pfile_uuid = $1`
	appStmt, err := tx.PrepareContext(ctx, appQuery)
	if err != nil {
		r.log.Err("IN [Delete] failed to prepare context ->", err)
		return
	}
	defer appStmt.Close()

	_, err = appStmt.ExecContext(ctx, uuid)
	if uuid == "" && err != nil {
		r.log.Err("IN [Delete] failed to exec statement ->", err)
		return
	}

	query := `DELETE FROM pfile WHERE uuid = $1`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [Delete] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, uuid)
	// WARN: This thing is to be deleted when you apply prev exists check in usecase
	if uuid == "" || err != nil {
		r.log.Err("IN [Delete] failed to exec statement ->", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		r.log.Err("IN [Delete] failed to commit changes -> ", err)
		return err
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

func (r *postgresPFileRepository) Approve(ctx context.Context, pfUuid, userUuid string) (err error) {
	query := `UPDATE approvation SET is_approved = $1 WHERE pfile_uuid = $2 AND user_uuid = $3`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [Approve] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, true, pfUuid, userUuid)
	if err != nil {
		r.log.Err("IN [Approve] failed to exec statement ->", err)
		return
	}

	return
}

func (r *postgresPFileRepository) Activate(ctx context.Context, pfUuid, userUuid string) (err error) {
	query := `UPDATE pfile SET state = 2 WHERE uuid = $1 AND resp_user = $2`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [Activate] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, pfUuid, userUuid)
	if err != nil {
		r.log.Err("IN [Activate] failed to exec statement ->", err)
		return
	}

	return
}

func (r *postgresPFileRepository) ApprExistsByPK(ctx context.Context, pfUuid, userUuid string) (res bool) {
	query := `SELECT COUNT(*) > 0 FROM approvation WHERE pfile_uuid = $1 AND user_uuid = $2`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ApprExistsByPK] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, pfUuid, userUuid).Scan(&res)
	if err != nil {
		r.log.Err("IN [ApprExistsByPK] failed to exec statement ->", err)
	}

	return
}

func (r *postgresPFileRepository) IsApproved(ctx context.Context, uuid string) (res bool) {
	query := `SELECT COUNT(*) = 0 FROM approvation WHERE pfile_uuid = $1 AND is_approved = FALSE`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [IsApproved] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, uuid).Scan(&res)
	if err != nil {
		r.log.Err("IN [IsApproved] failed to exec statement ->", err)
	}

	return
}
