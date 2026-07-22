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

// HandleRegister godoc
// @Summary      Register a new user
// @Tags         auth
// @Accept       x-www-form-urlencoded,json,mpfd
// @Produce      json
// @Param        body  body      model.User  true  "New user"
// @Success      201   {object}  service.AuthResult
// @Failure      400   {string}  string  "invalid request body"
// @Failure      409   {string}  string  "email already exists"
// @Failure      500   {string}  string  "internal error"
// @Router       /auth/register [post]
func (h *UserHandler) HandleRegister(ctx *gin.Context) {
	var new model.User
	if err := ctx.ShouldBind(&new); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Register(ctx, new)
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

// HandleLogin godoc
// @Summary      Log in
// @Tags         auth
// @Accept       x-www-form-urlencoded,json,mpfd
// @Produce      json
// @Param        body  body      model.Credentials  true  "Credentials"
// @Success      200   {object}  service.AuthResult
// @Failure      400   {string}  string  "invalid request body"
// @Failure      401   {string}  string  "email not registered"
// @Failure      422   {string}  string  "incorrect password"
// @Failure      500   {string}  string  "internal error"
// @Router       /auth/login [post]
func (h *UserHandler) HandleLogin(ctx *gin.Context) {
	var cre model.Credentials
	if err := ctx.ShouldBind(&cre); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Login(ctx, cre)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrEmailUnregistered) {
			status = http.StatusUnauthorized
		}
		if errors.Is(err, service.ErrPasswordIncorrect) {
			status = http.StatusUnauthorized
		}
		ctx.AbortWithStatusJSON(status, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// HandleList godoc
// @Summary      List users
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   model.UserIdentified
// @Failure      500  {string}  string  "internal error"
// @Router       /users/ [get]
func (h *UserHandler) HandleList(ctx *gin.Context) {
	_, exists := ctx.Get("user.id")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	users, err := h.service.List(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// HandlePatch godoc
// @Summary      Update a user
// @Tags         users
// @Accept       x-www-form-urlencoded,json,mpfd
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int         true  "User ID"
// @Param        body  body      model.User  true  "User fields"
// @Success      200   {object}  model.UserIdentified
// @Failure      400   {string}  string  "invalid request"
// @Failure      500   {string}  string  "internal error"
// @Router       /users/{id} [patch]
func (h *UserHandler) HandlePatch(ctx *gin.Context) {
	var id model.Id
	if err := ctx.ShouldBindUri(&id); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if idCtx, exists := ctx.Get("user.id"); !exists || idCtx != id {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var new model.User
	if err := ctx.ShouldBind(&new); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Edit(ctx, id, new)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// HandlePutPicture godoc
// @Summary      Upload or clear a user's profile picture
// @Tags         users
// @Accept       image/*,mpfd
// @Accept       octet-stream
// @Security     BearerAuth
// @Param        id       path  int   true  "User ID"
// @Param        picture  formData  file  false  "Picture file (omit or send empty body to clear)"
// @Success      200  "OK"
// @Failure      400  {string}  string  "invalid request"
// @Failure      413  {string}  string  "file too large"
// @Failure      422  {string}  string  "unsupported image format"
// @Failure      500  {string}  string  "internal error"
// @Router       /users/{id}/picture [put]
func (h *UserHandler) HandlePutPicture(ctx *gin.Context) {
	var id model.Id
	if err := ctx.ShouldBindUri(&id); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if idCtx, exists := ctx.Get("user.id"); !exists || idCtx != id {
		ctx.AbortWithStatus(http.StatusUnauthorized)
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

	if err := h.service.UpdatePicture(ctx, id, data); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrImageUnsupported) {
			status = http.StatusUnprocessableEntity
		}
		ctx.AbortWithStatusJSON(status, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

// HandleDelete godoc
// @Summary      Delete a user
// @Tags         users
// @Security     BearerAuth
// @Param        id  path  int  true  "User ID"
// @Success      200  "OK"
// @Failure      400  {string}  string  "invalid request"
// @Failure      500  {string}  string  "internal error"
// @Router       /users/{id} [delete]
func (h *UserHandler) HandleDelete(ctx *gin.Context) {
	var id model.Id
	if err := ctx.ShouldBindUri(&id); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if idCtx, exists := ctx.Get("user.id"); !exists || idCtx != id {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err := h.service.Delete(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}
