package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/model"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) HandleList(ctx *gin.Context) {
	users, err := h.service.List()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (h *UserHandler) HandleRegister(ctx *gin.Context) {
	var new model.User
	if err := ctx.ShouldBind(&new); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Register(new)
	if err != nil {
		var status int
		if errors.Is(err, service.ErrEmailConflict) {
			status = http.StatusConflict
		} else {
			status = http.StatusInternalServerError
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
		var status int
		if errors.Is(err, service.ErrEmailUnregistered) {
			status = http.StatusUnauthorized
		} else if errors.Is(err, service.ErrPasswordIncorrect) {
			status = http.StatusUnprocessableEntity
		} else {
			status = http.StatusInternalServerError
		}
		ctx.AbortWithStatusJSON(status, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
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
