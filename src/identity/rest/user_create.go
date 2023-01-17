package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/jbenzshawel/go-sandbox/common/cerror"
	crest "github.com/jbenzshawel/go-sandbox/common/rest"
	"github.com/jbenzshawel/go-sandbox/identity/app/command"
	"github.com/jbenzshawel/go-sandbox/identity/app/query"
)

type createUserRequest struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type createUserResponse struct {
	ID        int32     `json:"id"`
	UUID      uuid.UUID `json:"uuid"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
}

func (h *HttpHandler) CreateUser(ctx *gin.Context) {
	var user createUserRequest
	if err := ctx.BindJSON(&user); err != nil {
		h.app.Logger.Error(err)
		ctx.JSON(http.StatusBadRequest, cerror.NewValidationError("invalid JSON", nil))
		return
	}

	err := h.app.Commands.CreateUser.Handle(ctx, command.UserCreate{
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Email:           user.Email,
		Password:        user.Password,
		ConfirmPassword: user.ConfirmPassword,
	})
	if err != nil {
		crest.HandleErrorResponse(ctx, err)
		return
	}

	createdUser, err := h.app.Queries.UserByEmail.Handle(ctx, query.UserByEmail{Email: user.Email})
	if err != nil {
		crest.HandleErrorResponse(ctx, err)
		return
	}

	err = h.app.Commands.SendVerificationEmail.Handle(ctx, command.SendVerificationEmail{
		UserUUID:  createdUser.UUID(),
		FirstName: createdUser.FirstName(),
		Email:     createdUser.Email(),
	})
	if err != nil {
		// just log error here since the user has been created. verification email can be resent
		h.app.Logger.WithError(err).Error("failed to send verification email during create user")
	}

	ctx.JSON(http.StatusCreated, &createUserResponse{
		ID:        createdUser.ID(),
		UUID:      createdUser.UUID(),
		FirstName: createdUser.FirstName(),
		LastName:  createdUser.LastName(),
		Email:     createdUser.Email(),
	})
}
