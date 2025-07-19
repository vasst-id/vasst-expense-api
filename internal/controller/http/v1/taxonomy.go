package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/middleware"
	"github.com/vasst-id/vasst-expense-api/internal/services"
)

type taxonomyRoutes struct {
	taxonomyService services.TaxonomyService
	auth            *middleware.AuthMiddleware
}

func newTaxonomyRoutes(handler *gin.RouterGroup, taxonomyService services.TaxonomyService, auth *middleware.AuthMiddleware) {
	r := &taxonomyRoutes{
		taxonomyService: taxonomyService,
		auth:            auth,
	}

	// Taxonomy endpoints - all require authentication
	taxonomies := handler.Group("/taxonomies")
	taxonomies.Use(auth.AuthRequired())
	{
		taxonomies.GET("/type/:type", r.GetTaxonomiesByType)
	}
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
func (r *taxonomyRoutes) GetTaxonomiesByType(c *gin.Context) {
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
