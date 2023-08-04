package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sicozz/papyrus/domain"
)

type postgresUserRepository struct {
	Conn *sql.DB
}

// NewPostgresUserRepository will create an object that represent the UserRepository interface
func NewPostgresUserRepository(conn *sql.DB) domain.UserRepository {
	return &postgresUserRepository{conn}
}

func (r *postgresUserRepository) fetch(ctx context.Context, query string, args ...interface{}) (res []domain.User, err error) {
	rows, err := r.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		domain.AgLog.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			domain.AgLog.Error(errRow)
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
			domain.AgLog.Error(err)
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
func (r *postgresUserRepository) Fetch(ctx context.Context) (res []domain.User, err error) {
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

// Get user by id
func (r *postgresUserRepository) GetByUsername(ctx context.Context, uname string) (res domain.User, err error) {
	// TODO: Refactor operations that expect only 1 row
	// TODO: Rename this function to GetByUsername
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

// Store a new user
func (r *postgresUserRepository) Store(ctx context.Context, u *domain.User) (err error) {
	query :=
		`INSERT INTO user_ (username, email, password, name, lastname, role, state)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING uuid`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		domain.AgLog.Error(err)
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

	return
}

// Delete a user
func (r *postgresUserRepository) Delete(ctx context.Context, uname string) (err error) {
	query := `DELETE FROM user_ WHERE username=$1 RETURNING uuid`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		domain.AgLog.Error(err)
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
