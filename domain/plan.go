package domain

import (
	"context"
	"time"

	"github.com/sicozz/papyrus/domain/dtos"
)

type Plan struct {
	Uuid         string
	Code         string
	Name         string
	Origin       string
	ActionType   string
	Term         int
	CreatorUser  string
	RespUser     string
	DateCreation time.Time
	DateClose    string
	Causes       string
	Conclusions  string
	State        string
	Stage        int
	Dir          string
	Tasks        []Task

	Action0desc string
	Action0date string
	Action0user string
	Action1desc string
	Action1date string
	Action1user string
	Action2desc string
	Action2date string
	Action2user string
	Action3desc string
	Action3date string
	Action3user string
	Action4desc string
	Action4date string
	Action4user string
	Action5desc string
	Action5date string
	Action5user string
}

type Action struct {
	Uuid        string
	description string
	CreatorUser string
	DateClose   time.Time
}

type PlanUsecase interface {
	GetAll(c context.Context) ([]dtos.PlanGetDto, RequestErr)
	GetByUuid(c context.Context, uuid string) (dtos.PlanGetDto, RequestErr)
	Store(c context.Context, p dtos.PlanStoreDto) (dtos.PlanGetDto, RequestErr)
	Update(c context.Context, uuid string, p dtos.PlanUpdateDto) (dtos.PlanGetDto, RequestErr)
	Delete(c context.Context, uuid string) RequestErr
}

type PlanRepository interface {
	GetAll(c context.Context) ([]Plan, error)
	GetByUuid(c context.Context, uuid string) (Plan, error)
	Store(c context.Context, p Plan) (string, error)
	Update(c context.Context, uuid string, p Plan) error
	Delete(c context.Context, uuid string) error

	ExistsByUuid(c context.Context, uuid string) bool
}
