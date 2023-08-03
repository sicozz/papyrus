package http

import (
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
	e.GET("/login/:uname/:passwd", handler.Login)
	e.GET("/user", handler.Fetch)
	e.POST("/user", handler.Store)
	e.GET("/user/:uname", handler.GetByUsername)
	e.DELETE("/user/:uname", handler.Delete)
	e.PATCH("/user/:uname/state/:desc", handler.ChangeState)
	e.PATCH("/user/:uname/role/:desc", handler.ChangeRole)
	e.PATCH("/user/:uname", handler.Update)
}

// TODO: Add failure responses when error

func isRequestValid(u *domain.User) (bool, error) {
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *UserHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()
	users, err := u.UUsecase.Fetch(ctx)
	if err != nil {
		domain.AgLog.Error("Could not retrieve users")
		errBody := dtos.NewErrDto("[failure] users fetch")
		return c.JSON(http.StatusInternalServerError, errBody)
	}

	return c.JSON(http.StatusOK, users)
}

func (u *UserHandler) GetByUsername(c echo.Context) error {
	ctx := c.Request().Context()
	uname := c.Param("uname")
	user, err := u.UUsecase.GetByUsername(ctx, uname)
	if err != nil {
		domain.AgLog.Error("Could not retrieve user")
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, user)
}

func (u *UserHandler) Store(c echo.Context) (err error) {
	var user domain.User
	err = c.Bind(&user)
	if err != nil {
		errBody := dtos.NewErrDto("[failure] user creation. could not bind payload")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	if ok, err := isRequestValid(&user); !ok {
		errBody, err := dtos.NewValidationErrDto(err.Error())
		if err != nil {
			return c.JSON(http.StatusBadRequest, errBody)
		}
		return c.JSON(http.StatusBadRequest, errBody)
	}

	ctx := c.Request().Context()
	err = u.UUsecase.Store(ctx, &user)
	if err != nil {
		domain.AgLog.Error("Could not store user: ", err)
	}

	return c.JSON(http.StatusCreated, user)
}

func (u *UserHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	uname := c.Param("uname")
	err := u.UUsecase.Delete(ctx, uname)
	if err != nil {
		domain.AgLog.Error("Could not delete user")
		return c.JSON(http.StatusNotFound, err)
	}

	return c.NoContent(http.StatusOK)
}

func (u *UserHandler) ChangeState(c echo.Context) error {
	ctx := c.Request().Context()
	uname := c.Param("uname")
	desc := c.Param("desc")
	err := u.UUsecase.ChangeState(ctx, uname, desc)
	if err != nil {
		domain.AgLog.Error("Could not patch user")
		return c.JSON(http.StatusNotModified, err)
	}

	return c.NoContent(http.StatusOK)
}

func (u *UserHandler) ChangeRole(c echo.Context) error {
	ctx := c.Request().Context()
	uname := c.Param("uname")
	desc := c.Param("desc")
	err := u.UUsecase.ChangeRole(ctx, uname, desc)
	if err != nil {
		domain.AgLog.Error("Could not patch user")
		return c.JSON(http.StatusNotModified, err)
	}

	return c.NoContent(http.StatusOK)
}

func (u *UserHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	uname := c.Param("uname")
	passwd := c.Param("passwd")
	user, err := u.UUsecase.Login(ctx, uname, passwd)
	if err != nil {
		domain.AgLog.Error("Could not login")
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, user)
}

func (u *UserHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	uname := c.Param("uname")

	var uUpDto dtos.UserUpdateDto
	err := c.Bind(&uUpDto)
	if err != nil {
		errBody := dtos.NewErrDto("[failure] user update. could not bind payload")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	// if ok, err := isRequestValid(&uUpDto); !ok {
	// 	errBody, err := dtos.NewValidationErrDto(err.Error())
	// 	if err != nil {
	// 		return c.JSON(http.StatusBadRequest, errBody)
	// 	}
	// 	return c.JSON(http.StatusBadRequest, errBody)
	// }

	user := domain.User{
		Username: uUpDto.Username,
		Lastname: uUpDto.Lastname,
		Role:     domain.Role{Description: uUpDto.Role},
		State:    domain.UserState{Description: uUpDto.State},
	}
	err = u.UUsecase.Update(ctx, uname, &user)
	if err != nil {
		domain.AgLog.Error("Error in user update: ", err)
	}

	return c.JSON(http.StatusCreated, user)
}
