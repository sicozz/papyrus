package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sicozz/papyrus/domain"
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
	// e.GET("/users/:id", handler.GetByID)
	e.DELETE("/user/:uname", handler.Delete)
	e.PATCH("/user/:uname/state/:desc", handler.ChangeState)
	e.PATCH("/user/:uname/role/:desc", handler.ChangeRole)
}

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
	}

	return c.JSON(http.StatusOK, users)
}

func (u *UserHandler) Store(c echo.Context) (err error) {
	var user domain.User
	err = c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&user); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
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
