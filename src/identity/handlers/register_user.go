package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jbenzshawel/go-sandbox/identity/app/command"
)

type registerUserRequest struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func (s *HttpServer) RegisterUser(c *gin.Context) {
	var user registerUserRequest
	if err := c.BindJSON(&user); err != nil {
		s.application.Logger.Error(err)
		c.IndentedJSON(http.StatusBadRequest, struct{ Error string }{Error: "invalid JSON"})
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
		s.application.Logger.Error(err)
		c.IndentedJSON(http.StatusBadRequest, struct{ Error string }{Error: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, nil)
}
