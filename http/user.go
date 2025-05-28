package http

import (
	backend "glider"
	"net/http"

	"github.com/eyalhagai12/hagio/handler"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) RegisterUserRoutes() {
	userGroup := s.apiRoutes.Group("/users")
	userGroup.POST("/register", handler.FromFunc(s.regiser, http.StatusCreated))
	userGroup.POST("/login", handler.FromFunc(s.login, http.StatusOK))
}

func (s *Server) regiser(c *gin.Context, request createUser) (*backend.User, error) {
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

func (s *Server) login(c *gin.Context, request loginUser) (*backend.User, error) {
	user, err := s.userService.GetByUsername(c.Request.Context(), request.Username)
	if err != nil {
		return nil, err
	}
	
	s.logger.Debug("logging in user", "username", user.Username, "password", user.HashedPassword, "request password", request.HashedPassword)

	if user.HashedPassword != request.HashedPassword {
		return nil, backend.ErrUnauthorized
	}

	return user, nil
}
