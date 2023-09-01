package http

import (
	"net/http"
	"os"

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
	e.GET("/file/:uuid", handler.Download)
	e.DELETE("/file/:uuid", handler.Delete)

	e.PATCH("/file/:file_uuid/user/:user_uuid", handler.Approve)
	e.PATCH("/file/:uuid/activate", handler.Activate)
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
	file, err := c.FormFile("file")
	if err != nil {
		errBody := dtos.NewErrDto("Upload file is required")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	var p dtos.PFileUploadDto
	err = utils.BindFormToPFileUploadDto(c, &p)
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
	pf, rErr := h.PFUsecase.Upload(ctx, p, file)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusCreated, pf)
}

func (h *PFileHandler) Download(c echo.Context) error {
	h.log.Inf("REQ: download")
	ctx := c.Request().Context()
	uuid := c.Param("uuid")
	if valid := utils.IsValidUUID(uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	pFile, rErr := h.PFUsecase.GetByUuid(ctx, uuid)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	file, err := os.Open(pFile.FsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return c.Stream(http.StatusOK, constants.RespTypeStream, file)
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

func (h *PFileHandler) Approve(c echo.Context) error {
	h.log.Inf("REQ: approve")
	ctx := c.Request().Context()
	pfUuid := c.Param("file_uuid")
	userUuid := c.Param("user_uuid")
	if valid := utils.IsValidUUID(pfUuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}
	if valid := utils.IsValidUUID(userUuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	rErr := h.PFUsecase.Approve(ctx, pfUuid, userUuid)

	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}

func (h *PFileHandler) Activate(c echo.Context) error {
	h.log.Inf("REQ: activate")
	ctx := c.Request().Context()
	uuid := c.Param("uuid")
	if valid := utils.IsValidUUID(uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	rErr := h.PFUsecase.Activate(ctx, uuid)

	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}
