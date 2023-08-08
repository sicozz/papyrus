package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
	"gopkg.in/go-playground/validator.v9"
)

// UserHandler will initialize the users/ resources endpoint
type UserHandler struct {
	UUsecase domain.UserUsecase
}

func NewUserHandler(e *echo.Echo, us domain.UserUsecase) {
	handler := &UserHandler{
		UUsecase: us,
	}
	e.GET("/user", handler.Fetch)
	e.POST("/user", handler.Store)
	e.GET("/user/:uname", handler.GetByUsername)
	e.DELETE("/user/:uname", handler.Delete)
	e.PATCH("/user/:uname", handler.Update)
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

func (u *UserHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()
	users, rErr := u.UUsecase.Fetch(ctx)
	if rErr != nil {
		domain.AgLog.Error("[failure] users fetch")
		errBody := dtos.NewErrDto("User fetch failed")
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, users)
}

func (u *UserHandler) GetByUsername(c echo.Context) error {
	ctx := c.Request().Context()
	uname := c.Param("uname")
	user, rErr := u.UUsecase.GetByUsername(ctx, uname)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, user)
}

func (u *UserHandler) Store(c echo.Context) (err error) {
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
	rErr := u.UUsecase.Store(ctx, &user)
	if rErr != nil {
		domain.AgLog.Error("Could not store user: ", err)
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusCreated, user)
}

func (u *UserHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	uname := c.Param("uname")
	rErr := u.UUsecase.Delete(ctx, uname)
	if rErr != nil {
		domain.AgLog.Error("Could not delete user")
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}

func (u *UserHandler) Login(c echo.Context) error {
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

	user, rErr := u.UUsecase.Login(ctx, lDto.Username, lDto.Password)
	if rErr != nil {
		domain.AgLog.Error("Login failed")
		errBody := dtos.NewErrDto("Wrong username or password")
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, user)
}

func (u *UserHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	uname := c.Param("uname")

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
		Name:     uUpDto.Name,
		Lastname: uUpDto.Lastname,
		Email:    uUpDto.Email,
		Role:     domain.Role{Description: uUpDto.Role},
		State:    domain.UserState{Description: uUpDto.State},
	}

	rErr := u.UUsecase.Update(ctx, uname, &user)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}
