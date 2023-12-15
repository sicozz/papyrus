package mapper

import (
	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
	"github.com/sicozz/papyrus/utils/constants"
)

func MapPlanToPlanGetDto(p domain.Plan) dtos.PlanGetDto {
	return dtos.PlanGetDto{
		Uuid:         p.Uuid,
		Code:         p.Code,
		Name:         p.Name,
		Origin:       p.Origin,
		ActionType:   p.ActionType,
		Term:         p.Term,
		CreatorUser:  p.CreatorUser,
		RespUser:     p.RespUser,
		DateCreation: p.DateCreation.Format(constants.LayoutDate),
		DateClose:    p.DateClose.Format(constants.LayoutDate),
		Causes:       p.Causes,
		Conclusions:  p.Conclusions,
		State:        p.State,
		Stage:        p.Stage,
		Dir:          p.Dir,
		Action0desc:  p.Action0desc,
		Action0date:  p.Action0date,
		Action0user:  p.Action0user,
		Action1desc:  p.Action1desc,
		Action1date:  p.Action1date,
		Action1user:  p.Action1user,
		Action2desc:  p.Action2desc,
		Action2date:  p.Action2date,
		Action2user:  p.Action2user,
		Action3desc:  p.Action3desc,
		Action3date:  p.Action3date,
		Action3user:  p.Action3user,
		Action4desc:  p.Action4desc,
		Action4date:  p.Action4date,
		Action4user:  p.Action4user,
		Action5desc:  p.Action5desc,
		Action5date:  p.Action5date,
		Action5user:  p.Action5user,
	}
}

func MapPlanToDirGetDto(p domain.Plan) dtos.DirGetDto {
	return dtos.DirGetDto{
		Uuid:      p.Uuid,
		Name:      p.Name,
		ParentDir: p.Dir,
		Term:      p.Term,
		RespUser:  p.RespUser,
		Type:      "plan",
	}
}
