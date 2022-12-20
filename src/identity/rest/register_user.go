package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/jbenzshawel/go-sandbox/common/cerror"
	"github.com/jbenzshawel/go-sandbox/identity/app/command"
	"github.com/jbenzshawel/go-sandbox/identity/app/query"
)

type registerUserRequest struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type registerUserResponse struct {
	ID        int32     `json:"id"`
	UUID      uuid.UUID `json:"uuid"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
}

func (s *HttpServer) RegisterUser(c *gin.Context) {
	var user registerUserRequest
	if err := c.BindJSON(&user); err != nil {
		s.application.Logger.Error(err)
		c.IndentedJSON(http.StatusBadRequest, cerror.NewValidationError("invalid JSON", nil))
		return
	}

	err := s.application.Commands.RegisterUser.Handle(c, command.RegisterUser{
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Email:           user.Email,
		Password:        user.Password,
		ConfirmPassword: user.ConfirmPassword,
	})
	if err != nil {
		cerror.HandleValidationError(c, err)
		return
	}

	createdUser, err := s.application.Queries.UserByEmail.Handle(c, query.UserByEmail{Email: user.Email})
	if err != nil {
		cerror.HandleValidationError(c, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, &registerUserResponse{
		ID:        createdUser.ID,
		UUID:      createdUser.UUID,
		FirstName: createdUser.FirstName,
		LastName:  createdUser.LastName,
		Email:     createdUser.Email,
	})
}
