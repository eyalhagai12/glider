package http

import (
	backend "glider"
	"net/http"

	"github.com/eyalhagai12/hagio/handler"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createUser struct {
	Username       string `json:"username" binding:"required"`
	Email          string `json:"email" binding:"required"`
	HashedPassword string `json:"hashedPassword" binding:"required"`
}

func (s *Server) RegisterUserRoutes() {
	userGroup := s.apiRoutes.Group("/users")
	userGroup.POST("/register", handler.FromFunc(s.createUser, http.StatusCreated))
}

func (s *Server) createUser(c *gin.Context, request createUser) (*backend.User, error) {
	user := &backend.User{
		ID:             uuid.New(),
		Username:       request.Username,
		Email:          request.Email,
		HashedPassword: request.HashedPassword,
		Role:           backend.RoleUser,
	}

	createdUser, err := s.userService.Create(c.Request.Context(), user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
