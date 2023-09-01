package utils

import (
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sicozz/papyrus/domain/dtos"
)

func BindFormToPFileUploadDto(c echo.Context, p *dtos.PFileUploadDto) (err error) {
	var val string
	var boolval bool

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
	p.RespUser = val

	if val = c.FormValue("approval_user1"); "" == val {
		return errors.New("File approval user 1 is required. approval_user1:")
	}
	p.AppUser1 = val

	if val = c.FormValue("user_check1"); "" == val {
		return errors.New("File user check 1 is required. user_check1:")
	}
	boolval, err = strconv.ParseBool(val)
	if err != nil {
		return errors.New("File user check 1 must be of type boolean")
	}
	p.Chk1 = boolval

	if val = c.FormValue("approval_user2"); "" != val {
		p.AppUser2 = val

		if val = c.FormValue("user_check2"); "" == val {
			return errors.New("File user check 2 is required. user_check2:")
		}
		boolval, err = strconv.ParseBool(val)
		if err != nil {
			return errors.New("File user check 2 must be of type boolean")
		}
		p.Chk2 = boolval
	}

	if val = c.FormValue("approval_user3"); "" != val {
		p.AppUser3 = val

		if val = c.FormValue("user_check3"); "" == val {
			return errors.New("File user check 3 is required. user_check3:")
		}
		boolval, err = strconv.ParseBool(val)
		if err != nil {
			return errors.New("File user check 3 must be of type boolean")
		}
		p.Chk3 = boolval
	}

	if val = c.FormValue("version"); "" == val {
		return errors.New("File version is required. version:")
	}
	p.Version = val

	if val = c.FormValue("term"); "" == val {
		return errors.New("File term is required. term:")
	}
	term, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return errors.New("File term must be of type integer")
	}
	p.Term = int(term)

	if val = c.FormValue("subtype"); "" == val {
		return errors.New("File version is required. version:")
	}
	p.Subtype = val

	return nil
}
