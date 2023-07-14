package postgres

import (
	"context"
	"database/sql"

	"github.com/sicozz/papyrus/domain"
)

type postgresUserRepository struct {
	Conn *sql.DB
}

// NewPostgresUserRepository will cerate an object that represent the user.Repository interface
func NewPostgresUserRepository(conn *sql.DB) domain.UserRepository {
	return &postgresUserRepository{conn}
}

func (r *postgresUserRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.User, err error) {
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

	result = make([]domain.User, 0)
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
		result = append(result, t)
	}

	return result, nil
}

// Retrieve all users
func (r *postgresUserRepository) Fetch(ctx context.Context) (res []domain.User, err error) {
	// query := `SELECT uuid, username, email, name, lastname, role_id, state_id`
	query := `SELECT uuid, username, email, name, lastname, role, state FROM user_`

	res, err = r.fetch(ctx, query)
	if err != nil {
		return nil, err
	}

	return
}
