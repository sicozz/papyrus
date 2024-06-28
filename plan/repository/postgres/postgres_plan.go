package postgres

import (
	"context"
	"database/sql"

	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/utils"
	"github.com/sicozz/papyrus/utils/constants"
)

type postgresPlanRepository struct {
	Conn *sql.DB
	log  utils.AggregatedLogger
}

// NewPostgresPlanRepository will create an object that represent the PlanRepository interface
func NewPostgresPlanRepository(conn *sql.DB) domain.PlanRepository {
	logger := utils.NewAggregatedLogger(constants.Repository, constants.Plan)
	return &postgresPlanRepository{conn, logger}
}

func (r *postgresPlanRepository) GetAll(ctx context.Context) (res []domain.Plan, err error) {
	query := `SELECT
			uuid,
			code,
			name,
			origin,
			action_type,
			term,
			creator_user,
			resp_user,
			date_create,
			date_close,
			causes,
			conclusions,
			state,
			stage,
			dir,
			action0_desc,
			action0_date,
			action0_user,
			action1_desc,
			action1_date,
			action1_user,
			action2_desc,
			action2_date,
			action2_user,
			action3_desc,
			action3_date,
			action3_user,
			action4_desc,
			action4_date,
			action4_user,
			action5_desc,
			action5_date,
			action5_user
		FROM plan p `

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

	res = make([]domain.Plan, 0)
	for rows.Next() {
		p := domain.Plan{}
		err = rows.Scan(
			&p.Uuid,
			&p.Code,
			&p.Name,
			&p.Origin,
			&p.ActionType,
			&p.Term,
			&p.CreatorUser,
			&p.RespUser,
			&p.DateCreation,
			&p.DateClose,
			&p.Causes,
			&p.Conclusions,
			&p.State,
			&p.Stage,
			&p.Dir,
			&p.Action0desc,
			&p.Action0date,
			&p.Action0user,
			&p.Action1desc,
			&p.Action1date,
			&p.Action1user,
			&p.Action2desc,
			&p.Action2date,
			&p.Action2user,
			&p.Action3desc,
			&p.Action3date,
			&p.Action3user,
			&p.Action4desc,
			&p.Action4date,
			&p.Action4user,
			&p.Action5desc,
			&p.Action5date,
			&p.Action5user,
		)

		if err != nil {
			r.log.Err("IN [GetAll] failed to scan plan ->", err)
		}
		res = append(res, p)
	}

	return
}

func (r *postgresPlanRepository) GetByUuid(ctx context.Context, uuid string) (res domain.Plan, err error) {
	query := `SELECT
			uuid,
			code,
			name,
			origin,
			action_type,
			term,
			creator_user,
			resp_user,
			date_create,
			date_close,
			causes,
			conclusions,
			state,
			stage,
			dir,
			action0_desc,
			action0_date,
			action0_user,
			action1_desc,
			action1_date,
			action1_user,
			action2_desc,
			action2_date,
			action2_user,
			action3_desc,
			action3_date,
			action3_user,
			action4_desc,
			action4_date,
			action4_user,
			action5_desc,
			action5_date,
			action5_user
		FROM plan p WHERE uuid = $1`
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
		&res.Origin,
		&res.ActionType,
		&res.Term,
		&res.CreatorUser,
		&res.RespUser,
		&res.DateCreation,
		&res.DateClose,
		&res.Causes,
		&res.Conclusions,
		&res.State,
		&res.Stage,
		&res.Dir,
		&res.Action0desc,
		&res.Action0date,
		&res.Action0user,
		&res.Action1desc,
		&res.Action1date,
		&res.Action1user,
		&res.Action2desc,
		&res.Action2date,
		&res.Action2user,
		&res.Action3desc,
		&res.Action3date,
		&res.Action3user,
		&res.Action4desc,
		&res.Action4date,
		&res.Action4user,
		&res.Action5desc,
		&res.Action5date,
		&res.Action5user,
	)

	if err != nil {
		r.log.Err("IN [GetByUuid] failed to scan rows ->", err)
		return
	}

	return
}

func (r *postgresPlanRepository) GetByUser(ctx context.Context, uuid string) (res []domain.Plan, err error) {
	query := `SELECT
			uuid,
			code,
			name,
			origin,
			action_type,
			term,
			creator_user,
			resp_user,
			date_create,
			date_close,
			causes,
			conclusions,
			state,
			stage,
			dir,
			action0_desc,
			action0_date,
			action0_user,
			action1_desc,
			action1_date,
			action1_user,
			action2_desc,
			action2_date,
			action2_user,
			action3_desc,
			action3_date,
			action3_user,
			action4_desc,
			action4_date,
			action4_user,
			action5_desc,
			action5_date,
			action5_user
		FROM plan p WHERE resp_user = $1`

	rows, err := r.Conn.QueryContext(ctx, query, uuid)
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

	res = make([]domain.Plan, 0)
	for rows.Next() {
		p := domain.Plan{}
		err = rows.Scan(
			&p.Uuid,
			&p.Code,
			&p.Name,
			&p.Origin,
			&p.ActionType,
			&p.Term,
			&p.CreatorUser,
			&p.RespUser,
			&p.DateCreation,
			&p.DateClose,
			&p.Causes,
			&p.Conclusions,
			&p.State,
			&p.Stage,
			&p.Dir,
			&p.Action0desc,
			&p.Action0date,
			&p.Action0user,
			&p.Action1desc,
			&p.Action1date,
			&p.Action1user,
			&p.Action2desc,
			&p.Action2date,
			&p.Action2user,
			&p.Action3desc,
			&p.Action3date,
			&p.Action3user,
			&p.Action4desc,
			&p.Action4date,
			&p.Action4user,
			&p.Action5desc,
			&p.Action5date,
			&p.Action5user,
		)

		if err != nil {
			r.log.Err("IN [GetAll] failed to scan plan ->", err)
		}
		res = append(res, p)
	}

	return
}

func (r *postgresPlanRepository) GetOwnedByUser(ctx context.Context, uuid string) (res []domain.Plan, err error) {
	query := `SELECT
			uuid,
			code,
			name,
			origin,
			action_type,
			term,
			creator_user,
			resp_user,
			date_create,
			date_close,
			causes,
			conclusions,
			state,
			stage,
			dir,
			action0_desc,
			action0_date,
			action0_user,
			action1_desc,
			action1_date,
			action1_user,
			action2_desc,
			action2_date,
			action2_user,
			action3_desc,
			action3_date,
			action3_user,
			action4_desc,
			action4_date,
			action4_user,
			action5_desc,
			action5_date,
			action5_user
		FROM plan p WHERE creator_user = $1`

	rows, err := r.Conn.QueryContext(ctx, query, uuid)
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

	res = make([]domain.Plan, 0)
	for rows.Next() {
		p := domain.Plan{}
		err = rows.Scan(
			&p.Uuid,
			&p.Code,
			&p.Name,
			&p.Origin,
			&p.ActionType,
			&p.Term,
			&p.CreatorUser,
			&p.RespUser,
			&p.DateCreation,
			&p.DateClose,
			&p.Causes,
			&p.Conclusions,
			&p.State,
			&p.Stage,
			&p.Dir,
			&p.Action0desc,
			&p.Action0date,
			&p.Action0user,
			&p.Action1desc,
			&p.Action1date,
			&p.Action1user,
			&p.Action2desc,
			&p.Action2date,
			&p.Action2user,
			&p.Action3desc,
			&p.Action3date,
			&p.Action3user,
			&p.Action4desc,
			&p.Action4date,
			&p.Action4user,
			&p.Action5desc,
			&p.Action5date,
			&p.Action5user,
		)

		if err != nil {
			r.log.Err("IN [GetAll] failed to scan plan ->", err)
		}
		res = append(res, p)
	}

	return
}

func (r *postgresPlanRepository) ExistsByUuid(ctx context.Context, uuid string) (res bool) {
	query := `SELECT COUNT(*) > 0 FROM plan WHERE uuid = $1`
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

func (r *postgresPlanRepository) Store(ctx context.Context, p domain.Plan) (uuid string, err error) {
	query := `INSERT INTO plan (
			code,
			name,
			origin,
			action_type,
			term,
			creator_user,
			resp_user,
			date_create,
			date_close,
			causes,
			conclusions,
			state,
			stage,
			dir,
			action0_desc,
			action0_date,
			action0_user,
			action1_desc,
			action1_date,
			action1_user,
			action2_desc,
			action2_date,
			action2_user,
			action3_desc,
			action3_date,
			action3_user,
			action4_desc,
			action4_date,
			action4_user,
			action5_desc,
			action5_date,
			action5_user
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
			$13,
			$14,
			$15,
			$16,
			$17,
			$18,
			$19,
			$20,
			$21,
			$22,
			$23,
			$24,
			$25,
			$26,
			$27,
			$28,
			$29,
			$30,
			$31,
			$32
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
		p.Code,
		p.Name,
		p.Origin,
		p.ActionType,
		p.Term,
		p.CreatorUser,
		p.RespUser,
		p.DateCreation,
		p.DateClose,
		p.Causes,
		p.Conclusions,
		p.State,
		p.Stage,
		p.Dir,
		p.Action0desc,
		p.Action0date,
		p.Action0user,
		p.Action1desc,
		p.Action1date,
		p.Action1user,
		p.Action2desc,
		p.Action2date,
		p.Action2user,
		p.Action3desc,
		p.Action3date,
		p.Action3user,
		p.Action4desc,
		p.Action4date,
		p.Action4user,
		p.Action5desc,
		p.Action5date,
		p.Action5user,
	).Scan(&uuid)

	if err != nil {
		r.log.Err("IN [Store] failed to scan rows ->", err)
		return
	}

	return
}

func (r *postgresPlanRepository) Update(ctx context.Context, uuid string, p domain.Plan) (err error) {
	query := `UPDATE plan SET
			name = $2,
			origin = $3,
			action_type = $4,
			term = $5,
			resp_user = $6,
			date_close = $7,
			causes = $8,
			conclusions = $9,
			dir = $10,
			state = $11,
			stage = $12,
			action0_desc = $13,
			action0_date = $14,
			action0_user = $15,
			action1_desc = $16,
			action1_date = $17,
			action1_user = $18,
			action2_desc = $19,
			action2_date = $20,
			action2_user = $21,
			action3_desc = $22,
			action3_date = $23,
			action3_user = $24,
			action4_desc = $25,
			action4_date = $26,
			action4_user = $27,
			action5_desc = $28,
			action5_date = $29,
			action5_user = $30
		WHERE
			uuid = $1`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [Update] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		uuid,
		p.Name,
		p.Origin,
		p.ActionType,
		p.Term,
		p.RespUser,
		p.DateClose,
		p.Causes,
		p.Conclusions,
		p.Dir,
		p.State,
		p.Stage,
		p.Action0desc,
		p.Action0date,
		p.Action0user,
		p.Action1desc,
		p.Action1date,
		p.Action1user,
		p.Action2desc,
		p.Action2date,
		p.Action2user,
		p.Action3desc,
		p.Action3date,
		p.Action3user,
		p.Action4desc,
		p.Action4date,
		p.Action4user,
		p.Action5desc,
		p.Action5date,
		p.Action5user,
	)

	if err != nil {
		r.log.Err("IN [Update] failed to scan rows ->", err)
		return
	}

	return
}

func (r *postgresPlanRepository) Delete(ctx context.Context, uuid string) (err error) {
	query := `DELETE FROM plan WHERE uuid = $1`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [Delete] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, uuid)
	if err != nil {
		r.log.Err("IN [Delete] failed to exec statement ->", err)
		return
	}

	return
}
