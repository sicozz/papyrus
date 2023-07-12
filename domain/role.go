package domain

import "context"

// Role is representing the Role data struct
type Role struct {
	Code        int64  `json:"Code"`
	Description string `json:"Description"`
}

// RoleRepository represents the role's repository contract
type RoleRepository interface {
	GetByDescription(ctx context.Context, desc string) (Role, error)
}
