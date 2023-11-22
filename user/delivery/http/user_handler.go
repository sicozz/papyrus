package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
	"github.com/sicozz/papyrus/utils"
	"github.com/sicozz/papyrus/utils/constants"
)

// UserHandler will initialize the user/ resources endpoint
type UserHandler struct {
	UUsecase domain.UserUsecase
	log      utils.AggregatedLogger
}

func NewUserHandler(e *echo.Echo, uu domain.UserUsecase) {
	logger := utils.NewAggregatedLogger(constants.Delivery, constants.User)
	handler := &UserHandler{uu, logger}
	e.GET("/user", handler.GetAll)
	e.POST("/user", handler.Store)
	e.GET("/user/:uname", handler.GetByUsername)
	e.GET("/user/:uuid/uuid", handler.GetByUuid)
	e.DELETE("/user/:uname", handler.Delete)
	e.PATCH("/user/:uuid", handler.Update)
	e.PATCH("/user/:uuid/chg_password", handler.ChgPasswd)
	e.PATCH("/user/:uuid/rst_password", handler.RstPasswd)
	e.POST("/login", handler.Login)

	e.GET("/user/:uuid/permission", handler.GetUserPermittedDirs)
	e.POST("/permission", handler.AddPermission)
	e.DELETE("/permission/user/:user_uuid/dir/:dir_uuid", handler.RevokePermission)

	e.GET("/user/:uuid/history/downloads", handler.GetHistoryDownloads)
}

func (h *UserHandler) GetAll(c echo.Context) error {
	h.log.Inf("REQ: get all")
	ctx := c.Request().Context()
	users, rErr := h.UUsecase.GetAll(ctx)
	if rErr != nil {
		errBody := dtos.NewErrDto("User fetch failed")
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetByUsername(c echo.Context) error {
	h.log.Inf("REQ: get by username")
	ctx := c.Request().Context()
	uname := c.Param("uname")
	user, rErr := h.UUsecase.GetByUsername(ctx, uname)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetByUuid(c echo.Context) error {
	h.log.Inf("REQ: get by uuid")
	ctx := c.Request().Context()
	uuid := c.Param("uuid")

	user, rErr := h.UUsecase.GetByUuid(ctx, uuid)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Store(c echo.Context) (err error) {
	h.log.Inf("REQ: store")
	var p dtos.UserStore
	err = c.Bind(&p)
	if err != nil {
		errBody := dtos.NewErrDto(err.Error())
		return c.JSON(http.StatusBadRequest, errBody)
	}

	if ok, err := utils.IsRequestValid(&p); !ok {
		errBody, err := dtos.NewValidationErrDto(err.Error())
		if err != nil {
			errParse := dtos.NewErrDto(err.Error())
			return c.JSON(http.StatusBadRequest, errParse)
		}
		return c.JSON(http.StatusBadRequest, errBody)
	}

	ctx := c.Request().Context()
	nUser, rErr := h.UUsecase.Store(ctx, p)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusCreated, nUser)
}

func (h *UserHandler) Delete(c echo.Context) error {
	h.log.Inf("REQ: delete")
	ctx := c.Request().Context()
	uname := c.Param("uname")
	rErr := h.UUsecase.Delete(ctx, uname)

	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}

func (h *UserHandler) Login(c echo.Context) error {
	h.log.Inf("REQ: login")
	ctx := c.Request().Context()
	var lDto dtos.LoginDto
	err := c.Bind(&lDto)
	if err != nil {
		errBody := dtos.NewErrDto(fmt.Sprint("Req body binding failed: ", err))
		return c.JSON(http.StatusBadRequest, errBody)
	}

	if ok, err := utils.IsRequestValid(&lDto); !ok {
		errBody, err := dtos.NewValidationErrDto(err.Error())
		if err != nil {
			errValid := dtos.NewErrDto(fmt.Sprint("Req body validation failed: ", err))
			return c.JSON(http.StatusBadRequest, errValid)
		}
		return c.JSON(http.StatusBadRequest, errBody)
	}

	user, rErr := h.UUsecase.Login(ctx, lDto.Username, lDto.Password)
	if rErr != nil {
		errBody := dtos.NewErrDto("Wrong username or password")
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Update(c echo.Context) error {
	h.log.Inf("REQ: update")
	ctx := c.Request().Context()

	uuid := c.Param("uuid")
	if valid := utils.IsValidUUID(uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	var uUpDto dtos.UserUpdateDto
	err := c.Bind(&uUpDto)
	if err != nil {
		errBody := dtos.NewErrDto(fmt.Sprint("Req body binding failed: ", err))
		return c.JSON(http.StatusBadRequest, errBody)
	}

	if ok, err := utils.IsRequestValid(&uUpDto); !ok {
		errBody, err := dtos.NewValidationErrDto(err.Error())
		if err != nil {
			errValid := dtos.NewErrDto(fmt.Sprint("Req body validation failed: ", err))
			return c.JSON(http.StatusBadRequest, errValid)
		}
		return c.JSON(http.StatusBadRequest, errBody)
	}

	rErr := h.UUsecase.Update(ctx, uuid, uUpDto)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}

func (h *UserHandler) ChgPasswd(c echo.Context) error {
	h.log.Inf("REQ: change password")
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

	var data dtos.UserChgPasswdDto
	err := c.Bind(&data)
	if err != nil {
		errBody := dtos.NewErrDto(fmt.Sprint("Req body binding failed: ", err))
		return c.JSON(http.StatusBadRequest, errBody)
	}

	if ok, err := utils.IsRequestValid(&data); !ok {
		errBody, err := dtos.NewValidationErrDto(err.Error())
		if err != nil {
			errValid := dtos.NewErrDto(fmt.Sprint("Req body validation failed: ", err))
			return c.JSON(http.StatusBadRequest, errValid)
		}
		return c.JSON(http.StatusBadRequest, errBody)
	}

	rErr := h.UUsecase.ChgPasswd(ctx, uuid, data)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}

func (h *UserHandler) RstPasswd(c echo.Context) error {
	h.log.Inf("REQ: reset password")
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

	var data dtos.UserChgPasswdDto
	err := c.Bind(&data)
	if err != nil {
		errBody := dtos.NewErrDto(fmt.Sprint("Req body binding failed: ", err))
		return c.JSON(http.StatusBadRequest, errBody)
	}

	if ok, err := utils.IsRequestValid(&data); !ok {
		errBody, err := dtos.NewValidationErrDto(err.Error())
		if err != nil {
			errValid := dtos.NewErrDto(fmt.Sprint("Req body validation failed: ", err))
			return c.JSON(http.StatusBadRequest, errValid)
		}
		return c.JSON(http.StatusBadRequest, errBody)
	}

	rErr := h.UUsecase.RstPasswd(ctx, uuid, data)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}

func (h *UserHandler) AddPermission(c echo.Context) error {
	h.log.Inf("REQ: change password")
	ctx := c.Request().Context()

	var data dtos.UserAddPermissionDto
	err := c.Bind(&data)
	if err != nil {
		errBody := dtos.NewErrDto(fmt.Sprint("Req body binding failed: ", err))
		return c.JSON(http.StatusBadRequest, errBody)
	}

	if ok, err := utils.IsRequestValid(&data); !ok {
		errBody, err := dtos.NewValidationErrDto(err.Error())
		if err != nil {
			errValid := dtos.NewErrDto(fmt.Sprint("Req body validation failed: ", err))
			return c.JSON(http.StatusBadRequest, errValid)
		}
		return c.JSON(http.StatusBadRequest, errBody)
	}

	rErr := h.UUsecase.AddPermission(ctx, data.UserUuid, data.DirUuid)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}

func (h *UserHandler) RevokePermission(c echo.Context) error {
	h.log.Inf("REQ: chg approvation")
	ctx := c.Request().Context()

	userUuid := c.Param("user_uuid")
	dirUuid := c.Param("dir_uuid")
	if valid := utils.IsValidUUID(userUuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}
	if valid := utils.IsValidUUID(dirUuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	rErr := h.UUsecase.RevokePermission(ctx, userUuid, dirUuid)

	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}

func (h *UserHandler) GetUserPermittedDirs(c echo.Context) error {
	h.log.Inf("REQ: get all")
	ctx := c.Request().Context()

	userUuid := c.Param("uuid")
	if valid := utils.IsValidUUID(userUuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	dirs, rErr := h.UUsecase.GetUserPermittedDirs(ctx, userUuid)
	if rErr != nil {
		errBody := dtos.NewErrDto("Dir fetch failed")
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, dirs)
}

func (h *UserHandler) GetHistoryDownloads(c echo.Context) error {
	h.log.Inf("REQ: get history downloads")
	ctx := c.Request().Context()

	userUuid := c.Param("uuid")
	if valid := utils.IsValidUUID(userUuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	hist, rErr := h.UUsecase.GetHistoryDownloads(ctx, userUuid)
	if rErr != nil {
		errBody := dtos.NewErrDto("History downloads fetch failed")
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, hist)
}
