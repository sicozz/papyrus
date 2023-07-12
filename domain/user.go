package domain

import "context"

// User is representing the User data struct
type User struct {
	Uuid     string    `json:"uuid"`
	UserName string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Name     string    `json:"name"`
	LastName string    `json:"lastname"`
	Role     Role      `json:"role"`
	State    UserState `json:"state"`
}

// UserUsecase represents the user's usecases
type UserUsecase interface {
	Fetch(ctx context.Context) ([]User, error)
	// GetByUuid(ctx context.Context, uuid string) (User, error)
	// GetByEmail(ctx context.Context, email string) (User, error)
	// GetByUserName(ctx context.Context, uname string) (User, error)
	// Update(ctx context.Context, u *User) error
	// Store(ctx context.Context, u *User) error
	// Delete(ctx context.Context, uuid string) error
}

// UserRepository represents the user's repository contract
type UserRepository interface {
	Fetch(ctx context.Context) ([]User, error)
	// GetByUuid(ctx context.Context, uuid string) (User, error)
	// GetByEmail(ctx context.Context, email string) (User, error)
	// GetByUserName(ctx context.Context, uname string) (User, error)
	// Update(ctx context.Context, u *User) error
	// Store(ctx context.Context, u *User) error
	// Delete(ctx context.Context, uuid string) error
}
