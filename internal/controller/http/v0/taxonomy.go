package v0

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/middleware"
	"github.com/vasst-id/vasst-expense-api/internal/services"
)

type taxonomyAdminRoutes struct {
	taxonomyService services.TaxonomyService
	auth            *middleware.AuthMiddleware
}

func newTaxonomyAdminRoutes(handler *gin.RouterGroup, taxonomyService services.TaxonomyService, auth *middleware.AuthMiddleware) {
	r := &taxonomyAdminRoutes{
		taxonomyService: taxonomyService,
		auth:            auth,
	}

	// Taxonomy endpoints - all require authentication
	taxonomies := handler.Group("/taxonomies")
	taxonomies.Use(auth.AuthRequired())
	{
		taxonomies.GET("", r.GetAllTaxonomies)
		taxonomies.POST("", r.CreateTaxonomy)
		taxonomies.GET("/active", r.GetActiveTaxonomies)
		taxonomies.GET("/type/:type", r.GetTaxonomiesByType)
		taxonomies.GET("/type/:type/value/:value", r.GetTaxonomyByTypeAndValue)
		taxonomies.GET("/:id", r.GetTaxonomyByID)
		taxonomies.PUT("/:id", r.UpdateTaxonomy)
		taxonomies.DELETE("/:id", r.DeleteTaxonomy)
	}
}

// @Summary Get all taxonomies
// @Description Get a list of all taxonomies with pagination
// @Tags taxonomies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit for pagination"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /taxonomies [get]
func (r *taxonomyAdminRoutes) GetAllTaxonomies(c *gin.Context) {
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

	taxonomies, err := r.taxonomyService.GetAllTaxonomies(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    taxonomies,
	})
}

// @Summary Create taxonomy
// @Description Create a new taxonomy
// @Tags taxonomies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body entities.CreateTaxonomyRequest true "Taxonomy details"
// @Success 201 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 409 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /taxonomies [post]
func (r *taxonomyAdminRoutes) CreateTaxonomy(c *gin.Context) {
	var input entities.CreateTaxonomyRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	taxonomy, err := r.taxonomyService.CreateTaxonomy(c.Request.Context(), &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "taxonomy with this type and value already exists" {
			status = http.StatusConflict
		} else if err.Error() == "label is required" ||
			err.Error() == "value is required" ||
			err.Error() == "type is required" ||
			err.Error() == "type label is required" {
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
		Data:    taxonomy,
	})
}

// @Summary Get active taxonomies
// @Description Get all active taxonomies with pagination
// @Tags taxonomies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit for pagination"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /taxonomies/active [get]
func (r *taxonomyAdminRoutes) GetActiveTaxonomies(c *gin.Context) {
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

	taxonomies, err := r.taxonomyService.GetActiveTaxonomies(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    taxonomies,
	})
}

// @Summary Get taxonomies by type
// @Description Get taxonomies by type with pagination
// @Tags taxonomies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type path string true "Taxonomy type"
// @Param limit query int false "Limit for pagination"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /taxonomies/type/{type} [get]
func (r *taxonomyAdminRoutes) GetTaxonomiesByType(c *gin.Context) {
	taxonomyType := c.Param("type")
	if taxonomyType == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "taxonomy type is required",
		})
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

	taxonomies, err := r.taxonomyService.GetTaxonomiesByType(c.Request.Context(), taxonomyType, limit, offset)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "taxonomy type is required" {
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
		Data:    taxonomies,
	})
}

// @Summary Get taxonomy by type and value
// @Description Get a taxonomy by its type and value
// @Tags taxonomies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type path string true "Taxonomy type"
// @Param value path string true "Taxonomy value"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /taxonomies/type/{type}/value/{value} [get]
func (r *taxonomyAdminRoutes) GetTaxonomyByTypeAndValue(c *gin.Context) {
	taxonomyType := c.Param("type")
	value := c.Param("value")

	if taxonomyType == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "taxonomy type is required",
		})
		return
	}

	if value == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "value is required",
		})
		return
	}

	taxonomy, err := r.taxonomyService.GetTaxonomyByTypeAndValue(c.Request.Context(), taxonomyType, value)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "taxonomy not found" {
			status = http.StatusNotFound
		} else if err.Error() == "taxonomy type is required" || err.Error() == "value is required" {
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
		Data:    taxonomy,
	})
}

// @Summary Get taxonomy by ID
// @Description Get a taxonomy by its ID
// @Tags taxonomies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Taxonomy ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /taxonomies/{id} [get]
func (r *taxonomyAdminRoutes) GetTaxonomyByID(c *gin.Context) {
	taxonomyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid taxonomy ID format",
		})
		return
	}

	taxonomy, err := r.taxonomyService.GetTaxonomyByID(c.Request.Context(), taxonomyID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "taxonomy not found" {
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
		Data:    taxonomy,
	})
}

// @Summary Update taxonomy
// @Description Update an existing taxonomy
// @Tags taxonomies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Taxonomy ID"
// @Param input body entities.UpdateTaxonomyRequest true "Updated taxonomy details"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 409 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /taxonomies/{id} [put]
func (r *taxonomyAdminRoutes) UpdateTaxonomy(c *gin.Context) {
	taxonomyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid taxonomy ID format",
		})
		return
	}

	var input entities.UpdateTaxonomyRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	taxonomy, err := r.taxonomyService.UpdateTaxonomy(c.Request.Context(), taxonomyID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "taxonomy not found" {
			status = http.StatusNotFound
		} else if err.Error() == "taxonomy with this type and value already exists" {
			status = http.StatusConflict
		} else if err.Error() == "label is required" ||
			err.Error() == "value is required" ||
			err.Error() == "type is required" ||
			err.Error() == "type label is required" {
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
		Data:    taxonomy,
	})
}

// @Summary Delete taxonomy
// @Description Delete a taxonomy by its ID (soft delete)
// @Tags taxonomies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Taxonomy ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /taxonomies/{id} [delete]
func (r *taxonomyAdminRoutes) DeleteTaxonomy(c *gin.Context) {
	taxonomyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid taxonomy ID format",
		})
		return
	}

	err = r.taxonomyService.DeleteTaxonomy(c.Request.Context(), taxonomyID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "taxonomy not found" {
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
		Message: "Taxonomy deleted successfully",
	})
}
