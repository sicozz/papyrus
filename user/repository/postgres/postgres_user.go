package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/utils"
	"github.com/sicozz/papyrus/utils/constants"
)

type postgresUserRepository struct {
	Conn *sql.DB
	log  utils.AggregatedLogger
}

// NewPostgresUserRepository will create an object that represent the UserRepository interface
func NewPostgresUserRepository(conn *sql.DB) domain.UserRepository {
	logger := utils.NewAggregatedLogger(constants.Repository, constants.User)
	return &postgresUserRepository{conn, logger}
}

// Retrieve all users
func (r *postgresUserRepository) GetAll(ctx context.Context) (res []domain.User, err error) {
	query :=
		`SELECT
			uuid,
			username,
			email,
			name,
			lastname,
			r.code as role_code,
			r.description as role_description,
			s.code as state_code,
			s.description as state_description
		FROM
			user_ u
			INNER JOIN role r ON(u.role = r.code)
			INNER JOIN user_state s ON(u.state = s.code)`

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

	res = make([]domain.User, 0)
	for rows.Next() {
		t := domain.User{}
		err = rows.Scan(
			&t.Uuid,
			&t.Username,
			&t.Email,
			&t.Name,
			&t.Lastname,
			&t.Role.Code,
			&t.Role.Description,
			&t.State.Code,
			&t.State.Description,
		)

		// TODO: Rename msgs with 'could not' to 'failed to'
		if err != nil {
			r.log.Err("IN [GetAll]: failed to scan user ->", err)
			return nil, err
		}

		res = append(res, t)
	}

	return
}

func (r *postgresUserRepository) GetByUuid(ctx context.Context, uuid string) (res domain.User, err error) {
	query :=
		`SELECT
			uuid,
			username,
			email,
			name,
			lastname,
			r.code as role_code,
			r.description as role_description,
			s.code as state_code,
			s.description as state_description
		FROM
			user_ u
			INNER JOIN role r ON(u.role = r.code)
			INNER JOIN user_state s ON(u.state = s.code)
		WHERE u.uuid = $1`

	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [GetByUuid] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, uuid).Scan(
		&res.Uuid,
		&res.Username,
		&res.Email,
		&res.Name,
		&res.Lastname,
		&res.Role.Code,
		&res.Role.Description,
		&res.State.Code,
		&res.State.Description,
	)

	if err != nil {
		r.log.Err("IN [GetByUuid] failed to scan row ->", err)
		return
	}

	return
}

func (r *postgresUserRepository) GetByUsername(ctx context.Context, uname string) (res domain.User, err error) {
	query :=
		`SELECT
			uuid,
			username,
			email,
			name,
			lastname,
			r.code as role_code,
			r.description as role_description,
			s.code as state_code,
			s.description as state_description
		FROM
			user_ u
			INNER JOIN role r ON(u.role = r.code)
			INNER JOIN user_state s ON(u.state = s.code)
		WHERE u.username = $1`

	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [GetByUsername] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, uname).Scan(
		&res.Uuid,
		&res.Username,
		&res.Email,
		&res.Name,
		&res.Lastname,
		&res.Role.Code,
		&res.Role.Description,
		&res.State.Code,
		&res.State.Description,
	)

	if err != nil {
		r.log.Err("IN [GetByUsername] failed to scan row ->", err)
		return
	}

	return
}

func (r *postgresUserRepository) ExistByUname(ctx context.Context, uname string) (res bool) {
	query := `SELECT COUNT(*) > 0 FROM user_ WHERE username = $1`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ExistByUname] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, uname).Scan(&res)

	return
}

func (r *postgresUserRepository) ExistByUuid(ctx context.Context, uuid string) (res bool) {
	query := `SELECT COUNT(*) > 0 FROM user_ WHERE uuid = $1`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ExistByUuid] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, uuid).Scan(&res)

	return
}

func (r *postgresUserRepository) ExistByEmail(ctx context.Context, email string) (res bool) {
	query := `SELECT COUNT(*) > 0 FROM user_ WHERE email = $1`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ExistByEmail] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, email).Scan(&res)

	return
}

// Store a new user
func (r *postgresUserRepository) Store(ctx context.Context, u *domain.User) (err error) {
	query :=
		`INSERT INTO user_ (username, email, password, name, lastname, role, state)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING uuid`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [Store] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(
		ctx,
		u.Username,
		u.Email,
		u.Password,
		u.Name,
		u.Lastname,
		u.Role.Code,
		u.State.Code,
	).Scan(&u.Uuid)

	if err != nil {
		r.log.Err("IN [Store] failed to scan rows ->", err)
		return
	}

	return
}

// Delete a user
func (r *postgresUserRepository) Delete(ctx context.Context, uname string) (err error) {
	query := `DELETE FROM user_ WHERE username=$1 RETURNING uuid`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [Delete] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	var uuid string
	err = stmt.QueryRowContext(ctx, uname).Scan(&uuid)
	if uuid == "" && err == nil {
		err = errors.New("Failed to delete user")
	}

	return
}

// Change user email
func (r *postgresUserRepository) ChgEmail(ctx context.Context, uname string, nEmail string) (err error) {
	query := `UPDATE user_ SET email=$1 WHERE username=$2`

	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ChgEmail] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, nEmail, uname)

	return
}

// Change user name
func (r *postgresUserRepository) ChgName(ctx context.Context, uname string, nName string) (err error) {
	query := `UPDATE user_ SET name=$1 WHERE username=$2`

	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ChgName] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, nName, uname)

	return
}

// Change user lastname
func (r *postgresUserRepository) ChgLstname(ctx context.Context, uname string, nLname string) (err error) {
	query := `UPDATE user_ SET lastname=$1 WHERE username=$2`

	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ChgLastname] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, nLname, uname)

	return
}

func (r *postgresUserRepository) ChgUsername(ctx context.Context, uname string, nUname string) (err error) {
	query := `UPDATE user_ SET username=$1 WHERE username=$2`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ChgUsername] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.QueryContext(ctx, nUname, uname)

	return
}

// Change user password
func (r *postgresUserRepository) ChgPasswd(ctx context.Context, uuid string, nPasswd string) (err error) {
	query := `UPDATE user_ SET lastname=$1 WHERE username=$2`

	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ChgLastname] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, nPasswd, uuid)

	return
}

// Change user role
func (r *postgresUserRepository) ChgRole(ctx context.Context, uname string, ro domain.Role) (err error) {
	query := `UPDATE user_ SET lastname=$1 WHERE username=$2`

	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ChgRole] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, ro.Code, uname)

	return
}

// Change user state
func (r *postgresUserRepository) ChgState(ctx context.Context, uname string, st domain.UserState) (err error) {
	query := `UPDATE user_ SET state=$1 WHERE username=$2`

	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ChgRole] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, st.Code, uname)

	return
}

// Authenticate a user
func (r *postgresUserRepository) Auth(ctx context.Context, uname string, passwd string) (res bool) {
	query := `SELECT COUNT(*) > 0 FROM user_ WHERE username = $1 AND password = $2`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [Login] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, uname, passwd).Scan(&res)

	return
}
