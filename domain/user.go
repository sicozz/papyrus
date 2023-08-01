package domain

import "context"

// User is representing the User data struct
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
	Fetch(c context.Context) ([]User, error)
	// GetByUuid(c context.Context, uuid string) (User, error)
	// GetByEmail(c context.Context, email string) (User, error)
	GetByUsername(c context.Context, uname string) (User, error)
	// Update(c context.Context, u *User) error
	Store(c context.Context, u *User) error
	Delete(c context.Context, uname string) (string /*uuid*/, error)
	Login(c context.Context, uname string, passwd string) (User, error)
	ChangeState(c context.Context, uname string, desc string) error
	ChangeRole(c context.Context, uname string, desc string) error
}

// UserRepository represents the user's repository contract
type UserRepository interface {
	Fetch(ctx context.Context) ([]User, error)
	// GetByUuid(ctx context.Context, uuid string) (User, error)
	// GetByEmail(ctx context.Context, email string) (User, error)
	GetByUsername(ctx context.Context, uname string) (User, error)
	// Update(ctx context.Context, u *User) error
	Store(ctx context.Context, u *User) error
	Delete(ctx context.Context, uuid string) (string /*uuid*/, error)
	Login(ctx context.Context, uname string, passwd string) (User, error)
	ChangeState(ctx context.Context, uname string, st UserState) error
	ChangeRole(ctx context.Context, uname string, ro Role) error
}
