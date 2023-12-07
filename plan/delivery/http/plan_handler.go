package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
	"github.com/sicozz/papyrus/utils"
	"github.com/sicozz/papyrus/utils/constants"
)

type PlanHandler struct {
	planUsecase domain.PlanUsecase
	log         utils.AggregatedLogger
}

func NewPlanHandler(e *echo.Echo, du domain.PlanUsecase) {
	logger := utils.NewAggregatedLogger(constants.Delivery, constants.Plan)
	handler := &PlanHandler{du, logger}
	e.GET("/plan", handler.GetAll)
	e.GET("/plan/:uuid", handler.GetByUuid)
	e.POST("/plan", handler.Store)
	e.PUT("/plan/:uuid", handler.Update)
	e.DELETE("/plan/:uuid", handler.Delete)
}

func (h *PlanHandler) GetAll(c echo.Context) error {
	h.log.Inf("REQ: get all")
	ctx := c.Request().Context()
	plans, rErr := h.planUsecase.GetAll(ctx)

	if rErr != nil {
		errBody := dtos.NewErrDto("Plan fetch failed")
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, plans)
}

func (h *PlanHandler) GetByUuid(c echo.Context) error {
	h.log.Inf("REQ: get by uuid")
	ctx := c.Request().Context()

	uuid := c.Param("uuid")
	if valid := utils.IsValidUUID(uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	plan, rErr := h.planUsecase.GetByUuid(ctx, uuid)

	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, plan)
}

func (h *PlanHandler) Store(c echo.Context) (err error) {
	h.log.Inf("REQ: store")
	var p dtos.PlanStoreDto
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
	plan, rErr := h.planUsecase.Store(ctx, p)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusCreated, plan)
}

func (h *PlanHandler) Update(c echo.Context) (err error) {
	h.log.Inf("REQ: update")
	ctx := c.Request().Context()
	uuid := c.Param("uuid")
	if valid := utils.IsValidUUID(uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	var p dtos.PlanUpdateDto
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

	plan, rErr := h.planUsecase.Update(ctx, uuid, p)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusCreated, plan)
}

func (h *PlanHandler) Delete(c echo.Context) error {
	h.log.Inf("REQ: delete")
	ctx := c.Request().Context()
	uuid := c.Param("uuid")
	if valid := utils.IsValidUUID(uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	rErr := h.planUsecase.Delete(ctx, uuid)

	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}
