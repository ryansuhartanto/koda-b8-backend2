package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/model"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/service"
)

const maxPictureSize = 2 << 20 // 2 MiB

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) HandleRegister(ctx *gin.Context) {
	var new model.User
	if err := ctx.ShouldBind(&new); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Register(new)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrEmailConflict) {
			status = http.StatusConflict
		}
		ctx.AbortWithStatusJSON(status, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (h *UserHandler) HandleLogin(ctx *gin.Context) {
	var cre model.Credentials
	if err := ctx.ShouldBind(&cre); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Login(cre)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrEmailUnregistered) {
			status = http.StatusUnauthorized
		}
		if errors.Is(err, service.ErrPasswordIncorrect) {
			status = http.StatusUnprocessableEntity
		}
		ctx.AbortWithStatusJSON(status, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *UserHandler) HandleList(ctx *gin.Context) {
	users, err := h.service.List()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (h *UserHandler) HandlePatch(ctx *gin.Context) {
	var id model.Id
	if err := ctx.ShouldBindUri(&id); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	var new model.User
	if err := ctx.ShouldBind(&new); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Edit(id, new)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *UserHandler) HandlePutPicture(ctx *gin.Context) {
	var id model.Id
	if err := ctx.ShouldBindUri(&id); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	var data []byte
	if fileHeader, err := ctx.FormFile("picture"); err == nil {
		if fileHeader.Size > maxPictureSize {
			ctx.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, "picture: file too large")
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}
		defer file.Close()

		data, err = io.ReadAll(file)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}
	} else {
		body, err := io.ReadAll(http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxPictureSize))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, "picture: file too large")
			return
		}
		data = body
	}

	if len(data) == 0 {
		data = nil
	}

	if err := h.service.UpdatePicture(id, data); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrImageUnsupported) {
			status = http.StatusUnprocessableEntity
		}
		ctx.AbortWithStatusJSON(status, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *UserHandler) HandleDelete(ctx *gin.Context) {
	var id model.Id
	if err := ctx.ShouldBindUri(&id); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err := h.service.Delete(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}
