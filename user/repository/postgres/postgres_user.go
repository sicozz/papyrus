package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
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

func (r *postgresUserRepository) ExistsByUname(ctx context.Context, uname string) (res bool) {
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

func (r *postgresUserRepository) ExistsByUuid(ctx context.Context, uuid string) (res bool) {
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

func (r *postgresUserRepository) ExistsByEmail(ctx context.Context, email string) (res bool) {
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
func (r *postgresUserRepository) ChgPasswd(ctx context.Context, uname string, nPasswd string) (err error) {
	query := `UPDATE user_ SET password=$1 WHERE username=$2`

	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ChgLastname] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, nPasswd, uname)

	return
}

// Change user role
func (r *postgresUserRepository) ChgRole(ctx context.Context, uname string, ro domain.Role) (err error) {
	query := `UPDATE user_ SET role=$1 WHERE username=$2`

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
		r.log.Err("IN [ChgState] failed to prepare context ->", err)
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
		r.log.Err("IN [Auth] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, uname, passwd).Scan(&res)

	return
}

// Update user data
func (r *postgresUserRepository) Update(ctx context.Context, uuid string, p dtos.UserUpdateDto) (err error) {
	defaultErr := errors.New("Failed to update user data")
	tx, err := r.Conn.Begin()
	if err != nil {
		r.log.Err("IN [Update] failed to begin transaction -> ", err)
		return defaultErr
	}
	defer tx.Rollback()

	// User.Username
	if p.Username != "" {
		usernameQuery := `UPDATE user_ SET username = $1 WHERE uuid = $2`
		_, err = tx.ExecContext(ctx, usernameQuery, p.Username, uuid)

		if err != nil {
			r.log.Err("IN [Update] failed to update username -> ", err)
			return defaultErr
		}
	}

	// User.Name
	if p.Name != "" {
		nameQuery := `UPDATE user_ SET name = $1 WHERE uuid = $2`
		_, err = tx.ExecContext(ctx, nameQuery, p.Name, uuid)

		if err != nil {
			r.log.Err("IN [Update] failed to update name -> ", err)
			return defaultErr
		}
	}

	// User.Lastname
	if p.Lastname != "" {
		lastnameQuery := `UPDATE user_ SET lastname = $1 WHERE uuid = $2`
		_, err = tx.ExecContext(ctx, lastnameQuery, p.Lastname, uuid)

		if err != nil {
			r.log.Err("IN [Update] failed to update lastname -> ", err)
			return defaultErr
		}
	}

	// User.Email
	if p.Email != "" {
		lastnameQuery := `UPDATE user_ SET email = $1 WHERE uuid = $2`
		_, err = tx.ExecContext(ctx, lastnameQuery, p.Email, uuid)

		if err != nil {
			r.log.Err("IN [Update] failed to update email -> ", err)
			return defaultErr
		}
	}

	// User.Role
	if p.Role != "" {
		var roleCode int64
		getRoleCodeQuery := `SELECT code FROM role WHERE description = $1`
		err = tx.QueryRowContext(ctx, getRoleCodeQuery, p.Role).Scan(&roleCode)

		if err != nil {
			r.log.Err("IN [Update] failed to get role code -> ", err)
			return defaultErr
		}

		roleQuery := `UPDATE user_ SET role = $1 WHERE uuid = $2`
		_, err = tx.ExecContext(ctx, roleQuery, roleCode, uuid)

		if err != nil {
			r.log.Err("IN [Update] failed to update role -> ", err)
			return defaultErr
		}
	}

	// User.State
	if p.State != "" {
		var stateCode int64
		getStateCodeQuery := `SELECT code FROM user_state WHERE description = $1`
		err = tx.QueryRowContext(ctx, getStateCodeQuery, p.State).Scan(&stateCode)

		if err != nil {
			r.log.Err("IN [Update] failed to get user_state code -> ", err)
			return defaultErr
		}

		stateQuery := `UPDATE user_ SET state = $1 WHERE uuid = $2`
		_, err = tx.ExecContext(ctx, stateQuery, stateCode, uuid)

		if err != nil {
			r.log.Err("IN [Update] failed to update state -> ", err)
			return defaultErr
		}
	}

	err = tx.Commit()
	if err != nil {
		r.log.Err("IN [Update] failed to commit changes -> ", err)
		return defaultErr
	}

	return
}

func (r *postgresUserRepository) ExistsPermission(ctx context.Context, uUuid, dUuid string) (res bool) {
	query := `SELECT COUNT(*) > 0 FROM permission WHERE user_uuid = $1 AND dir_uuid = $2`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [ExistPermission] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, uUuid, dUuid).Scan(&res)

	return
}

func (r *postgresUserRepository) AddPermission(ctx context.Context, p domain.Permission) (err error) {
	query :=
		`INSERT INTO permission (user_uuid, dir_uuid)
		VALUES ($1, $2)`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [AddPermission] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		p.UserUuid,
		p.DirUuid,
	)

	if err != nil {
		r.log.Err("IN [AddPermission] failed to scan rows ->", err)
		return
	}

	return
}

func (r *postgresUserRepository) RevokePermission(ctx context.Context, uUuid, dUuid string) (err error) {
	query := `DELETE FROM permission WHERE user_uuid = $1 AND dir_uuid = $2`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		r.log.Err("IN [RevokePermission] failed to prepare context ->", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, uUuid, dUuid)
	if err != nil {
		r.log.Err("IN [RevokePermission] failed to exec statement ->", err)
		return
	}

	return
}

func (r *postgresUserRepository) GetPermissionsByUserUuid(ctx context.Context, uUuid string) (res []domain.Permission, err error) {
	query := `SELECT user_uuid, dir_uuid FROM permission WHERE user_uuid = $1`

	rows, err := r.Conn.QueryContext(ctx, query, uUuid)
	if err != nil {
		res = nil
		return
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			r.log.Err("IN [GetPermissionsByUserUuid] failed to close *rows ->", err)
		}
	}()

	res = []domain.Permission{}
	for rows.Next() {
		t := domain.Permission{}
		err = rows.Scan(
			&t.UserUuid,
			&t.DirUuid,
		)

		if err != nil {
			r.log.Err("IN [GetPermissionsByUserUuid] failed to scan permission ->", err)
		}
		res = append(res, t)
	}

	return
}

func (r *postgresUserRepository) GetHistoryDownloads(ctx context.Context, uUuid string) (res []dtos.UserHistoryGetDto, err error) {
	query :=
		`SELECT
			dwn.uuid,
			dwn.date,
			dwn.user_,
			pf.uuid,
			pf.code,
			pf.version,
			pf.term,
			pf.name,
			pft.description AS pfile_type,
			pf.date_input,
			pf.dir
		FROM
			download dwn
			INNER JOIN pfile pf ON(dwn.pfile = pf.uuid)
			INNER JOIN pfile_type pft ON(pf.type = pft.code)
			INNER JOIN dir ON(dir.uuid = pf.dir)
		WHERE
			dwn.user_ = $1;`

	rows, err := r.Conn.QueryContext(ctx, query, uUuid)
	if err != nil {
		res = nil
		return
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			r.log.Err("IN [GetHistoryDownloads] failed to close *rows ->", err)
		}
	}()

	res = []dtos.UserHistoryGetDto{}
	for rows.Next() {
		t := dtos.UserHistoryGetDto{}
		err = rows.Scan(
			&t.DownloadUuid,
			&t.Date,
			&t.UserUuid,
			&t.PFileUuid,
			&t.PFileCode,
			&t.PFileVersion,
			&t.PFileTerm,
			&t.PFileName,
			&t.PFileType,
			&t.PFileDateInput,
			&t.PFileDir,
		)

		if err != nil {
			r.log.Err("IN [GetHistoryDownloads] failed to scan history download ->", err)
		}
		res = append(res, t)
	}

	return
}
