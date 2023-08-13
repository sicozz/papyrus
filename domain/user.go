package domain

import "context"

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

// UserUsecase represents the user's usecases
type UserUsecase interface {
	GetAll(c context.Context) ([]User, RequestErr)
	// GetByUuid(c context.Context, uuid string) (User, error)
	// GetByEmail(c context.Context, email string) (User, error)
	GetByUsername(c context.Context, uname string) (User, RequestErr)
	// TODO: Change *User recievers to *dto used in handler
	Store(c context.Context, u *User) RequestErr
	Update(c context.Context, uname string, uUp *User) RequestErr
	Delete(c context.Context, uname string) RequestErr
	Login(c context.Context, uname string, passwd string) (User, RequestErr)
}

// UserRepository represents the user's repository contract
type UserRepository interface {
	// TODO reorganize functions
	GetAll(ctx context.Context) ([]User, error)
	// GetByUuid(ctx context.Context, uuid string) (User, error)
	// GetByEmail(ctx context.Context, email string) (User, error)
	GetByUsername(ctx context.Context, uname string) (User, error)
	ExistByUname(ctx context.Context, uname string) bool
	ExistByEmail(ctx context.Context, email string) bool
	Store(ctx context.Context, u *User) error
	Delete(ctx context.Context, uuid string) error
	Login(ctx context.Context, uname string, passwd string) (User, error)
	ChgEmail(ctx context.Context, uname string, email string) error
	ChgName(ctx context.Context, uname string, nName string) error
	ChgLstname(ctx context.Context, uname string, nLname string) error
	ChgRole(ctx context.Context, uname string, ro Role) error
	ChgState(ctx context.Context, uname string, st UserState) error
}
