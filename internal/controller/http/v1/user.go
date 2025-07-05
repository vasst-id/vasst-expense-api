package v1

import (
	"net/http"
	"strconv"

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

	// Organization-scoped endpoints
	org := handler.Group("/users")
	{
		org.GET("", r.ListUsersByOrganization)
		org.GET("/:id", r.GetUserByIDAndOrganization)
		org.GET("/phone/:phone", r.GetUserByPhoneNumberAndOrganization)
		org.GET("/username/:username", r.GetUserByUsernameAndOrganization)
		org.POST("/login", r.Login)
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
func (r *userRoutes) ListAllUsers(c *gin.Context) {
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
func (r *userRoutes) GetUserByID(c *gin.Context) {
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
// @Param input body struct{ PhoneNumber string `json:"phone_number" binding:"required"`; FullName string `json:"full_name" binding:"required"` } true "User details"
// @Success 201 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /users [post]
func (r *userRoutes) CreateUser(c *gin.Context) {
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
		if err.Error() == "user with this phone number already exists" || err.Error() == "user with this username already exists" {
			status = http.StatusConflict
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
// @Param input body struct{ PhoneNumber string `json:"phone_number" binding:"required"`; FullName string `json:"full_name" binding:"required"` } true "Updated user details"
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
		} else if err.Error() == "phone number already in use" || err.Error() == "username already in use" {
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
func (r *userRoutes) DeleteUser(c *gin.Context) {
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
		Data:    "User deleted successfully",
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
func (r *userRoutes) GetUserByPhoneNumber(c *gin.Context) {
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

func (r *userRoutes) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{Success: false, Error: "username is required"})
		return
	}
	user, err := r.userService.GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, &entities.ApiResponse{Success: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, &entities.ApiResponse{Success: true, Data: user})
}

// @Summary Get users by organization
// @Description Get a list of users for a specific organization
// @Tags users
// @Accept json
// @Produce json
// @Param orgID path string true "Organization ID"
// @Param limit query int false "Limit for pagination"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /organizations/{orgID}/users [get]
func (r *userRoutes) ListUsersByOrganization(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("orgID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{Success: false, Error: "invalid organization ID format"})
		return
	}
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
	users, err := r.userService.ListUsersByOrganization(c.Request.Context(), orgID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{Success: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, &entities.ApiResponse{Success: true, Data: users})
}

// @Summary Get user by ID and organization
// @Description Get a user by their ID and organization
// @Tags users
// @Accept json
// @Produce json
// @Param orgID path string true "Organization ID"
// @Param id path string true "User ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /organizations/{orgID}/users/{id} [get]
func (r *userRoutes) GetUserByIDAndOrganization(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("orgID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{Success: false, Error: "invalid organization ID format"})
		return
	}
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{Success: false, Error: "invalid user ID format"})
		return
	}
	user, err := r.userService.GetUserByIDAndOrganization(c.Request.Context(), userID, orgID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, &entities.ApiResponse{Success: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, &entities.ApiResponse{Success: true, Data: user})
}

// @Summary Get user by phone number and organization
// @Description Get a user by their phone number and organization
// @Tags users
// @Accept json
// @Produce json
// @Param orgID path string true "Organization ID"
// @Param phone path string true "Phone Number"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /organizations/{orgID}/users/phone/{phone} [get]
func (r *userRoutes) GetUserByPhoneNumberAndOrganization(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("orgID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{Success: false, Error: "invalid organization ID format"})
		return
	}
	phoneNumber := c.Param("phone")
	if phoneNumber == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{Success: false, Error: "phone number is required"})
		return
	}
	user, err := r.userService.GetUserByPhoneNumberAndOrganization(c.Request.Context(), phoneNumber, orgID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, &entities.ApiResponse{Success: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, &entities.ApiResponse{Success: true, Data: user})
}

// @Summary Get user by username and organization
// @Description Get a user by their username and organization
// @Tags users
// @Accept json
// @Produce json
// @Param orgID path string true "Organization ID"
// @Param username path string true "Username"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /organizations/{orgID}/users/username/{username} [get]
func (r *userRoutes) GetUserByUsernameAndOrganization(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("orgID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{Success: false, Error: "invalid organization ID format"})
		return
	}
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{Success: false, Error: "username is required"})
		return
	}
	user, err := r.userService.GetUserByUsernameAndOrganization(c.Request.Context(), username, orgID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, &entities.ApiResponse{Success: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, &entities.ApiResponse{Success: true, Data: user})
}

// @Summary Login
// @Description Login a user
// @Tags users
// @Accept json
// @Produce json
// @Param input body entities.LoginInput true "Login details"
// @Success 200 {object} entities.ApiResponse
func (r *userRoutes) Login(c *gin.Context) {
	var input entities.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if input.Username == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "username and password are required",
		})
		return
	}

	user, err := r.userService.Login(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{Success: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, &entities.ApiResponse{Success: true, Data: user})
}
