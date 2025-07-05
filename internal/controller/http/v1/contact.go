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

type contactRoutes struct {
	contactService services.ContactService
	auth           *middleware.AuthMiddleware
}

func newContactRoutes(handler *gin.RouterGroup, contactService services.ContactService, auth *middleware.AuthMiddleware) {
	r := &contactRoutes{
		contactService: contactService,
		auth:           auth,
	}

	h := handler.Group("/contacts")
	h.Use(r.auth.AuthRequired())
	{
		authGroup := h.Group("")
		authGroup.GET("", r.ListContacts)
		authGroup.GET("/:id", r.GetContactByID)
		authGroup.GET("/phone/:phone", r.GetContactByPhoneNumber)
		authGroup.GET("/organization/:id", r.ListContactsByOrganizationID)
		authGroup.POST("", r.CreateContact)
		authGroup.PUT("/:id", r.UpdateContact)
		authGroup.DELETE("/:id", r.DeleteContact)
	}
}

// @Summary Get all contacts
// @Description Get a list of contacts with optional filtering
// @Tags contacts
// @Accept json
// @Produce json
// @Param limit query int false "Limit for pagination"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /contacts [get]
func (r *contactRoutes) ListContacts(c *gin.Context) {
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

	contacts, err := r.contactService.ListContacts(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    contacts,
	})
}

// @Summary Get contact by ID
// @Description Get a contact by their ID
// @Tags contacts
// @Accept json
// @Produce json
// @Param id path string true "Contact ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /contacts/{id} [get]
func (r *contactRoutes) GetContactByID(c *gin.Context) {
	contactID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid contact ID format",
		})
		return
	}

	contact, err := r.contactService.GetContactByID(c.Request.Context(), contactID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "contact not found" {
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
		Data:    contact,
	})
}

// @Summary Get contact by phone number
// @Description Get a contact by their phone number
// @Tags contacts
// @Accept json
// @Produce json
// @Param phone path string true "Phone Number"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /contacts/phone/{phone} [get]
func (r *contactRoutes) GetContactByPhoneNumber(c *gin.Context) {
	phoneNumber := c.Param("phone")
	if phoneNumber == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "phone number is required",
		})
		return
	}

	contact, err := r.contactService.GetContactByPhoneNumber(c.Request.Context(), phoneNumber)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "contact not found" {
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
		Data:    contact,
	})
}

// @Summary List contacts by organization ID
// @Description Get a list of contacts for a specific organization
// @Tags contacts
// @Accept json
// @Produce json
// @Param id path string true "Organization ID"
// @Param limit query int false "Limit for pagination"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /contacts/organization/{id} [get]
func (r *contactRoutes) ListContactsByOrganizationID(c *gin.Context) {
	organizationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
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

	contacts, err := r.contactService.ListContactsByOrganizationID(c.Request.Context(), organizationID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    contacts,
	})
}

// @Summary Create a new contact
// @Description Create a new contact with the provided details
// @Tags contacts
// @Accept json
// @Produce json
// @Param input body entities.CreateContactInput true "Contact details"
// @Success 201 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /contacts [post]
func (r *contactRoutes) CreateContact(c *gin.Context) {
	var input entities.CreateContactInput
	organizationID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, &entities.ApiResponse{
			Success: false,
			Error:   "organization ID not found in context",
		})
		return
	}

	orgID, ok := organizationID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}
	input.OrganizationID = orgID

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	contact, err := r.contactService.CreateContact(c.Request.Context(), &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "contact with this phone number already exists" {
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
		Data:    contact,
	})
}

// @Summary Update a contact
// @Description Update an existing contact's details
// @Tags contacts
// @Accept json
// @Produce json
// @Param id path string true "Contact ID"
// @Param input body entities.UpdateContactInput true "Updated contact details"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /contacts/{id} [put]
func (r *contactRoutes) UpdateContact(c *gin.Context) {
	contactID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid contact ID format",
		})
		return
	}

	var input entities.UpdateContactInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	contact, err := r.contactService.UpdateContact(c.Request.Context(), contactID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "contact not found" {
			status = http.StatusNotFound
		} else if err.Error() == "phone number already in use" {
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
		Data:    contact,
	})
}

// @Summary Delete a contact
// @Description Delete a contact by their ID
// @Tags contacts
// @Accept json
// @Produce json
// @Param id path string true "Contact ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /contacts/{id} [delete]
func (r *contactRoutes) DeleteContact(c *gin.Context) {
	contactID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid contact ID format",
		})
		return
	}

	err = r.contactService.DeleteContact(c.Request.Context(), contactID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "contact not found" {
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
		Data:    "Contact deleted successfully",
	})
}
