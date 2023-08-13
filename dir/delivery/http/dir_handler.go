package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
	"github.com/sicozz/papyrus/utils"
	"github.com/sicozz/papyrus/utils/constants"
)

// DirHandler will initialize the dir/ resources endpoint
type DirHandler struct {
	DUsecase domain.DirUsecase
	log      utils.AggregatedLogger
}

func NewDirHandler(e *echo.Echo, du domain.DirUsecase) {
	logger := utils.NewAggregatedLogger(constants.Delivery, constants.Dir)
	handler := &DirHandler{du, logger}
	e.GET("/dir", handler.GetAll)
	// e.POST("/dir", handler.Store)
	// e.PATCH("/dir", handler.Update)
	// e.DELETE("/dir/:uuid", handler.Delete)
	// e.GET("/dir/:uuid", handler.GetByUuid)
	// e.PATCH("/dir/:uuid/move", handler.Move)
}

func (h *DirHandler) GetAll(c echo.Context) error {
	h.log.Inf("REQ: get all")
	ctx := c.Request().Context()
	dirs, rErr := h.DUsecase.GetAll(ctx)

	if rErr != nil {
		errBody := dtos.NewErrDto("Dir fetch failed")
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, dirs)
}
