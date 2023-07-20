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
	query := `SELECT uuid, username, email, name, lastname, role, state FROM user_`

	res, err = r.fetch(ctx, query)
	if err != nil {
		return nil, err
	}

	return
}

func (r *postgresUserRepository) Store(ctx context.Context, u *domain.User) (err error) {
	query := `INSERT INTO user_ (username, email, password, name, lastname, role, state) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		domain.AgLog.Error(err)
		return
	}

	res, err := stmt.ExecContext(
		ctx,
		`Solarized`,
		`simozuluaga@gmail.com`,
		`hardHash21*`,
		`Simon`,
		`Zuluaga`,
		1,
		1,
	)
	// res, err := stmt.ExecContext(
	// 	ctx,
	// 	u.Uuid,
	// 	u.Username,
	// 	u.Email,
	// 	u.Password,
	// 	u.Name,
	// 	u.Lastname,
	// 	u.Role,
	// 	u.State,
	// )

	rowsAffected, err := res.RowsAffected()
	domain.AgLog.Warn("ROWS AFFECTED:\t", rowsAffected)
	if err != nil {
		return
	}

	// u.Uuid = lastUuid
	// TODO: defer stmt.Close() everywhere
	return
}
