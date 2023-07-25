package postgres

import (
	"context"
	"database/sql"

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
	// TODO: Verify deletion
	query :=
		`DELETE FROM user_ WHERE username=$1`

	_, err = r.fetch(ctx, query, uname)

	return
}

// Change user state
func (r *postgresUserRepository) ChangeState(ctx context.Context, uname string, st domain.UserState) (err error) {
	query := `UPDATE user_ SET state=$1 WHERE username=$2`

	_, err = r.fetch(ctx, query, st.Code, uname)

	return
}

// Change user state
func (r *postgresUserRepository) ChangeRole(ctx context.Context, uname string, ro domain.Role) (err error) {
	query := `UPDATE user_ SET role=$1 WHERE username=$2`

	res, err := r.fetch(ctx, query, ro.Code, uname)
	domain.AgLog.Info("PATCH RES:", res)

	return
}
