package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/light-group/light-ex-backend/internal/modules/auth/service"
)

type Handler struct {
	service *service.Service
}

// New creates a new authentication HTTP handler.
func New(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers all authentication routes.
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {

	auth := router.Group("/auth")

	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.Refresh)
}

// Register handles user registration.
func (h *Handler) Register(c *gin.Context) {

	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})

		return
	}

	user, err := h.service.Register(
		c.Request.Context(),
		req.Email,
		req.Password,
	)

	if err != nil {

		switch {

		case errors.Is(err, service.ErrEmailAlreadyExists):

			c.JSON(http.StatusConflict, ErrorResponse{
				Error: err.Error(),
			})

		case errors.Is(err, service.ErrPasswordTooShort):

			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: err.Error(),
			})

		default:

			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "registration failed",
			})
		}

		return
	}

	c.JSON(http.StatusCreated, RegisterResponse{
		Message: "Registration successful. Please verify your email.",
		User: UserResponse{
			ID:            user.ID.String(),
			Email:         user.Email,
			Status:        string(user.Status),
			EmailVerified: user.EmailVerified,
			CreatedAt:     user.CreatedAt,
		},
	})
}

// Login handles user authentication.
func (h *Handler) Login(c *gin.Context) {

	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})

		return
	}

	result, err := h.service.Login(
		c.Request.Context(),
		req.Email,
		req.Password,
	)

	if err != nil {

		switch {

		case errors.Is(err, service.ErrInvalidCredentials):

			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error: err.Error(),
			})

		case errors.Is(err, service.ErrAccountLocked):

			c.JSON(http.StatusForbidden, ErrorResponse{
				Error: err.Error(),
			})

		case errors.Is(err, service.ErrAccountSuspended):

			c.JSON(http.StatusForbidden, ErrorResponse{
				Error: err.Error(),
			})

		default:

			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "login failed",
			})
		}

		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    result.ExpiresIn,

		User: UserResponse{
			ID:            result.User.ID.String(),
			Email:         result.User.Email,
			Status:        string(result.User.Status),
			EmailVerified: result.User.EmailVerified,
			CreatedAt:     result.User.CreatedAt,
		},
	})
}

func (h *Handler) Refresh(c *gin.Context) {

	var req RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid request body",
		})

		return
	}

	result, err := h.service.Refresh(
		c.Request.Context(),
		req.RefreshToken,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, RefreshTokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		TokenType:    result.TokenType,
		ExpiresIn:    result.ExpiresIn,
	})
}
