package domain

import "context"

// UserState is representing the userState data struct
type UserState struct {
	Code        int64  `json:"Code"`
	Description string `json:"Description"`
}

// UserStateRepository represents the userStates's repository contract
type UserStateRepository interface {
	GetByCode(ctx context.Context, code int64) (UserState, error)
	GetByDescription(ctx context.Context, desc string) (UserState, error)
}
