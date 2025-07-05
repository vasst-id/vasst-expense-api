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

type leadsRoutes struct {
	leadsService services.LeadsService
	auth         *middleware.AuthMiddleware
}

func newLeadsRoutes(handler *gin.RouterGroup, leadsService services.LeadsService, auth *middleware.AuthMiddleware) {
	r := &leadsRoutes{
		leadsService: leadsService,
		auth:         auth,
	}

	// Public endpoint (no authentication)
	leads := handler.Group("/leads")
	{
		leads.POST("", r.CreateLead) // Public endpoint for lead creation
	}

	// Protected endpoints (superadmin only)
	adminLeads := handler.Group("/admin/leads")
	adminLeads.Use(r.auth.SuperAdminRequired())
	{
		adminLeads.GET("", r.ListAllLeads)
		adminLeads.GET("/:id", r.GetLeadByID)
		adminLeads.PUT("/:id", r.UpdateLead)
		adminLeads.GET("/phone/:phone", r.GetLeadByPhoneNumber)
		adminLeads.GET("/email/:email", r.GetLeadByEmail)
	}
}

// @Summary Create a new lead
// @Description Create a new lead (public endpoint)
// @Tags leads
// @Accept json
// @Produce json
// @Param input body entities.CreateLeadInput true "Lead details"
// @Success 201 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 409 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /leads [post]
func (r *leadsRoutes) CreateLead(c *gin.Context) {
	var input entities.CreateLeadInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	lead, err := r.leadsService.CreateLead(c.Request.Context(), &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "lead with this phone number already exists" || err.Error() == "lead with this email already exists" {
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
		Data:    lead,
	})
}

// @Summary Get all leads
// @Description Get a list of leads with pagination (superadmin only)
// @Tags leads
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit for pagination"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /admin/leads [get]
func (r *leadsRoutes) ListAllLeads(c *gin.Context) {
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

	leads, err := r.leadsService.ListAllLeads(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    leads,
	})
}

// @Summary Get lead by ID
// @Description Get a lead by their ID (superadmin only)
// @Tags leads
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Lead ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /admin/leads/{id} [get]
func (r *leadsRoutes) GetLeadByID(c *gin.Context) {
	leadID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid lead ID format",
		})
		return
	}

	lead, err := r.leadsService.GetLeadByID(c.Request.Context(), leadID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "lead not found" {
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
		Data:    lead,
	})
}

// @Summary Update a lead
// @Description Update an existing lead's details (superadmin only)
// @Tags leads
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Lead ID"
// @Param input body entities.UpdateLeadInput true "Updated lead details"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 409 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /admin/leads/{id} [put]
func (r *leadsRoutes) UpdateLead(c *gin.Context) {
	leadID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid lead ID format",
		})
		return
	}

	var input entities.UpdateLeadInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	lead, err := r.leadsService.UpdateLead(c.Request.Context(), leadID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "lead not found" {
			status = http.StatusNotFound
		} else if err.Error() == "phone number already in use" || err.Error() == "email already in use" {
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
		Data:    lead,
	})
}

// @Summary Get lead by phone number
// @Description Get a lead by their phone number (superadmin only)
// @Tags leads
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param phone path string true "Phone Number"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /admin/leads/phone/{phone} [get]
func (r *leadsRoutes) GetLeadByPhoneNumber(c *gin.Context) {
	phoneNumber := c.Param("phone")
	if phoneNumber == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "phone number is required",
		})
		return
	}

	lead, err := r.leadsService.GetLeadByPhoneNumber(c.Request.Context(), phoneNumber)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "lead not found" {
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
		Data:    lead,
	})
}

// @Summary Get lead by email
// @Description Get a lead by their email (superadmin only)
// @Tags leads
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param email path string true "Email"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /admin/leads/email/{email} [get]
func (r *leadsRoutes) GetLeadByEmail(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "email is required",
		})
		return
	}

	lead, err := r.leadsService.GetLeadByEmail(c.Request.Context(), email)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "lead not found" {
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
		Data:    lead,
	})
}
