package domain

import (
	"context"
	"time"

	"github.com/sicozz/papyrus/domain/dtos"
)

type Task struct {
	Uuid         string
	Name         string
	Procedure    string
	DateCreation time.Time
	DateCheck    string
	DateClose    string
	Term         int
	State        string
	Dir          string
	CreatorUser  string
	RecvUser     string
	Check        bool
	Plan         string
}

type TaskUsecase interface {
	GetAll(c context.Context) ([]dtos.TaskGetDto, RequestErr)
	GetByUuid(c context.Context, uuid string) (dtos.TaskGetDto, RequestErr)
	Store(c context.Context, t dtos.TaskStoreDto) (dtos.TaskGetDto, RequestErr)
	StoreMultiple(c context.Context, inputDtos []dtos.TaskStoreDto) ([]dtos.TaskGetDto, RequestErr)
	ChgCheck(c context.Context, tUuid, uUuid string, chk bool) RequestErr
	ChgState(c context.Context, tUuid, uUuid string, desc string) RequestErr
	Delete(c context.Context, tUuid, uUuid string) RequestErr

	GetByUser(c context.Context, uuid string) ([]dtos.TaskGetDto, RequestErr)
	GetByPlan(c context.Context, uuid string) ([]dtos.TaskGetDto, RequestErr)
}

type TaskRepository interface {
	GetAll(ctx context.Context) ([]Task, error)
	GetByUuid(ctx context.Context, uuid string) (Task, error)
	ExistsByUuid(ctx context.Context, uuid string) bool
	Store(ctx context.Context, t Task) (string, error)
	StoreMultiple(ctx context.Context, tasks []Task) error
	ChgCheck(ctx context.Context, tUuid string, uUuid string, chk bool) error
	ChgState(ctx context.Context, tUuid string, uUuid string, desc string) error
	Delete(ctx context.Context, tUuid, uUuid string) error

	ExistsStateByDesc(ctx context.Context, desc string) bool
	SetDateCheck(ctx context.Context, tUuid string) error

	GetByUser(c context.Context, uuid string) ([]Task, error)
	GetOwnedByUser(c context.Context, uuid string) ([]Task, error)
	GetByCreatorOrRecv(c context.Context, uuid string) ([]Task, error)
	GetByPlan(c context.Context, uuid string) ([]Task, error)
}
