package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
	"github.com/sicozz/papyrus/utils"
	"github.com/sicozz/papyrus/utils/constants"
)

type PFileHandler struct {
	PFUsecase domain.PFileUsecase
	log       utils.AggregatedLogger
}

func NewPFileHandler(e *echo.Echo, uu domain.PFileUsecase) {
	logger := utils.NewAggregatedLogger(constants.Delivery, constants.PFile)
	handler := &PFileHandler{uu, logger}
	e.GET("/file", handler.GetAll)
	e.POST("/file", handler.Upload)
	e.DELETE("/file/:uuid", handler.Delete)
}

func (h *PFileHandler) GetAll(c echo.Context) error {
	h.log.Inf("REQ: get all")
	ctx := c.Request().Context()
	users, rErr := h.PFUsecase.GetAll(ctx)
	if rErr != nil {
		errBody := dtos.NewErrDto("User fetch failed")
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, users)
}

func (h *PFileHandler) Upload(c echo.Context) (err error) {
	h.log.Inf("REQ: upload")
	var p dtos.PFileUploadDto
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
	pf, rErr := h.PFUsecase.Upload(ctx, p)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusCreated, pf)
}

func (h *PFileHandler) Delete(c echo.Context) error {
	h.log.Inf("REQ: delete")
	ctx := c.Request().Context()
	uuid := c.Param("uuid")
	if valid := utils.IsValidUUID(uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	rErr := h.PFUsecase.Delete(ctx, uuid)

	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}
