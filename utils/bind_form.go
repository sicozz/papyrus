package utils

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/sicozz/papyrus/domain/dtos"
)

func BindFormToPFileUploadDto(c echo.Context, p *dtos.PFileUploadDto) (err error) {
	var val string
	if val = c.FormValue("code"); "" == val {
		return errors.New("File code is required. code:")
	}
	p.Code = val

	if val = c.FormValue("name"); "" == val {
		return errors.New("File name is required. name:")
	}
	p.Name = val

	if val = c.FormValue("date_create"); "" == val {
		return errors.New("File creation date is required. date_create:")
	}
	p.DateCreation = val

	if val = c.FormValue("type"); "" == val {
		return errors.New("File type is required. type:")
	}
	p.Type = val

	if val = c.FormValue("dir"); "" == val {
		return errors.New("File destination directory is required. dir:")
	}
	p.Dir = val

	if val = c.FormValue("responsible_user"); "" == val {
		return errors.New("File responsible user is required. responsible_user:")
	}
	p.RevUser = val

	if val = c.FormValue("approval_user"); "" == val {
		return errors.New("File approval user is required. approval_user:")
	}
	p.AppUser = c.FormValue("approval_user")

	return nil
}
