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

type postgresUserRepository struct {
	Conn *sql.DB
	log  utils.AggregatedLogger
}

// NewPostgresUserRepository will create an object that represent the UserRepository interface
func NewPostgresUserRepository(conn *sql.DB) domain.UserRepository {
	logger := utils.NewAggregatedLogger(constants.Repository, constants.User)
	return &postgresUserRepository{conn, logger}
}

func (r *postgresUserRepository) fetch(ctx context.Context, query string, args ...interface{}) (res []domain.User, err error) {
	rows, err := r.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		r.log.Err(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			r.log.Err(errRow)
		}
	}()

	res = make([]domain.User, 0)
	for rows.Next() {
		t := domain.User{}
		roleCode := int64(0)
		stateCode := int64(0)
		// Get from db
		err = rows.Scan(
			&t.Uuid,
			&t.Username,
			&t.Email,
			&t.Password,
			&t.Name,
			&t.Lastname,
			&roleCode,
			&stateCode,
		)

		if err != nil {
			r.log.Err("IN [fetch]:", err)
			return nil, err
		}
		t.Role = domain.Role{
			Code: roleCode,
		}
		t.State = domain.UserState{
			Code: stateCode,
		}
		res = append(res, t)
	}

	return res, nil
}

// Retrieve all users
func (r *postgresUserRepository) GetAll(ctx context.Context) (res []domain.User, err error) {
	query :=
		`SELECT uuid, username, email, password, name, lastname, role, state
		FROM user_`

	res, err = r.fetch(ctx, query)
	if err != nil {
		return nil, err
	}
	for _, u := range res {
		u.Password = ""
	}

	return
}

func (r *postgresUserRepository) GetByUuid(ctx context.Context, uuid string) (res domain.User, err error) {
	// TODO: Add role and state columns
	query :=
		`SELECT uuid, username, email, password, name, lastname
		FROM user_
		WHERE uuid = $1`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [GetByUuid]: could not prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, uuid).Scan(
		&res.Uuid,
		&res.Username,
		&res.Email,
		&res.Password,
		&res.Name,
		&res.Lastname,
	)

	if err != nil {
		r.log.Err("IN [GetByUuid]: could not scan rows ->", err)
		return
	}

	return
}

func (r *postgresUserRepository) GetByUsername(ctx context.Context, uname string) (res domain.User, err error) {
	// TODO: Refactor operations that expect only 1 row
	query :=
		`SELECT uuid, username, email, password, name, lastname, role, state
		FROM user_
		WHERE username = $1`

	users, err := r.fetch(ctx, query, uname)
	if err != nil {
		return domain.User{}, err
	}

	if len(users) < 1 {
		return domain.User{}, errors.New(fmt.Sprintln("No user with username:", uname))
	}

	for _, u := range users {
		u.Password = ""
	}

	res = users[0]

	return
}

func (r *postgresUserRepository) ExistByUname(ctx context.Context, uname string) (res bool) {
	query := `SELECT COUNT(*) > 0 FROM user_ WHERE username = $1`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ExistByUname]: could not prepare context ->", err)
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
		r.log.Err("IN [ExistByUuid]: could not prepare context ->", err)
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
		r.log.Err("IN [ExistByEmail]: could not prepare context ->", err)
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
		r.log.Err("IN [Store]: could not prepare context ->", err)
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
		r.log.Err("IN [Store]: could not scan rows ->", err)
		return
	}

	return
}

// Delete a user
func (r *postgresUserRepository) Delete(ctx context.Context, uname string) (err error) {
	query := `DELETE FROM user_ WHERE username=$1 RETURNING uuid`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [Delete]: could not prepare context ->", err)
		return
	}
	defer stmt.Close()

	var uuid string
	err = stmt.QueryRowContext(ctx, uname).Scan(&uuid)
	if uuid == "" && err == nil {
		err = errors.New("Could not delete user")
	}

	return
}

// Change user email
func (r *postgresUserRepository) ChgEmail(ctx context.Context, uname string, nEmail string) (err error) {
	query := `UPDATE user_ SET email=$1 WHERE username=$2`

	_, err = r.fetch(ctx, query, nEmail, uname)

	return
}

// Change user name
func (r *postgresUserRepository) ChgName(ctx context.Context, uname string, nName string) (err error) {
	query := `UPDATE user_ SET name=$1 WHERE username=$2`

	_, err = r.fetch(ctx, query, nName, uname)

	return
}

// Change user lastname
func (r *postgresUserRepository) ChgLstname(ctx context.Context, uname string, nLname string) (err error) {
	query := `UPDATE user_ SET lastname=$1 WHERE username=$2`

	_, err = r.fetch(ctx, query, nLname, uname)

	return
}

func (r *postgresUserRepository) ChgUsername(ctx context.Context, uname string, nUname string) (err error) {
	query := `UPDATE user_ SET username=$1 WHERE username=$2`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ChgUsername]: could not prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.QueryContext(ctx, nUname, uname)

	return
}

// Change user password
func (r *postgresUserRepository) ChgPasswd(ctx context.Context, uuid string, nPasswd string) (err error) {
	query := `UPDATE user_ SET password=$1 WHERE uuid=$2`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ChgPasswd]: could not prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.QueryContext(ctx, nPasswd, uuid)

	return
}

// Change user role
func (r *postgresUserRepository) ChgRole(ctx context.Context, uname string, ro domain.Role) (err error) {
	query := `UPDATE user_ SET role=$1 WHERE username=$2`

	_, err = r.fetch(ctx, query, ro.Code, uname)

	return
}

// Change user state
func (r *postgresUserRepository) ChgState(ctx context.Context, uname string, st domain.UserState) (err error) {
	query := `UPDATE user_ SET state=$1 WHERE username=$2`

	_, err = r.fetch(ctx, query, st.Code, uname)

	return
}

// Authenticate a user
func (r *postgresUserRepository) Login(ctx context.Context, uname string, passwd string) (res domain.User, err error) {
	user, err := r.GetByUsername(ctx, uname)
	if err != nil {
		return domain.User{}, err
	}

	if user.Password != passwd {
		return domain.User{}, errors.New("Incorrect password or username")
	}

	return user, nil
}
