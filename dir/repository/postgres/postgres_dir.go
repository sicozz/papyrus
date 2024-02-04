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

// Get the number of children of a directory
func (r *postgresDirRepository) GetNChild(ctx context.Context, uuid string) (nChild int, err error) {
	query := `SELECT COUNT(*) FROM dir WHERE parent_dir = $1`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [GetNchild] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, uuid).Scan(&nChild)
	if err != nil {
		r.log.Err("IN [GetNchild] failed to exec statement ->", err)
	}

	// Avoid problems with root directory
	if uuid == constants.RootDirUuid {
		nChild -= 1
	}

	return
}

// Get the display path of a directory
func (r *postgresDirRepository) GetPath(ctx context.Context, uuid string) (path string, err error) {
	query := `SELECT sp_getPath($1)`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [GetPath] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, uuid).Scan(&path)
	if err != nil {
		r.log.Err("IN [GetPath] failed to exec statement ->", err)
	}

	return
}

// Get the depth of a directory
func (r *postgresDirRepository) GetDepth(ctx context.Context, uuid string) (depth int, err error) {
	query := `SELECT sp_getDepth($1)`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [GetDepth] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, uuid).Scan(&depth)
	if err != nil {
		r.log.Err("IN [GetDepth] failed to exec statement ->", err)
	}

	return
}

// Check if a dir with that name already exists in the parent dir
func (r *postgresDirRepository) IsNameTaken(ctx context.Context, name string, destUuid string) (res bool) {
	query := `SELECT COUNT(*) > 0 FROM dir WHERE parent_dir = $1 AND name = $2`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [IsNameTaken] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, destUuid, name).Scan(&res)
	if err != nil {
		r.log.Err("IN [IsNameTaken] failed to exec statement ->", err)
	}

	return
}

// Check if a dir with that name already exists in the parent dir
func (r *postgresDirRepository) IsSubDir(ctx context.Context, uuid string, destUuid string) (res bool) {
	query := `SELECT sp_isSubDir($1, $2)`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [IsSubDir] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, uuid, destUuid).Scan(&res)
	if err != nil {
		r.log.Err("IN [IsSubDir] failed to exec statement ->", err)
	}

	return
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
			r.log.Err("IN [GetAll] failed to close *rows ->", err)
		}
	}()

	res = make([]domain.Dir, 0)
	for rows.Next() {
		t := domain.Dir{}
		err = rows.Scan(
			&t.Uuid,
			&t.Name,
			&t.ParentDir,
		)

		if err != nil {
			r.log.Err("IN [GetAll] failed to scan dir ->", err)
		}
		res = append(res, t)
	}

	return
}

// Know if a dir exists by uuid
func (r *postgresDirRepository) ExistsByUuid(ctx context.Context, uuid string) (res bool) {
	query := `SELECT COUNT(*) > 0 FROM dir WHERE uuid = $1`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ExistsByUuid] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, uuid).Scan(&res)
	if err != nil {
		r.log.Err("IN [ExistsByUuid] failed to exec statement ->", err)
	}

	return
}

// Get dir by uuid
func (r *postgresDirRepository) GetByUuid(ctx context.Context, uuid string) (res domain.Dir, err error) {
	query := `SELECT * FROM dir WHERE uuid = $1`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [GetByUuid] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, uuid).Scan(
		&res.Uuid,
		&res.Name,
		&res.ParentDir,
	)

	if err != nil {
		r.log.Err("IN [GetByUuid] failed to scan rows ->", err)
		return
	}

	return
}

// Store a new dir
func (r *postgresDirRepository) Store(ctx context.Context, d *domain.Dir) (uuid string, err error) {
	query :=
		`INSERT INTO dir (name, parent_dir)
		VALUES ($1, $2)
		RETURNING uuid`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [Store] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(
		ctx,
		d.Name,
		d.ParentDir,
	).Scan(&uuid)

	if err != nil {
		r.log.Err("IN [Store] failed to scan rows ->", err)
		return
	}

	return
}

// Delete dir by uuid
func (r *postgresDirRepository) Delete(ctx context.Context, uuid string) (err error) {
	tx, err := r.Conn.Begin()
	if err != nil {
		r.log.Err("IN [Delete] failed to begin transaction ->", err)
	}
	defer tx.Rollback()

	// Pfile delete
	pfileQuery := `DELETE FROM pfile WHERE dir = $1`
	pfileStmt, err := r.Conn.PrepareContext(ctx, pfileQuery)
	if err != nil {
		r.log.Err("IN [Delete] failed to prepare PFILE context ->", err)
		return
	}
	defer pfileStmt.Close()

	_, err = pfileStmt.ExecContext(ctx, uuid)
	if err != nil {
		r.log.Err("IN [Delete] failed to exec PFILE statement ->", err)
		return
	}

	// Task delete
	taskQuery := `DELETE FROM task WHERE dir = $1`
	taskStmt, err := r.Conn.PrepareContext(ctx, taskQuery)
	if err != nil {
		r.log.Err("IN [Delete] failed to prepare TASK context ->", err)
		return
	}
	defer taskStmt.Close()

	_, err = taskStmt.ExecContext(ctx, uuid)
	if err != nil {
		r.log.Err("IN [Delete] failed to exec TASK statement ->", err)
		return
	}

	// Plan delete
	planQuery := `DELETE FROM plan WHERE dir = $1`
	planStmt, err := r.Conn.PrepareContext(ctx, planQuery)
	if err != nil {
		r.log.Err("IN [Delete] failed to prepare PLAN context ->", err)
		return
	}
	defer planStmt.Close()

	_, err = planStmt.ExecContext(ctx, uuid)
	if err != nil {
		r.log.Err("IN [Delete] failed to exec PLAN statement ->", err)
		return
	}

	// Permission delete
	permissionQuery := `DELETE FROM permission WHERE dir_uuid = $1`
	permissionStmt, err := r.Conn.PrepareContext(ctx, permissionQuery)
	if err != nil {
		r.log.Err("IN [Delete] failed to prepare PERMISSON context ->", err)
		return
	}
	defer permissionStmt.Close()

	_, err = permissionStmt.ExecContext(ctx, uuid)
	if err != nil {
		r.log.Err("IN [Delete] failed to exec PERMISSON statement ->", err)
		return
	}

	// Delete the dir itself
	dirQuery := `DELETE FROM dir WHERE uuid = $1`
	dirStmt, err := r.Conn.PrepareContext(ctx, dirQuery)
	if err != nil {
		r.log.Err("IN [Delete] failed to prepare DIR context ->", err)
		return
	}
	defer dirStmt.Close()

	_, err = dirStmt.ExecContext(ctx, uuid)
	if err != nil {
		r.log.Err("IN [Delete] failed to exec DIR statement ->", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		r.log.Err("IN [Delete] failed to commit changes -> ", err)
		return err
	}

	return
}

func (r *postgresDirRepository) ChgName(ctx context.Context, uuid string, nName string) (err error) {
	query := `UPDATE dir SET name = $1 WHERE uuid = $2`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ChgName] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, nName, uuid)
	if err != nil {
		r.log.Err("IN [ChgName] failed to exec statement ->", err)
		return
	}

	return
}

func (r *postgresDirRepository) ChgParentDir(ctx context.Context, uuid string, nPUuid string) (err error) {
	query := `UPDATE dir SET parent_dir=$1 WHERE uuid=$2`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ChgParentDir] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, nPUuid, uuid)
	if err != nil {
		r.log.Err("IN [ChgParentDir] failed to exec statement ->", err)
	}

	return
}

func (r *postgresDirRepository) InsertDirs(ctx context.Context, dirs []domain.Dir) (err error) {
	defaultErr := errors.New("Failed to insert new directories")
	query :=
		`INSERT INTO dir (uuid, name, parent_dir)
		VALUES ($1, $2, $3)`

	tx, err := r.Conn.Begin()
	if err != nil {
		r.log.Err("IN [InsertDirs] failed to begin transaction -> ", err)
		return defaultErr
	}
	defer tx.Rollback()

	for _, d := range dirs {
		_, err = tx.ExecContext(
			ctx,
			query,
			d.Uuid,
			d.Name,
			d.ParentDir,
		)

		if err != nil {
			r.log.Err("IN [InsertDirs] failed insert dir ->", err)
			return defaultErr
		}
	}

	err = tx.Commit()
	if err != nil {
		r.log.Err("IN [InsertDirs] failed to commit changes -> ", err)
		return defaultErr
	}

	return
}
