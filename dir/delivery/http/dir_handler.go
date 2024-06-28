package http

import (
	"fmt"
	"net/http"
	"net/smtp"

	"github.com/labstack/echo/v4"
	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
	"github.com/sicozz/papyrus/utils"
	"github.com/sicozz/papyrus/utils/constants"
)

// TODO: Sort functions by endpoint order
// DirHandler will initialize the dir/ resources endpoint
type DirHandler struct {
	DUsecase domain.DirUsecase
	log      utils.AggregatedLogger
}

func NewDirHandler(e *echo.Echo, du domain.DirUsecase) {
	logger := utils.NewAggregatedLogger(constants.Delivery, constants.Dir)
	handler := &DirHandler{du, logger}
	e.GET("/dir", handler.GetAll)
	e.POST("/dir", handler.Store)
	e.GET("/dir/:uuid", handler.GetByUuid)
	e.PATCH("/dir/:uuid", handler.Update)
	e.DELETE("/dir/:uuid", handler.Delete)
	e.PATCH("/dir/:uuid/move", handler.Move)

	e.POST("/dir/duplicate", handler.Duplicate)
	// e.POST("/file/upload", handler.StoreDoc)

	e.GET("/dir/user_docs/:uuid", handler.GetDocsNotDirByUser)
	e.GET("/dir/user_owned_docs/:uuid", handler.GetOwnedDocsNotDirByUser)

	e.GET("/dir/:uuid/size", handler.GetDirSize)

	e.POST("/dir/recursive_permission", handler.AddRecursivePermission)

	e.GET("/email", handler.Email)
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

func (h *DirHandler) GetByUuid(c echo.Context) error {
	h.log.Inf("REQ: get by uuid")
	ctx := c.Request().Context()

	uuid := c.Param("uuid")
	if valid := utils.IsValidUUID(uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	dir, rErr := h.DUsecase.GetByUuid(ctx, uuid)

	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, dir)
}

func (h *DirHandler) Store(c echo.Context) (err error) {
	h.log.Inf("REQ: store")
	var p dtos.DirStoreDto
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
	dir, rErr := h.DUsecase.Store(ctx, p)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusCreated, dir)
}

func (h *DirHandler) Update(c echo.Context) error {
	h.log.Inf("REQ: update")
	ctx := c.Request().Context()

	uuid := c.Param("uuid")
	if valid := utils.IsValidUUID(uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	var p dtos.DirUpdateDto
	err := c.Bind(&p)
	if err != nil {
		errBody := dtos.NewErrDto(fmt.Sprint("Req body binding failed: ", err))
		return c.JSON(http.StatusBadRequest, errBody)
	}

	if ok, err := utils.IsRequestValid(&p); !ok {
		errBody, err := dtos.NewValidationErrDto(err.Error())
		if err != nil {
			errValid := dtos.NewErrDto(fmt.Sprint("Req body validation failed: ", err))
			return c.JSON(http.StatusBadRequest, errValid)
		}
		return c.JSON(http.StatusBadRequest, errBody)
	}

	rErr := h.DUsecase.Update(ctx, uuid, p)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}

func (h *DirHandler) Delete(c echo.Context) error {
	h.log.Inf("REQ: delete")
	ctx := c.Request().Context()
	uuid := c.Param("uuid")
	if valid := utils.IsValidUUID(uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	rErr := h.DUsecase.Delete(ctx, uuid)

	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}

func (h *DirHandler) Move(c echo.Context) error {
	h.log.Inf("REQ: move")
	ctx := c.Request().Context()

	uuid := c.Param("uuid")
	if valid := utils.IsValidUUID(uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	var dMDto dtos.DirMoveDto
	err := c.Bind(&dMDto)
	if err != nil {
		errBody := dtos.NewErrDto(fmt.Sprint("Req body binding failed: ", err))
		return c.JSON(http.StatusBadRequest, errBody)
	}

	if ok, err := utils.IsRequestValid(&dMDto); !ok {
		errBody, err := dtos.NewValidationErrDto(err.Error())
		if err != nil {
			errValid := dtos.NewErrDto(fmt.Sprint("Req body validation failed: ", err))
			return c.JSON(http.StatusBadRequest, errValid)
		}
		return c.JSON(http.StatusBadRequest, errBody)
	}

	rErr := h.DUsecase.Move(ctx, uuid, dMDto.ParentDir)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}

func (h *DirHandler) Duplicate(c echo.Context) error {
	h.log.Inf("REQ: duplicate")
	ctx := c.Request().Context()

	var p dtos.DirDuplicateDto
	err := c.Bind(&p)
	if err != nil {
		errBody := dtos.NewErrDto(fmt.Sprint("Req body binding failed: ", err))
		return c.JSON(http.StatusBadRequest, errBody)
	}

	if ok, err := utils.IsRequestValid(&p); !ok {
		errBody, err := dtos.NewValidationErrDto(err.Error())
		if err != nil {
			errValid := dtos.NewErrDto(fmt.Sprint("Req body validation failed: ", err))
			return c.JSON(http.StatusBadRequest, errValid)
		}
		return c.JSON(http.StatusBadRequest, errBody)
	}

	res, rErr := h.DUsecase.Duplicate(ctx, p)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, res)
}

func (h *DirHandler) GetDocsNotDirByUser(c echo.Context) error {
	h.log.Inf("REQ: get docs not dir by user")
	ctx := c.Request().Context()

	uuid := c.Param("uuid")
	if valid := utils.IsValidUUID(uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	dir, rErr := h.DUsecase.GetDocsByUser(ctx, uuid)

	if rErr != nil {
		return c.JSON(rErr.GetStatus(), rErr.Error())
	}

	return c.JSON(http.StatusOK, dir)
}

func (h *DirHandler) GetOwnedDocsNotDirByUser(c echo.Context) error {
	h.log.Inf("REQ: get owned docs not dir by user")
	ctx := c.Request().Context()

	uuid := c.Param("uuid")
	if valid := utils.IsValidUUID(uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	dir, rErr := h.DUsecase.GetOwnedDocsByUser(ctx, uuid)

	if rErr != nil {
		return c.JSON(rErr.GetStatus(), rErr.Error())
	}

	return c.JSON(http.StatusOK, dir)
}

func (h *DirHandler) GetDirSize(c echo.Context) error {
	h.log.Inf("REQ: get by uuid")
	ctx := c.Request().Context()

	uuid := c.Param("uuid")
	if valid := utils.IsValidUUID(uuid); !valid {
		errBody := dtos.NewErrDto("Uuid does not conform to the uuid format")
		return c.JSON(http.StatusBadRequest, errBody)
	}

	res, rErr := h.DUsecase.GetDirSize(ctx, uuid)
	if rErr != nil {
		return c.JSON(rErr.GetStatus(), rErr.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func (h *DirHandler) AddRecursivePermission(c echo.Context) error {
	h.log.Inf("REQ: AddRecursivePermission")
	ctx := c.Request().Context()

	var data dtos.UserAddPermissionDto
	err := c.Bind(&data)
	if err != nil {
		errBody := dtos.NewErrDto(fmt.Sprint("Req body binding failed: ", err))
		return c.JSON(http.StatusBadRequest, errBody)
	}

	if ok, err := utils.IsRequestValid(&data); !ok {
		errBody, err := dtos.NewValidationErrDto(err.Error())
		if err != nil {
			errValid := dtos.NewErrDto(fmt.Sprint("Req body validation failed: ", err))
			return c.JSON(http.StatusBadRequest, errValid)
		}
		return c.JSON(http.StatusBadRequest, errBody)
	}

	rErr := h.DUsecase.AddRecursivePermission(ctx, data)
	if rErr != nil {
		errBody := dtos.NewErrDto(rErr.Error())
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.NoContent(http.StatusOK)
}

func (h *DirHandler) Email(c echo.Context) error {
	h.log.Inf("Send email test")

	Server := "smtp.masterplac.com"
	Sender := "info@masterplac.com"
	Passwd := "Mail.2012+"
	Subject := "Subject: Mensaje de PAPYRUS\r\n\r\n"
	Port := 465
	to := "simozuluaga@gmail.com"

	h.log.Inf("Trying to authenticate...")
	auth := smtp.PlainAuth("", Sender, Passwd, Server)
	h.log.Inf(auth)

	h.log.Inf("Successful authentication...")
	sliceTo := []string{to}
	byteMsg := []byte(fmt.Sprintf("%v%v", Subject, "Email sending test email"))

	h.log.Inf("Sending email...")
	err := smtp.SendMail(
		fmt.Sprintf("%v:%v", Server, Port),
		auth,
		Sender,
		sliceTo,
		byteMsg,
	)

	if err != nil {
		h.log.Err("ERROR: ", err)
	} else {
		h.log.Inf("Email sent...")
	}

	return err
}
