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
	e.GET("/file/:uuid", handler.GetByUuid)
	e.DELETE("/file/:uuid", handler.Delete)

	e.POST("/file/:file_uuid/download", handler.Download)
	e.PATCH("/file/:file_uuid/user/:user_uuid/check", handler.ChgApprovation)
	e.PATCH("/file/:file_uuid/user/:user_uuid/state", handler.ChgState)

	e.GET("/evidence/task/:uuid", handler.GetAllEvidence)
	e.POST("/evidence/task/:uuid", handler.UploadEvidence)
	e.DELETE("/evidence/task/:task_uuid/file/:file_uuid", handler.DeleteEvidence)
}

func (h *PFileHandler) GetAll(c echo.Context) error {
	h.log.Inf("REQ: get all")
	ctx := c.Request().Context()
	pfiles, rErr := h.PFUsecase.GetAll(ctx)
	if rErr != nil {
		errBody := dtos.NewErrDto("File fetch failed")
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, pfiles)
}

func (h *PFileHandler) GetByUuid(c echo.Context) error {
	h.log.Inf("REQ: get by uuid")
	ctx := c.Request().Context()

	uuid := c.Param("uuid")
	if valid := utils.IsValidUUID(uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	pfDto, rErr := h.PFUsecase.GetByUuid(ctx, uuid)
	if rErr != nil {
		errBody := dtos.NewErrDto("File fetch failed")
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, pfDto)
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

	var p dtos.PFileDownloadDto
	err := c.Bind(&p)
	if err != nil {
		errBody := dtos.NewErrDto(err.Error())
		return c.JSON(http.StatusBadRequest, errBody)
	}

	pfUuid := c.Param("file_uuid")
	if valid := utils.IsValidUUID(pfUuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	pfPath, rErr := h.PFUsecase.RequestDownload(ctx, pfUuid, p.UserUuid)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	file, err := os.Open(pfPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if !p.Registered {
		return c.Stream(http.StatusOK, constants.RespTypeStream, file)
	}

	rErr = h.PFUsecase.AddDwnHistory(ctx, pfUuid, p.UserUuid)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

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

func (h *PFileHandler) ChgApprovation(c echo.Context) error {
	h.log.Inf("REQ: chg approvation")
	ctx := c.Request().Context()

	var p dtos.PFileChgCheckDto
	err := c.Bind(&p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// if ok, err := utils.IsRequestValid(&p); !ok {
	// 	errBody, err := dtos.NewValidationErrDto(err.Error())
	// 	if err != nil {
	// 		errParse := dtos.NewErrDto(err.Error())
	// 		return c.JSON(http.StatusBadRequest, errParse)
	// 	}
	// 	return c.JSON(http.StatusBadRequest, errBody)
	// }

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

	rErr := h.PFUsecase.ChgApprovation(ctx, pfUuid, userUuid, p.Chk)

	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}

func (h *PFileHandler) ChgState(c echo.Context) error {
	h.log.Inf("REQ: activate")
	ctx := c.Request().Context()

	var p dtos.PFileChgStateDto
	err := c.Bind(&p)
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

	pfile_uuid := c.Param("file_uuid")
	user_uuid := c.Param("user_uuid")
	if valid := utils.IsValidUUID(pfile_uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}
	if valid := utils.IsValidUUID(user_uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	rErr := h.PFUsecase.ChgState(ctx, pfile_uuid, user_uuid, p.StateDesc)

	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}

func (h *PFileHandler) UploadEvidence(c echo.Context) (err error) {
	h.log.Inf("REQ: upload evidence")
	tUuid := c.Param("uuid")
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
	pf, rErr := h.PFUsecase.UploadEvidence(ctx, tUuid, p, file)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusCreated, pf)
}

func (h *PFileHandler) DeleteEvidence(c echo.Context) error {
	h.log.Inf("REQ: delete evidence")
	ctx := c.Request().Context()
	tUuid := c.Param("task_uuid")
	pfUuid := c.Param("file_uuid")
	if valid := utils.IsValidUUID(tUuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}
	if valid := utils.IsValidUUID(pfUuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	rErr := h.PFUsecase.DeleteEvidence(ctx, tUuid, pfUuid)

	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}

func (h *PFileHandler) GetAllEvidence(c echo.Context) error {
	h.log.Inf("REQ: get all")
	ctx := c.Request().Context()
	tUuid := c.Param("uuid")
	if valid := utils.IsValidUUID(tUuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	evidences, rErr := h.PFUsecase.GetEvidence(ctx, tUuid)
	if rErr != nil {
		errBody := dtos.NewErrDto("File fetch failed")
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, evidences)
}
