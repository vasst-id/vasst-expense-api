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

type workspaceRoutes struct {
	workspaceService services.WorkspaceService
	auth             *middleware.AuthMiddleware
}

func newWorkspaceRoutes(handler *gin.RouterGroup, workspaceService services.WorkspaceService, auth *middleware.AuthMiddleware) {
	r := &workspaceRoutes{
		workspaceService: workspaceService,
		auth:             auth,
	}

	// Workspace management endpoints
	workspaces := handler.Group("/workspaces")
	workspaces.Use(r.auth.AuthRequired())
	{
		workspaces.GET("", r.ListAllWorkspaces)
		workspaces.POST("", r.CreateWorkspace)
		workspaces.GET("/:id", r.GetWorkspaceByID)
		workspaces.PUT("/:id", r.UpdateWorkspace)
		workspaces.DELETE("/:id", r.DeleteWorkspace)
	}
}

// @Summary Get all workspaces
// @Description Get a list of workspaces with optional filtering
// @Tags workspaces
// @Accept json
// @Produce json
// @Param limit query int false "Limit for pagination"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /workspaces [get]
func (r *workspaceRoutes) ListAllWorkspaces(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return // Error response already sent by GetAuthenticatedUserID
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

	// Get workspaces for the authenticated user
	workspaces, err := r.workspaceService.ListAllWorkspaces(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    workspaces,
	})
}

// @Summary Get workspace by ID
// @Description Get a workspace by its ID
// @Tags workspaces
// @Accept json
// @Produce json
// @Param id path string true "Workspace ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /workspaces/{id} [get]
func (r *workspaceRoutes) GetWorkspaceByID(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return // Error response already sent by GetAuthenticatedUserID
	}

	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid workspace ID format",
		})
		return
	}

	workspace, err := r.workspaceService.GetWorkspaceByID(c.Request.Context(), workspaceID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "workspace not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Check if the user owns this workspace
	if workspace.CreatedBy != userID {
		c.JSON(http.StatusForbidden, &entities.ApiResponse{
			Success: false,
			Error:   "access denied: you can only access your own workspaces",
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    workspace,
	})
}

// @Summary Create a new workspace
// @Description Create a new workspace with the provided details
// @Tags workspaces
// @Accept json
// @Produce json
// @Param input body entities.CreateWorkspaceInput true "Workspace details"
// @Success 201 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /workspaces [post]
func (r *workspaceRoutes) CreateWorkspace(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return // Error response already sent by GetAuthenticatedUserID
	}

	var input entities.CreateWorkspaceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	workspace, err := r.workspaceService.CreateWorkspace(c.Request.Context(), userID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "workspace with this name already exists" {
			status = http.StatusConflict
		} else if err.Error() == "workspace name is required" ||
			err.Error() == "workspace type is required" ||
			err.Error() == "currency ID is required" {
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
		Data:    workspace,
	})
}

// @Summary Update a workspace
// @Description Update an existing workspace's details
// @Tags workspaces
// @Accept json
// @Produce json
// @Param id path string true "Workspace ID"
// @Param input body entities.UpdateWorkspaceInput true "Updated workspace details"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /workspaces/{id} [put]
func (r *workspaceRoutes) UpdateWorkspace(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return // Error response already sent by GetAuthenticatedUserID
	}

	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid workspace ID format",
		})
		return
	}

	// Check if workspace exists and user owns it
	existingWorkspace, err := r.workspaceService.GetWorkspaceByID(c.Request.Context(), workspaceID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "workspace not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if existingWorkspace.CreatedBy != userID {
		c.JSON(http.StatusForbidden, &entities.ApiResponse{
			Success: false,
			Error:   "access denied: you can only update your own workspaces",
		})
		return
	}

	var input entities.UpdateWorkspaceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	workspace, err := r.workspaceService.UpdateWorkspace(c.Request.Context(), workspaceID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "workspace not found" {
			status = http.StatusNotFound
		} else if err.Error() == "workspace name already in use" {
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
		Data:    workspace,
	})
}

// @Summary Delete a workspace
// @Description Delete a workspace by its ID
// @Tags workspaces
// @Accept json
// @Produce json
// @Param id path string true "Workspace ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /workspaces/{id} [delete]
func (r *workspaceRoutes) DeleteWorkspace(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return // Error response already sent by GetAuthenticatedUserID
	}

	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid workspace ID format",
		})
		return
	}

	// Check if workspace exists and user owns it
	existingWorkspace, err := r.workspaceService.GetWorkspaceByID(c.Request.Context(), workspaceID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "workspace not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if existingWorkspace.CreatedBy != userID {
		c.JSON(http.StatusForbidden, &entities.ApiResponse{
			Success: false,
			Error:   "access denied: you can only delete your own workspaces",
		})
		return
	}

	err = r.workspaceService.DeleteWorkspace(c.Request.Context(), workspaceID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "workspace not found" {
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
		Message: "Workspace deleted successfully",
	})
}
