package domain

import (
	"context"

	"github.com/sicozz/papyrus/domain/dtos"
)

// User represents the User data struct
type User struct {
	Uuid     string    `json:"uuid"`
	Username string    `json:"username" validate:"required,ascii"`
	Email    string    `json:"email" validate:"required,email,ascii"`
	Password string    `json:"password" validate:"required,ascii"`
	Name     string    `json:"name" validate:"required,ascii"`
	Lastname string    `json:"lastname" validate:"required,ascii"`
	Role     Role      `json:"role"`
	State    UserState `json:"state"`
}

type Permission struct {
	UserUuid string
	DirUuid  string
}

// UserUsecase represents the user's usecases
type UserUsecase interface {
	GetAll(c context.Context) ([]dtos.UserGetDto, RequestErr)
	GetByUsername(c context.Context, uname string) (dtos.UserGetDto, RequestErr)
	Login(c context.Context, uname string, passwd string) (dtos.UserGetDto, RequestErr)
	Store(c context.Context, p dtos.UserStore) (dtos.UserGetDto, RequestErr)
	Delete(c context.Context, uname string) RequestErr
	Update(c context.Context, uuid string, d dtos.UserUpdateDto) RequestErr
	ChgPasswd(c context.Context, uuid string, data dtos.UserChgPasswdDto) RequestErr

	GetUserPermittedDirs(c context.Context, uUuid string) ([]dtos.DirGetDto, RequestErr)
	AddPermission(c context.Context, uUuid, dUuid string) RequestErr
	RevokePermission(c context.Context, uUuid string, dUuid string) RequestErr
}

// UserRepository represents the user's repository contract
type UserRepository interface {
	GetAll(ctx context.Context) ([]User, error)
	GetByUuid(ctx context.Context, uuid string) (User, error)
	GetByUsername(ctx context.Context, uname string) (User, error)
	ExistsByUuid(ctx context.Context, uuid string) bool
	ExistsByUname(ctx context.Context, uname string) bool
	ExistsByEmail(ctx context.Context, email string) bool
	Auth(ctx context.Context, uname string, passwd string) bool
	Store(ctx context.Context, u *User) error
	Delete(ctx context.Context, uuid string) error
	ChgUsername(ctx context.Context, uname string, nUname string) error
	ChgEmail(ctx context.Context, uname string, email string) error
	ChgName(ctx context.Context, uname string, nName string) error
	ChgLstname(ctx context.Context, uname string, nLname string) error
	ChgRole(ctx context.Context, uname string, ro Role) error
	ChgState(ctx context.Context, uname string, st UserState) error
	ChgPasswd(ctx context.Context, uuid string, nPasswd string) error

	// Transactions
	Update(ctx context.Context, uuid string, p dtos.UserUpdateDto) error

	GetPermissionsByUserUuid(c context.Context, uUuid string) ([]Permission, error)
	ExistsPermission(c context.Context, uUuid, dUuid string) bool
	AddPermission(c context.Context, p Permission) error
	RevokePermission(c context.Context, uUuid string, dUuid string) error
}
