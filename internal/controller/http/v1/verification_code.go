package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/middleware"
	"github.com/vasst-id/vasst-expense-api/internal/services"
)

type verificationCodeRoutes struct {
	verificationCodeService services.VerificationCodeService
	auth                    *middleware.AuthMiddleware
}

func newVerificationCodeRoutes(handler *gin.RouterGroup, verificationCodeService services.VerificationCodeService, auth *middleware.AuthMiddleware) {
	r := &verificationCodeRoutes{
		verificationCodeService: verificationCodeService,
		auth:                    auth,
	}

	// Verification code endpoints
	verificationCodes := handler.Group("/verification-codes")
	{
		verificationCodes.POST("/create", r.CreateVerificationCode)
		verificationCodes.POST("/verify", r.VerifyCode)
		verificationCodes.POST("/resend", r.ResendVerificationCode)
	}
}

// @Summary Create verification code
// @Description Create a new verification code for phone number
// @Tags verification-codes
// @Accept json
// @Produce json
// @Param input body entities.CreateVerificationCodeRequest true "Verification code request"
// @Success 201 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 429 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /verification-codes/create [post]
func (r *verificationCodeRoutes) CreateVerificationCode(c *gin.Context) {
	var input entities.CreateVerificationCodeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	verificationCode, err := r.verificationCodeService.CreateVerificationCode(c.Request.Context(), &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "phone number is required" ||
			err.Error() == "code type is required" {
			status = http.StatusBadRequest
		} else if err.Error() == "please wait before requesting another code" {
			status = http.StatusTooManyRequests
		}
		c.JSON(status, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, &entities.ApiResponse{
		Success: true,
		Data:    verificationCode,
		Message: "Verification code sent successfully",
	})
}

// @Summary Verify code
// @Description Verify a verification code for phone number
// @Tags verification-codes
// @Accept json
// @Produce json
// @Param input body entities.VerifyVerificationCodeRequest true "Verification request"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /verification-codes/verify [post]
func (r *verificationCodeRoutes) VerifyCode(c *gin.Context) {
	var input entities.VerifyVerificationCodeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	err := r.verificationCodeService.VerifyCode(c.Request.Context(), &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "phone number is required" ||
			err.Error() == "code is required" ||
			err.Error() == "invalid verification code" ||
			err.Error() == "verification code has expired" ||
			err.Error() == "maximum verification attempts exceeded" {
			status = http.StatusBadRequest
		} else if err.Error() == "no active verification code found" ||
			err.Error() == "user not found" {
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

// @Summary Resend verification code
// @Description Resend a verification code for phone number
// @Tags verification-codes
// @Accept json
// @Produce json
// @Param phone_number query string true "Phone number"
// @Param code_type query string true "Code type"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 429 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /verification-codes/resend [post]
func (r *verificationCodeRoutes) ResendVerificationCode(c *gin.Context) {
	phoneNumber := c.Query("phone_number")
	codeType := c.Query("code_type")

	if phoneNumber == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "phone number is required",
		})
		return
	}

	if codeType == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "code type is required",
		})
		return
	}

	err := r.verificationCodeService.ResendVerificationCode(c.Request.Context(), phoneNumber, codeType)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "please wait before requesting another code" {
			status = http.StatusTooManyRequests
		}
		c.JSON(status, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Message: "Verification code resent successfully",
	})
}
