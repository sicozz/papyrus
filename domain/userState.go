package domain

import "context"

// UserState is representing the userState data struct
type UserState struct {
	Code        int64  `json:"code"`
	Description string `json:"description"`
}

// UserStateRepository represents the userStates's repository contract
type UserStateRepository interface {
	GetByCode(ctx context.Context, code int64) (UserState, error)
	GetAll(ctx context.Context) ([]UserState, error)
	GetByDescription(ctx context.Context, desc string) (UserState, error)
}
