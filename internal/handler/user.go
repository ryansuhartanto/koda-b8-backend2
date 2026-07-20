package handler

import (
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
	users := h.service.List()

	ctx.JSON(http.StatusOK, users)
}

func (h *UserHandler) HandleRegister(ctx *gin.Context) {
	var new model.User
	if err := ctx.ShouldBindJSON(&new); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Register(new); err != nil {
		ctx.JSON(http.StatusConflict, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, new)
}

func (h *UserHandler) HandleAuth(ctx *gin.Context) {
	var auth struct {
		Email    string         `json:"email" binding:"required"`
		Password model.Password `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&auth); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Auth(auth.Email, auth.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, *user)
}
