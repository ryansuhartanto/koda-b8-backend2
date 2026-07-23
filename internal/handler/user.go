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
// @Summary  Register a new user
// @Tags     auth
// @Accept   mpfd
// @Param    formData formData model.User true "New user"
// @Param    body     body     model.User true "New user"
// @Success  201      {object} model.AuthResult
// @Failure  400      {object} model.Problem
// @Failure  409      {object} model.Problem
// @Failure  500      {object} model.Problem
// @Router   /auth/register [post]
func (h *UserHandler) HandleRegister(ctx *gin.Context) {
	var new model.User
	if err := ctx.ShouldBind(&new); err != nil {
		model.AbortProblem(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Register(ctx, new)
	if err != nil {
		status := http.StatusInternalServerError
		detail := err.Error()

		if errors.Is(err, service.ErrEmailConflict) {
			status = http.StatusConflict
			detail = "Email already registered"
		}

		model.AbortProblem(ctx, status, detail)
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

// HandleLogin godoc
// @Summary  Log in
// @Tags     auth
// @Param    formData formData model.Credentials true "Credentials"
// @Param    body     body     model.Credentials true "Credentials"
// @Success  200      {object} model.AuthResult
// @Failure  400      {object} model.Problem
// @Failure  401      {object} model.Problem
// @Failure  500      {object} model.Problem
// @Router   /auth/login [post]
func (h *UserHandler) HandleLogin(ctx *gin.Context) {
	var cre model.Credentials
	if err := ctx.ShouldBind(&cre); err != nil {
		model.AbortProblem(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Login(ctx, cre)
	if err != nil {
		status := http.StatusInternalServerError
		detail := err.Error()

		if errors.Is(err, service.ErrEmailUnregistered) {
			status = http.StatusUnauthorized
			detail = "Email is not registered"
		}
		if errors.Is(err, service.ErrPasswordIncorrect) {
			status = http.StatusUnauthorized
			detail = "Password is incorrect"
		}

		model.AbortProblem(ctx, status, detail)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// HandleList godoc
// @Summary  List users
// @Tags     users
// @Produce  json
// @Security BearerAuth
// @Param    limit    query    int false "Max results (default 20, max 100)"
// @Param    offset   query    int false "Results to skip"
// @Success  200      {array}  model.UserIdentified
// @Failure  400      {object} model.Problem
// @Failure  401      {object} model.Problem
// @Failure  500      {object} model.Problem
// @Router   /users/ [get]
func (h *UserHandler) HandleList(ctx *gin.Context) {
	if _, exists := ctx.Get("user.id"); !exists {
		model.AbortProblem(ctx, http.StatusUnauthorized, "")
		return
	}

	var pagination model.Pagination
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		model.AbortProblem(ctx, http.StatusBadRequest, err.Error())
		return
	}

	users, err := h.service.List(ctx, pagination)
	if err != nil {
		model.AbortProblem(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// HandlePatch godoc
// @Summary  Update a user
// @Tags     users
// @Security BearerAuth
// @Param    id       path     int        true "User ID"
// @Param    formData formData model.User true "User fields"
// @Param    body     body     model.User true "User fields"
// @Success  200      {object} model.UserIdentified
// @Failure  400      {object} model.Problem
// @Failure  401      {object} model.Problem
// @Failure  500      {object} model.Problem
// @Router   /users/{id} [patch]
func (h *UserHandler) HandlePatch(ctx *gin.Context) {
	var id model.Id
	if err := ctx.ShouldBindUri(&id); err != nil {
		model.AbortProblem(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if idCtx, exists := ctx.Get("user.id"); !exists || idCtx != id {
		model.AbortProblem(ctx, http.StatusUnauthorized, "")
		return
	}

	var new model.User
	if err := ctx.ShouldBind(&new); err != nil {
		model.AbortProblem(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Edit(ctx, id, new)
	if err != nil {
		model.AbortProblem(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// HandlePutPicture godoc
// @Summary  Upload or clear a user's profile picture
// @Tags     users
// @Security BearerAuth
// @Accept   mpfd
// @Param    id       path     int  true  "User ID"
// @Param    picture  formData file false "Picture file (omit or send empty body to clear)"
// @Success  200
// @Failure  400      {object} model.Problem
// @Failure  401      {object} model.Problem
// @Failure  413      {object} model.Problem
// @Failure  422      {object} model.Problem
// @Failure  500      {object} model.Problem
// @Router   /users/{id}/picture [put]
func (h *UserHandler) HandlePutPicture(ctx *gin.Context) {
	var id model.Id
	if err := ctx.ShouldBindUri(&id); err != nil {
		model.AbortProblem(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if idCtx, exists := ctx.Get("user.id"); !exists || idCtx != id {
		model.AbortProblem(ctx, http.StatusUnauthorized, "")
		return
	}

	var data []byte
	if fileHeader, err := ctx.FormFile("picture"); err == nil {
		if fileHeader.Size > maxPictureSize {
			model.AbortProblem(ctx, http.StatusRequestEntityTooLarge, "picture: file too large")
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			model.AbortProblem(ctx, http.StatusBadRequest, err.Error())
			return
		}
		defer file.Close()

		data, err = io.ReadAll(file)
		if err != nil {
			model.AbortProblem(ctx, http.StatusBadRequest, err.Error())
			return
		}
	} else {
		body, err := io.ReadAll(http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxPictureSize))
		if err != nil {
			model.AbortProblem(ctx, http.StatusRequestEntityTooLarge, "picture: file too large")
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
		model.AbortProblem(ctx, status, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

// HandleDelete godoc
// @Summary  Delete a user
// @Tags     users
// @Security BearerAuth
// @Param    id       path     int true  "User ID"
// @Success  200
// @Failure  400      {object} model.Problem
// @Failure  401      {object} model.Problem
// @Failure  500      {object} model.Problem
// @Router   /users/{id} [delete]
func (h *UserHandler) HandleDelete(ctx *gin.Context) {
	var id model.Id
	if err := ctx.ShouldBindUri(&id); err != nil {
		model.AbortProblem(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if idCtx, exists := ctx.Get("user.id"); !exists || idCtx != id {
		model.AbortProblem(ctx, http.StatusUnauthorized, "")
		return
	}

	err := h.service.Delete(ctx, id)
	if err != nil {
		model.AbortProblem(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}
