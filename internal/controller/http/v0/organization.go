package v0

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/middleware"
	"github.com/vasst-id/vasst-expense-api/internal/services"
)

type organizationRoutes struct {
	orgService services.OrganizationService
	auth       *middleware.AuthMiddleware
}

func newOrganizationRoutes(handler *gin.RouterGroup, orgService services.OrganizationService, auth *middleware.AuthMiddleware) {
	r := &organizationRoutes{
		orgService: orgService,
		auth:       auth,
	}

	h := handler.Group("/organizations")
	h.Use(r.auth.SuperAdminRequired()) // Apply SuperAdmin protection to all organization endpoints
	{
		h.GET("", r.ListOrganizations)
		h.GET("/:id", r.GetOrganizationByID)
		h.GET("/code/:code", r.GetOrganizationByCode)
		h.POST("", r.CreateOrganization)
		h.PUT("/:id", r.UpdateOrganization)
		h.DELETE("/:id", r.DeleteOrganization)

		// Category routes
		h.GET("/categories", r.ListCategories)
		h.GET("/categories/:id", r.GetCategoryByID)
		h.POST("/categories", r.CreateCategory)
		h.PUT("/categories/:id", r.UpdateCategory)
		h.DELETE("/categories/:id", r.DeleteCategory)

		// Setting routes
		h.GET("/:id/settings", r.GetSettingByOrgID)
		h.PUT("/:id/settings", r.UpdateSetting)

		// Knowledge routes
		h.GET("/:id/knowledge", r.ListKnowledgeByOrgID)
		h.GET("/knowledge/:id", r.GetKnowledgeByID)
		h.POST("/:id/knowledge", r.CreateKnowledge)
		h.POST("/:id/knowledge/with-file", r.CreateKnowledgeWithFile)
		h.PUT("/knowledge/:id", r.UpdateKnowledge)
		h.PUT("/knowledge/:id/with-file", r.UpdateKnowledgeWithFile)
		h.DELETE("/knowledge/:id", r.DeleteKnowledge)

		// File upload routes
		h.POST("/:id/upload", r.UploadFile)
		h.POST("/:id/upload/multiple", r.UploadMultipleFiles)

		// Model routes
		h.GET("/:id/models", r.ListModelsByOrgID)

		// Integration routes
		h.GET("/:id/integrations", r.ListIntegrationsByOrgID)
		h.GET("/integrations/:id", r.GetIntegrationByID)
		h.POST("/:id/integrations", r.CreateIntegration)
		h.PUT("/integrations/:id", r.UpdateIntegration)
		h.DELETE("/integrations/:id", r.DeleteIntegration)
	}
}

// @Summary Get all organizations
// @Description Get a list of organizations with optional filtering
// @Tags organizations
// @Accept json
// @Produce json
// @Param limit query int false "Limit for pagination"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /organizations [get]
func (r *organizationRoutes) ListOrganizations(c *gin.Context) {
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

	organizations, err := r.orgService.ListOrganizations(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    organizations,
	})
}

// @Summary Get organization by ID
// @Description Get an organization by their ID
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path string true "Organization ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /organizations/{id} [get]
func (r *organizationRoutes) GetOrganizationByID(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}

	org, err := r.orgService.GetOrganizationByID(c.Request.Context(), orgID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "organization not found" {
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
		Data:    org,
	})
}

// @Summary Get organization by code
// @Description Get an organization by their code
// @Tags organizations
// @Accept json
// @Produce json
// @Param code path string true "Organization Code"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /organizations/code/{code} [get]
func (r *organizationRoutes) GetOrganizationByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "organization code is required",
		})
		return
	}

	org, err := r.orgService.GetOrganizationByCode(c.Request.Context(), code)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "organization not found" {
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
		Data:    org,
	})
}

// @Summary Create a new organization
// @Description Create a new organization with the provided details
// @Tags organizations
// @Accept json
// @Produce json
// @Param input body entities.CreateOrganizationInput true "Organization details"
// @Success 201 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /organizations [post]
func (r *organizationRoutes) CreateOrganization(c *gin.Context) {
	var input entities.CreateOrganizationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	org, err := r.orgService.CreateOrganization(c.Request.Context(), &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "organization with this code already exists" {
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
		Data:    org,
	})
}

// @Summary Update an organization
// @Description Update an existing organization's details
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path string true "Organization ID"
// @Param input body entities.UpdateOrganizationInput true "Updated organization details"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /organizations/{id} [put]
func (r *organizationRoutes) UpdateOrganization(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}

	var input entities.UpdateOrganizationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	org, err := r.orgService.UpdateOrganization(c.Request.Context(), orgID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "organization not found" {
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
		Data:    org,
	})
}

// @Summary Delete an organization
// @Description Delete an organization by their ID
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path string true "Organization ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /organizations/{id} [delete]
func (r *organizationRoutes) DeleteOrganization(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}

	err = r.orgService.DeleteOrganization(c.Request.Context(), orgID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "organization not found" {
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
		Data:    "Organization deleted successfully",
	})
}

// Category handlers
func (r *organizationRoutes) ListCategories(c *gin.Context) {
	categories, err := r.orgService.ListCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    categories,
	})
}

func (r *organizationRoutes) GetCategoryByID(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid category ID format",
		})
		return
	}

	category, err := r.orgService.GetCategoryByID(c.Request.Context(), categoryID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "category not found" {
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
		Data:    category,
	})
}

func (r *organizationRoutes) CreateCategory(c *gin.Context) {
	var input entities.CreateOrganizationCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	category, err := r.orgService.CreateCategory(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, &entities.ApiResponse{
		Success: true,
		Data:    category,
	})
}

func (r *organizationRoutes) UpdateCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid category ID format",
		})
		return
	}

	var input entities.UpdateOrganizationCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	category, err := r.orgService.UpdateCategory(c.Request.Context(), categoryID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "category not found" {
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
		Data:    category,
	})
}

func (r *organizationRoutes) DeleteCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid category ID format",
		})
		return
	}

	err = r.orgService.DeleteCategory(c.Request.Context(), categoryID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "category not found" {
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
		Data:    "Category deleted successfully",
	})
}

// Setting handlers
func (r *organizationRoutes) GetSettingByOrgID(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}

	setting, err := r.orgService.GetSettingByOrgID(c.Request.Context(), orgID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "setting not found" {
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
		Data:    setting,
	})
}

func (r *organizationRoutes) UpdateSetting(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}

	var input entities.UpdateOrganizationSettingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	setting, err := r.orgService.UpdateSetting(c.Request.Context(), orgID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "setting not found" {
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
		Data:    setting,
	})
}

// Knowledge handlers
func (r *organizationRoutes) ListKnowledgeByOrgID(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}

	knowledge, err := r.orgService.ListKnowledgeByOrgID(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    knowledge,
	})
}

func (r *organizationRoutes) GetKnowledgeByID(c *gin.Context) {
	knowledgeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid knowledge ID format",
		})
		return
	}

	knowledge, err := r.orgService.GetKnowledgeByID(c.Request.Context(), knowledgeID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "knowledge not found" {
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
		Data:    knowledge,
	})
}

func (r *organizationRoutes) CreateKnowledge(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}

	var input entities.CreateOrganizationKnowledgeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	input.OrganizationID = orgID
	knowledge, err := r.orgService.CreateKnowledge(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, &entities.ApiResponse{
		Success: true,
		Data:    knowledge,
	})
}

func (r *organizationRoutes) CreateKnowledgeWithFile(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}

	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "file is required",
		})
		return
	}

	// Parse form data
	knowledgeType, err := strconv.Atoi(c.PostForm("knowledge_type"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid knowledge type",
		})
		return
	}

	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "title is required",
		})
		return
	}

	content := c.PostForm("content")
	if content == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "content is required",
		})
		return
	}

	description := c.PostForm("description")
	isActive := c.PostForm("is_active") == "true"

	input := &entities.CreateOrganizationKnowledgeInput{
		OrganizationID: orgID,
		KnowledgeType:  knowledgeType,
		Title:          title,
		Content:        content,
		Description:    &description,
		IsActive:       isActive,
	}

	knowledge, err := r.orgService.CreateKnowledgeWithFile(c.Request.Context(), input, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, &entities.ApiResponse{
		Success: true,
		Data:    knowledge,
	})
}

func (r *organizationRoutes) UpdateKnowledge(c *gin.Context) {
	knowledgeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid knowledge ID format",
		})
		return
	}

	var input entities.UpdateOrganizationKnowledgeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	knowledge, err := r.orgService.UpdateKnowledge(c.Request.Context(), knowledgeID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "knowledge not found" {
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
		Data:    knowledge,
	})
}

func (r *organizationRoutes) UpdateKnowledgeWithFile(c *gin.Context) {
	knowledgeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid knowledge ID format",
		})
		return
	}

	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "file is required",
		})
		return
	}

	// Parse form data
	knowledgeType, err := strconv.Atoi(c.PostForm("knowledge_type"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid knowledge type",
		})
		return
	}

	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "title is required",
		})
		return
	}

	content := c.PostForm("content")
	if content == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "content is required",
		})
		return
	}

	description := c.PostForm("description")
	isActive := c.PostForm("is_active") == "true"

	input := &entities.UpdateOrganizationKnowledgeInput{
		KnowledgeType: knowledgeType,
		Title:         &title,
		Content:       content,
		Description:   &description,
		IsActive:      &isActive,
	}

	knowledge, err := r.orgService.UpdateKnowledgeWithFile(c.Request.Context(), knowledgeID, input, file)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "knowledge not found" {
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
		Data:    knowledge,
	})
}

func (r *organizationRoutes) DeleteKnowledge(c *gin.Context) {
	knowledgeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid knowledge ID format",
		})
		return
	}

	err = r.orgService.DeleteKnowledge(c.Request.Context(), knowledgeID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "knowledge not found" {
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
		Data:    "Knowledge deleted successfully",
	})
}

// Model handlers
func (r *organizationRoutes) ListModelsByOrgID(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}

	models, err := r.orgService.ListModelsByOrgID(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    models,
	})
}

// Integration handlers
func (r *organizationRoutes) ListIntegrationsByOrgID(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}

	integrations, err := r.orgService.ListIntegrationsByOrgID(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    integrations,
	})
}

func (r *organizationRoutes) GetIntegrationByID(c *gin.Context) {
	integrationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid integration ID format",
		})
		return
	}

	integration, err := r.orgService.GetIntegrationByID(c.Request.Context(), integrationID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "integration not found" {
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
		Data:    integration,
	})
}

func (r *organizationRoutes) CreateIntegration(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}

	var input entities.CreateOrganizationIntegrationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	input.OrganizationID = orgID
	integration, err := r.orgService.CreateIntegration(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, &entities.ApiResponse{
		Success: true,
		Data:    integration,
	})
}

func (r *organizationRoutes) UpdateIntegration(c *gin.Context) {
	integrationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid integration ID format",
		})
		return
	}

	var input entities.UpdateOrganizationIntegrationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	integration, err := r.orgService.UpdateIntegration(c.Request.Context(), integrationID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "integration not found" {
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
		Data:    integration,
	})
}

func (r *organizationRoutes) DeleteIntegration(c *gin.Context) {
	integrationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid integration ID format",
		})
		return
	}

	err = r.orgService.DeleteIntegration(c.Request.Context(), integrationID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "integration not found" {
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
		Data:    "Integration deleted successfully",
	})
}

// File upload handlers
func (r *organizationRoutes) UploadFile(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}

	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "file is required",
		})
		return
	}

	// Validate file size (e.g., max 10MB)
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "file size too large, maximum 10MB allowed",
		})
		return
	}

	// Validate file type
	allowedTypes := map[string]bool{
		"application/pdf":    true,
		"image/jpeg":         true,
		"image/png":          true,
		"image/gif":          true,
		"text/plain":         true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	}

	if !allowedTypes[file.Header.Get("Content-Type")] {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "file type not allowed. Allowed types: PDF, JPEG, PNG, GIF, TXT, DOC, DOCX",
		})
		return
	}

	// Generate a unique ID for the file
	fileID := uuid.New()

	// Upload file to Google Cloud Storage
	uploadResult, err := r.orgService.UploadFile(c.Request.Context(), orgID, fileID, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    uploadResult,
	})
}

func (r *organizationRoutes) UploadMultipleFiles(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "failed to parse form data",
		})
		return
	}

	files := c.Request.MultipartForm.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "no files provided",
		})
		return
	}

	if len(files) > 10 {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "too many files, maximum 10 files allowed",
		})
		return
	}

	var uploadResults []*entities.FileUploadResult
	var errors []string

	for _, file := range files {
		// Validate file size (e.g., max 10MB per file)
		if file.Size > 10*1024*1024 {
			errors = append(errors, fmt.Sprintf("File %s: size too large", file.Filename))
			continue
		}

		// Validate file type
		allowedTypes := map[string]bool{
			"application/pdf":    true,
			"image/jpeg":         true,
			"image/png":          true,
			"image/gif":          true,
			"text/plain":         true,
			"application/msword": true,
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		}

		if !allowedTypes[file.Header.Get("Content-Type")] {
			errors = append(errors, fmt.Sprintf("File %s: type not allowed", file.Filename))
			continue
		}

		// Generate a unique ID for the file
		fileID := uuid.New()

		// Upload file to Google Cloud Storage
		uploadResult, err := r.orgService.UploadFile(c.Request.Context(), orgID, fileID, file)
		if err != nil {
			errors = append(errors, fmt.Sprintf("File %s: %s", file.Filename, err.Error()))
			continue
		}

		uploadResults = append(uploadResults, uploadResult)
	}

	response := map[string]interface{}{
		"uploaded_files": uploadResults,
		"total_files":    len(files),
		"success_count":  len(uploadResults),
		"error_count":    len(errors),
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    response,
	})
}
