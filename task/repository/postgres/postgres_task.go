package postgres

import (
	"context"
	"database/sql"

	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/utils"
	"github.com/sicozz/papyrus/utils/constants"
)

type postgresTaskRepository struct {
	Conn *sql.DB
	log  utils.AggregatedLogger
}

// NewPostgresTaskRepository will create an object that represent the TaskRepository interface
func NewPostgresTaskRepository(conn *sql.DB) domain.TaskRepository {
	logger := utils.NewAggregatedLogger(constants.Repository, constants.Task)
	return &postgresTaskRepository{conn, logger}
}

func (r *postgresTaskRepository) GetAll(ctx context.Context) (res []domain.Task, err error) {
	query :=
		`SELECT
			uuid,
			name,
			procedure,
			date_creation,
			term,
			ts.description AS state,
			dir,
			creator_user,
			recv_user,
			chk
		FROM
			task t
			INNER JOIN task_state ts ON(t.state = ts.code)`

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

	res = make([]domain.Task, 0)
	for rows.Next() {
		t := domain.Task{}
		err = rows.Scan(
			&t.Uuid,
			&t.Name,
			&t.Procedure,
			&t.DateCreation,
			&t.Term,
			&t.State,
			&t.Dir,
			&t.CreatorUser,
			&t.RecvUser,
			&t.Check,
		)

		if err != nil {
			r.log.Err("IN [GetAll] failed to scan dir ->", err)
		}
		res = append(res, t)
	}

	return
}

func (r *postgresTaskRepository) GetByUuid(ctx context.Context, uuid string) (res domain.Task, err error) {
	query :=
		`SELECT
			uuid,
			name,
			procedure,
			date_creation,
			term,
			ts.description AS state,
			dir,
			creator_user,
			recv_user,
			chk
		FROM
			task t
			INNER JOIN task_state ts ON(t.state = ts.code)
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
		&res.Name,
		&res.Procedure,
		&res.DateCreation,
		&res.Term,
		&res.State,
		&res.Dir,
		&res.CreatorUser,
		&res.RecvUser,
		&res.Check,
	)

	if err != nil {
		r.log.Err("IN [GetByUuid] failed to scan row ->", err)
		return
	}

	return
}

func (r *postgresTaskRepository) GetByUser(ctx context.Context, uuid string) (res []domain.Task, err error) {
	query :=
		`SELECT
			uuid,
			name,
			procedure,
			date_creation,
			term,
			ts.description AS state,
			dir,
			creator_user,
			recv_user,
			chk
		FROM
			task t
			INNER JOIN task_state ts ON(t.state = ts.code)
		WHERE
			recv_user = $1`

	rows, err := r.Conn.QueryContext(ctx, query, uuid)
	if err != nil {
		res = nil
		return
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			r.log.Err("IN [GetByUser] failed to close *rows ->", err)
		}
	}()

	res = make([]domain.Task, 0)
	for rows.Next() {
		t := domain.Task{}
		err = rows.Scan(
			&t.Uuid,
			&t.Name,
			&t.Procedure,
			&t.DateCreation,
			&t.Term,
			&t.State,
			&t.Dir,
			&t.CreatorUser,
			&t.RecvUser,
			&t.Check,
		)

		if err != nil {
			r.log.Err("IN [GetByUser] failed to scan dir ->", err)
		}
		res = append(res, t)
	}

	return
}

func (r *postgresTaskRepository) ExistsByUuid(ctx context.Context, uuid string) (res bool) {
	query := `SELECT COUNT(*) > 0 FROM task WHERE uuid = $1`
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

func (r *postgresTaskRepository) ChgCheck(ctx context.Context, tUuid string, uUuid string, chk bool) (err error) {
	query := `UPDATE task SET chk = $1 WHERE uuid = $2 AND recv_user = $3`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ChgCheck] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, chk, tUuid, uUuid)
	if err != nil {
		r.log.Err("IN [ChgCheck] failed to exec statement ->", err)
		return
	}

	return
}

func (r *postgresTaskRepository) ChgState(ctx context.Context, tUuid string, uUuid string, desc string) (err error) {
	query :=
		`UPDATE task AS t
		SET state = ts.code
		FROM task_state AS ts
		WHERE
			ts.description = $3 AND
			t.uuid = $1 AND
			t.creator_user = $2`

	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ChgState] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, tUuid, uUuid, desc)
	if err != nil {
		r.log.Err("IN [ChgState] failed to exec statement ->", err)
		return
	}

	return
}

func (r *postgresTaskRepository) Store(ctx context.Context, t domain.Task) (uuid string, err error) {
	query :=
		`INSERT INTO task (
			name,
			procedure,
			date_creation,
			term,
			state,
			dir,
			creator_user,
			recv_user,
			chk
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
			$9
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
		t.Name,
		t.Procedure,
		t.DateCreation,
		t.Term,
		1,
		t.Dir,
		t.CreatorUser,
		t.RecvUser,
		false,
	).Scan(&uuid)

	if err != nil {
		r.log.Err("IN [Store] failed to scan rows ->", err)
		return
	}

	return
}

func (r *postgresTaskRepository) Delete(ctx context.Context, tUuid, uUuid string) (err error) {
	query := `DELETE FROM task WHERE uuid = $1 AND creator_user = $2`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [Delete] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, tUuid, uUuid)
	if err == nil {
		r.log.Err("IN [Delete] failed to exec statement ->", err)
		return
	}

	return
}

func (r *postgresTaskRepository) ExistsStateByDesc(ctx context.Context, desc string) (res bool) {
	query := `SELECT COUNT(*) > 0 FROM task_state WHERE description = $1`
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
