package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
	"github.com/sicozz/papyrus/utils"
	"github.com/sicozz/papyrus/utils/constants"
)

// TaskHandler will initialize the task endpoint
type TaskHandler struct {
	TUsecase domain.TaskUsecase
	log      utils.AggregatedLogger
}

func NewTaskHandler(e *echo.Echo, tu domain.TaskUsecase) {
	logger := utils.NewAggregatedLogger(constants.Delivery, constants.Task)
	handler := &TaskHandler{tu, logger}
	e.GET("/task", handler.GetAll)
	e.POST("/task", handler.Store)
	e.POST("/task/multiple", handler.StoreMultiple)
	e.GET("/task/:uuid", handler.GetByUuid)
	e.GET("/task/user/:uuid", handler.GetByUser)
	e.GET("/task/plan/:uuid", handler.GetByPlan)
	e.DELETE("/task/:task_uuid/user/:user_uuid", handler.Delete)
	e.PATCH("/task/:task_uuid/user/:user_uuid/check", handler.ChgCheck)
	e.PATCH("/task/:task_uuid/user/:user_uuid/state", handler.ChgState)
}

func (h *TaskHandler) GetAll(c echo.Context) error {
	h.log.Inf("REQ: get all")
	ctx := c.Request().Context()
	dirs, rErr := h.TUsecase.GetAll(ctx)

	if rErr != nil {
		errBody := dtos.NewErrDto("Task fetch failed")
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, dirs)
}

func (h *TaskHandler) GetByUuid(c echo.Context) error {
	h.log.Inf("REQ: get by uuid")
	ctx := c.Request().Context()

	uuid := c.Param("uuid")
	if valid := utils.IsValidUUID(uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	dir, rErr := h.TUsecase.GetByUuid(ctx, uuid)

	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, dir)
}

func (h *TaskHandler) Store(c echo.Context) (err error) {
	h.log.Inf("REQ: store")
	var p dtos.TaskStoreDto
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
	dir, rErr := h.TUsecase.Store(ctx, p)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusCreated, dir)
}

func (h *TaskHandler) StoreMultiple(c echo.Context) (err error) {
	h.log.Inf("REQ: store")
	var inputDtos []dtos.TaskStoreDto
	err = c.Bind(&inputDtos)
	if err != nil {
		errBody := dtos.NewErrDto(err.Error())
		return c.JSON(http.StatusBadRequest, errBody)
	}

	for _, d := range inputDtos {
		if ok, err := utils.IsRequestValid(&d); !ok {
			errBody, err := dtos.NewValidationErrDto(err.Error())
			if err != nil {
				errParse := dtos.NewErrDto(err.Error())
				return c.JSON(http.StatusBadRequest, errParse)
			}
			return c.JSON(http.StatusBadRequest, errBody)
		}
	}

	ctx := c.Request().Context()
	newTasks, rErr := h.TUsecase.StoreMultiple(ctx, inputDtos)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusCreated, newTasks)
}

func (h *TaskHandler) GetByUser(c echo.Context) error {
	h.log.Inf("REQ: get by user")
	ctx := c.Request().Context()

	uuid := c.Param("uuid")
	if valid := utils.IsValidUUID(uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	dir, rErr := h.TUsecase.GetByUser(ctx, uuid)

	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, dir)
}

func (h *TaskHandler) GetByPlan(c echo.Context) error {
	h.log.Inf("REQ: get by plan")
	ctx := c.Request().Context()

	uuid := c.Param("uuid")
	if valid := utils.IsValidUUID(uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	tasks, rErr := h.TUsecase.GetByPlan(ctx, uuid)

	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) ChgCheck(c echo.Context) (err error) {
	h.log.Inf("REQ: change check")
	ctx := c.Request().Context()

	var p dtos.TaskChgCheck
	err = c.Bind(&p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if ok, err := utils.IsRequestValid(&p); !ok {
		errBody, err := dtos.NewValidationErrDto(err.Error())
		if err != nil {
			errParse := dtos.NewErrDto(err.Error())
			return c.JSON(http.StatusBadRequest, errParse)
		}
		return c.JSON(http.StatusBadRequest, errBody)
	}

	tUuid := c.Param("task_uuid")
	userUuid := c.Param("user_uuid")
	if valid := utils.IsValidUUID(tUuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}
	if valid := utils.IsValidUUID(userUuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	rErr := h.TUsecase.ChgCheck(ctx, tUuid, userUuid, p.Chk)

	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}

func (h *TaskHandler) ChgState(c echo.Context) (err error) {
	h.log.Inf("REQ: change state")
	ctx := c.Request().Context()

	var p dtos.TaskChgStateDto
	err = c.Bind(&p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if ok, err := utils.IsRequestValid(&p); !ok {
		errBody, err := dtos.NewValidationErrDto(err.Error())
		if err != nil {
			errParse := dtos.NewErrDto(err.Error())
			return c.JSON(http.StatusBadRequest, errParse)
		}
		return c.JSON(http.StatusBadRequest, errBody)
	}

	tUuid := c.Param("task_uuid")
	userUuid := c.Param("user_uuid")
	if valid := utils.IsValidUUID(tUuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}
	if valid := utils.IsValidUUID(userUuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	rErr := h.TUsecase.ChgState(ctx, tUuid, userUuid, p.StateDesc)

	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}

func (h *TaskHandler) Delete(c echo.Context) error {
	h.log.Inf("REQ: delete")
	ctx := c.Request().Context()
	tUuid := c.Param("task_uuid")
	uUuid := c.Param("user_uuid")
	if valid := utils.IsValidUUID(tUuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}
	if valid := utils.IsValidUUID(uUuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	rErr := h.TUsecase.Delete(ctx, tUuid, uUuid)

	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}
