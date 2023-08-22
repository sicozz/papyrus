package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
	"github.com/sicozz/papyrus/domain/mapper"
	"github.com/sicozz/papyrus/utils"
	"github.com/sicozz/papyrus/utils/constants"
	"gopkg.in/go-playground/validator.v9"
)

// TODO: Sort functions by endpoint order
// DirHandler will initialize the dir/ resources endpoint
type DirHandler struct {
	DUsecase domain.DirUsecase
	log      utils.AggregatedLogger
}

func NewDirHandler(e *echo.Echo, du domain.DirUsecase) {
	logger := utils.NewAggregatedLogger(constants.Delivery, constants.Dir)
	handler := &DirHandler{du, logger}
	e.GET("/dir", handler.GetAll)
	e.POST("/dir", handler.Store)
	e.PATCH("/dir/:uuid", handler.Update)
	e.GET("/dir/:uuid", handler.GetByUuid)
	e.DELETE("/dir/:uuid", handler.Delete)
	e.PATCH("/dir/:uuid/move", handler.Move)

	e.POST("/dir/duplicate", handler.Duplicate)
}

func isRequestValid(p any) (bool, error) {
	validate := validator.New()
	err := validate.Struct(p)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (h *DirHandler) GetAll(c echo.Context) error {
	h.log.Inf("REQ: get all")
	ctx := c.Request().Context()
	dirs, rErr := h.DUsecase.GetAll(ctx)

	if rErr != nil {
		errBody := dtos.NewErrDto("Dir fetch failed")
		return c.JSON(rErr.GetStatus(), errBody)
	}

	dirDtos := make([]dtos.DirGetDto, len(dirs), len(dirs))
	for i, dir := range dirs {
		dirDtos[i] = mapper.MapDirToDirGetDto(dir)
	}

	return c.JSON(http.StatusOK, dirDtos)
}

func (h *DirHandler) GetByUuid(c echo.Context) error {
	h.log.Inf("REQ: get by uuid")
	ctx := c.Request().Context()

	uuid := c.Param("uuid")
	if valid := utils.IsValidUUID(uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	if valid := utils.IsValidUUID(uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	dir, rErr := h.DUsecase.GetByUuid(ctx, uuid)

	if rErr != nil {
		errBody := dtos.NewErrDto("Dir fetch failed")
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, dir)
}

func (h *DirHandler) Store(c echo.Context) (err error) {
	h.log.Inf("REQ: store")
	var dir domain.Dir
	err = c.Bind(&dir)
	if err != nil {
		errBody := dtos.NewErrDto(err.Error())
		return c.JSON(http.StatusBadRequest, errBody)
	}

	if ok, err := isRequestValid(&dir); !ok {
		errBody, err := dtos.NewValidationErrDto(err.Error())
		if err != nil {
			errParse := dtos.NewErrDto(err.Error())
			return c.JSON(http.StatusBadRequest, errParse)
		}
		return c.JSON(http.StatusBadRequest, errBody)
	}

	ctx := c.Request().Context()
	rErr := h.DUsecase.Store(ctx, &dir)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusCreated, dir)
}

func (h *DirHandler) Update(c echo.Context) error {
	h.log.Inf("REQ: update")
	ctx := c.Request().Context()

	uuid := c.Param("uuid")
	if valid := utils.IsValidUUID(uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	// TODO: Change domain entities recvrs for dedicated dtos EVERYWHERE!
	var dUpDto dtos.DirUpdateDto
	err := c.Bind(&dUpDto)
	if err != nil {
		errBody := dtos.NewErrDto(fmt.Sprint("Req body binding failed: ", err))
		return c.JSON(http.StatusBadRequest, errBody)
	}

	if ok, err := isRequestValid(&dUpDto); !ok {
		errBody, err := dtos.NewValidationErrDto(err.Error())
		if err != nil {
			errValid := dtos.NewErrDto(fmt.Sprint("Req body validation failed: ", err))
			return c.JSON(http.StatusBadRequest, errValid)
		}
		return c.JSON(http.StatusBadRequest, errBody)
	}

	dir := domain.Dir{
		Name: dUpDto.Name,
	}

	rErr := h.DUsecase.Update(ctx, uuid, &dir)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}

func (h *DirHandler) Delete(c echo.Context) error {
	h.log.Inf("REQ: delete")
	ctx := c.Request().Context()
	uuid := c.Param("uuid")
	if valid := utils.IsValidUUID(uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	rErr := h.DUsecase.Delete(ctx, uuid)

	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}

func (h *DirHandler) Move(c echo.Context) error {
	// WARN: Add validation to avoid ciclical references
	h.log.Inf("REQ: move")
	ctx := c.Request().Context()

	uuid := c.Param("uuid")
	if valid := utils.IsValidUUID(uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	var dMDto dtos.DirMoveDto
	err := c.Bind(&dMDto)
	if err != nil {
		errBody := dtos.NewErrDto(fmt.Sprint("Req body binding failed: ", err))
		return c.JSON(http.StatusBadRequest, errBody)
	}

	if ok, err := isRequestValid(&dMDto); !ok {
		errBody, err := dtos.NewValidationErrDto(err.Error())
		if err != nil {
			errValid := dtos.NewErrDto(fmt.Sprint("Req body validation failed: ", err))
			return c.JSON(http.StatusBadRequest, errValid)
		}
		return c.JSON(http.StatusBadRequest, errBody)
	}

	rErr := h.DUsecase.Move(ctx, uuid, dMDto.ParentDir)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}

func (h *DirHandler) Duplicate(c echo.Context) error {
	// WARN: Add validation to avoid ciclical references
	h.log.Inf("REQ: duplicate")
	ctx := c.Request().Context()

	var dDDto dtos.DirDuplicateDto
	err := c.Bind(&dDDto)
	if err != nil {
		errBody := dtos.NewErrDto(fmt.Sprint("Req body binding failed: ", err))
		return c.JSON(http.StatusBadRequest, errBody)
	}

	if ok, err := isRequestValid(&dDDto); !ok {
		errBody, err := dtos.NewValidationErrDto(err.Error())
		if err != nil {
			errValid := dtos.NewErrDto(fmt.Sprint("Req body validation failed: ", err))
			return c.JSON(http.StatusBadRequest, errValid)
		}
		return c.JSON(http.StatusBadRequest, errBody)
	}

	res, rErr := h.DUsecase.Duplicate(ctx, dDDto.Uuid, dDDto.Name, dDDto.ParentDir)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, res)
}
