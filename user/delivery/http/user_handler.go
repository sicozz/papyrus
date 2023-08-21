package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
	"github.com/sicozz/papyrus/utils"
	"github.com/sicozz/papyrus/utils/constants"
	"gopkg.in/go-playground/validator.v9"
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
	e.DELETE("/user/:uname", handler.Delete)
	e.PATCH("/user/:uuid", handler.Update)
	e.PATCH("/user/:uuid/chg_password", handler.ChgPasswd)
	e.POST("/login", handler.Login)
}

func isRequestValid(u any) (bool, error) {
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		return false, err
	}
	return true, nil
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

func (h *UserHandler) Store(c echo.Context) (err error) {
	h.log.Inf("REQ: store")
	var user domain.User
	err = c.Bind(&user)
	if err != nil {
		errBody := dtos.NewErrDto(err.Error())
		return c.JSON(http.StatusBadRequest, errBody)
	}

	if ok, err := isRequestValid(&user); !ok {
		errBody, err := dtos.NewValidationErrDto(err.Error())
		if err != nil {
			errParse := dtos.NewErrDto(err.Error())
			return c.JSON(http.StatusBadRequest, errParse)
		}
		return c.JSON(http.StatusBadRequest, errBody)
	}

	ctx := c.Request().Context()
	rErr := h.UUsecase.Store(ctx, &user)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusCreated, user)
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

	if ok, err := isRequestValid(&lDto); !ok {
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

	if ok, err := isRequestValid(&uUpDto); !ok {
		errBody, err := dtos.NewValidationErrDto(err.Error())
		if err != nil {
			errValid := dtos.NewErrDto(fmt.Sprint("Req body validation failed: ", err))
			return c.JSON(http.StatusBadRequest, errValid)
		}
		return c.JSON(http.StatusBadRequest, errBody)
	}

	user := domain.User{
		Username: uUpDto.UserName,
		Name:     uUpDto.Name,
		Lastname: uUpDto.Lastname,
		Email:    uUpDto.Email,
		Role:     domain.Role{Description: uUpDto.Role},
		State:    domain.UserState{Description: uUpDto.State},
	}

	rErr := h.UUsecase.Update(ctx, uuid, &user)
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

	var data dtos.ChgPasswdDto
	err := c.Bind(&data)
	if err != nil {
		errBody := dtos.NewErrDto(fmt.Sprint("Req body binding failed: ", err))
		return c.JSON(http.StatusBadRequest, errBody)
	}

	if ok, err := isRequestValid(&data); !ok {
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
