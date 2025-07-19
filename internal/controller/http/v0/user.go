package v0

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/middleware"
	"github.com/vasst-id/vasst-expense-api/internal/services"
)

type userAdminRoutes struct {
	userService services.UserService
	auth        *middleware.AuthMiddleware
}

func newUserAdminRoutes(handler *gin.RouterGroup, userService services.UserService, auth *middleware.AuthMiddleware) {
	r := &userAdminRoutes{
		userService: userService,
		auth:        auth,
	}

	// User management endpoints
	users := handler.Group("/users")
	{
		users.GET("", r.ListAllUsers)
		users.POST("", r.CreateUser)
		users.GET("/:id", r.GetUserByID)
		users.PUT("/:id", r.UpdateUser)
		users.DELETE("/:id", r.DeleteUser)
		users.GET("/email/:email", r.GetUserByEmail)
		users.GET("/phone/:phone", r.GetUserByPhoneNumber)
	}

	// Authentication endpoints
	authRoutes := handler.Group("/auth")
	{
		authRoutes.POST("/login", r.Login)
		authRoutes.POST("/forgot-password", r.ForgotPassword)
		authRoutes.POST("/reset-password", r.ResetPassword)
		authRoutes.POST("/change-password", r.ChangePassword)
	}
}

// @Summary Get all users
// @Description Get a list of users with optional filtering
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Limit for pagination"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /users [get]
func (r *userAdminRoutes) ListAllUsers(c *gin.Context) {
	limit := 10
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil {
			limit = val
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil {
			offset = val
		}
	}

	users, err := r.userService.ListAllUsers(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    users,
	})
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
func (r *userAdminRoutes) GetUserByID(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid user ID format",
		})
		return
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
func (r *userAdminRoutes) CreateUser(c *gin.Context) {
	var input entities.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	user, err := r.userService.CreateUser(c.Request.Context(), &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user with this email already exists" ||
			err.Error() == "user with this phone number already exists" {
			status = http.StatusConflict
		} else if err.Error() == "email is required" ||
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
		Data:    user,
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
func (r *userAdminRoutes) UpdateUser(c *gin.Context) {
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

// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /users/{id} [delete]
func (r *userAdminRoutes) DeleteUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid user ID format",
		})
		return
	}

	err = r.userService.DeleteUser(c.Request.Context(), userID)
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
		Message: "User deleted successfully",
	})
}

// @Summary Get user by email
// @Description Get a user by their email address
// @Tags users
// @Accept json
// @Produce json
// @Param email path string true "Email Address"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /users/email/{email} [get]
func (r *userAdminRoutes) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "email is required",
		})
		return
	}

	user, err := r.userService.GetUserByEmail(c.Request.Context(), email)
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

// @Summary Get user by phone number
// @Description Get a user by their phone number
// @Tags users
// @Accept json
// @Produce json
// @Param phone path string true "Phone Number"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /users/phone/{phone} [get]
func (r *userAdminRoutes) GetUserByPhoneNumber(c *gin.Context) {
	phoneNumber := c.Param("phone")
	if phoneNumber == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "phone number is required",
		})
		return
	}

	user, err := r.userService.GetUserByPhoneNumber(c.Request.Context(), phoneNumber)
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
func (r *userAdminRoutes) Login(c *gin.Context) {
	var input entities.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if (input.Email == "" && input.Phone == "") || input.Password == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "email or phone and password are required",
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
func (r *userAdminRoutes) ForgotPassword(c *gin.Context) {
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
func (r *userAdminRoutes) ResetPassword(c *gin.Context) {
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
func (r *userAdminRoutes) ChangePassword(c *gin.Context) {
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
func (r *userAdminRoutes) VerifyPhone(c *gin.Context) {
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
func (r *userAdminRoutes) ResendVerificationCode(c *gin.Context) {
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
func (r *userAdminRoutes) VerifyEmail(c *gin.Context) {
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
func (r *userAdminRoutes) ResendVerificationEmail(c *gin.Context) {
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
