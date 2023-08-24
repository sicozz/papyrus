package domain

import "context"

// Role is representing the Role data struct
type Role struct {
	Code        int64  `json:"code"`
	Description string `json:"description"`
}

// RoleRepository represents the role's repository contract
type RoleRepository interface {
	GetByCode(ctx context.Context, code int64) (Role, error)
	GetAll(ctx context.Context) ([]Role, error)
	GetByDescription(ctx context.Context, desc string) (Role, error)
	ExistsByDescription(ctx context.Context, desc string) bool
}
