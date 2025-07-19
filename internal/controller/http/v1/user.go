package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/middleware"
	"github.com/vasst-id/vasst-expense-api/internal/services"
)

type userRoutes struct {
	userService services.UserService
	auth        *middleware.AuthMiddleware
}

func newUserRoutes(handler *gin.RouterGroup, userService services.UserService, auth *middleware.AuthMiddleware) {
	r := &userRoutes{
		userService: userService,
		auth:        auth,
	}

	// User management endpoints
	users := handler.Group("/users")
	{
		users.GET("/", r.auth.AuthRequired(), r.GetUserByID)
		users.PUT("/", r.auth.AuthRequired(), r.UpdateUser)
	}

	// Authentication endpoints
	authRoutes := handler.Group("/auth")
	{
		authRoutes.POST("/login", r.Login)
		authRoutes.POST("/register", r.Register)
		authRoutes.POST("/forgot-password", r.ForgotPassword)
		authRoutes.POST("/reset-password", r.ResetPassword)
		authRoutes.POST("/change-password", r.ChangePassword)
		authRoutes.POST("/verify-phone", r.VerifyPhone)
		authRoutes.POST("/resend-verification-code", r.ResendVerificationCode)
		authRoutes.POST("/verify-email", r.VerifyEmail)
		authRoutes.POST("/resend-verification-email", r.ResendVerificationEmail)
	}
}

// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /users/{id} [get]
func (r *userRoutes) GetUserByID(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return // Error response already sent by GetAuthenticatedUserID
	}

	user, err := r.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    user,
	})
}

// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param input body entities.CreateUserInput true "User details"
// @Success 201 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /users [post]
func (r *userRoutes) Register(c *gin.Context) {
	var input entities.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	loginResponse, err := r.userService.CreateUser(c.Request.Context(), &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user with this email already exists" ||
			err.Error() == "user with this phone number already exists" {
			status = http.StatusConflict
		} else if err.Error() == "password must be exactly 6 digits" ||
			err.Error() == "email is required" ||
			err.Error() == "phone number is required" ||
			err.Error() == "first name is required" ||
			err.Error() == "last name is required" ||
			err.Error() == "password is required" {
			status = http.StatusBadRequest
		}
		c.JSON(status, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, &entities.ApiResponse{
		Success: true,
		Data:    loginResponse,
		Message: "User registered successfully",
	})
}

// @Summary Update a user
// @Description Update an existing user's details
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param input body entities.UpdateUserInput true "Updated user details"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /users/{id} [put]
func (r *userRoutes) UpdateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid user ID format",
		})
		return
	}

	var input entities.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	user, err := r.userService.UpdateUser(c.Request.Context(), userID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		} else if err.Error() == "phone number already in use" ||
			err.Error() == "email already in use" {
			status = http.StatusConflict
		}
		c.JSON(status, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    user,
	})
}

// @Summary Login
// @Description Authenticate a user with email/phone and password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body entities.LoginRequest true "Login details"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Router /auth/login [post]
func (r *userRoutes) Login(c *gin.Context) {
	var input entities.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if input.PhoneNumber == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "phone number and password are required",
		})
		return
	}

	loginResponse, err := r.userService.Login(c.Request.Context(), &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		} else if err.Error() == "password is incorrect" {
			status = http.StatusUnauthorized
		} else if err.Error() == "user account is inactive" {
			status = http.StatusForbidden
		} else if err.Error() == "email or phone is required" {
			status = http.StatusBadRequest
		}
		c.JSON(status, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    loginResponse,
	})
}

// @Summary Forgot Password
// @Description Initiate password reset process
// @Tags auth
// @Accept json
// @Produce json
// @Param input body entities.ForgotPasswordRequest true "Forgot password request"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Router /auth/forgot-password [post]
func (r *userRoutes) ForgotPassword(c *gin.Context) {
	var input entities.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	err := r.userService.ForgotPassword(c.Request.Context(), &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Message: "Password reset instructions sent to your email",
	})
}

// @Summary Reset Password
// @Description Reset password using token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body entities.ResetPasswordRequest true "Reset password request"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Router /auth/reset-password [post]
func (r *userRoutes) ResetPassword(c *gin.Context) {
	var input entities.ResetPasswordRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	err := r.userService.ResetPassword(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Message: "Password reset successfully",
	})
}

// @Summary Change Password
// @Description Change user's password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body entities.ChangePasswordRequest true "Change password request"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 501 {object} entities.ApiResponse
// @Router /auth/change-password [post]
func (r *userRoutes) ChangePassword(c *gin.Context) {
	var input entities.ChangePasswordRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	err := r.userService.ChangePassword(c.Request.Context(), &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "change password not implemented" {
			status = http.StatusNotImplemented
		}
		c.JSON(status, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}

// @Summary Verify Phone
// @Description Verify phone number with code
// @Tags auth
// @Accept json
// @Produce json
// @Param input body entities.VerifyPhoneRequest true "Phone verification request"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Router /auth/verify-phone [post]
func (r *userRoutes) VerifyPhone(c *gin.Context) {
	var input entities.VerifyPhoneRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	err := r.userService.VerifyPhone(c.Request.Context(), &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Message: "Phone number verified successfully",
	})
}

// @Summary Resend Verification Code
// @Description Resend SMS verification code
// @Tags auth
// @Accept json
// @Produce json
// @Param input body entities.ResendVerificationCodeRequest true "Resend code request"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Router /auth/resend-verification-code [post]
func (r *userRoutes) ResendVerificationCode(c *gin.Context) {
	var input entities.ResendVerificationCodeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	err := r.userService.ResendVerificationCode(c.Request.Context(), &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Message: "Verification code sent successfully",
	})
}

// @Summary Verify Email
// @Description Verify email address with token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body entities.VerifyEmailRequest true "Email verification request"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Router /auth/verify-email [post]
func (r *userRoutes) VerifyEmail(c *gin.Context) {
	var input entities.VerifyEmailRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	err := r.userService.VerifyEmail(c.Request.Context(), &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Message: "Email verified successfully",
	})
}

// @Summary Resend Verification Email
// @Description Resend email verification link
// @Tags auth
// @Accept json
// @Produce json
// @Param input body entities.ResendVerificationEmailRequest true "Resend email request"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Router /auth/resend-verification-email [post]
func (r *userRoutes) ResendVerificationEmail(c *gin.Context) {
	var input entities.ResendVerificationEmailRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	err := r.userService.ResendVerificationEmail(c.Request.Context(), &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Message: "Verification email sent successfully",
	})
}
