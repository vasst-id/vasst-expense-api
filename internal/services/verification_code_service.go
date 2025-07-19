package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/repositories"
	errorsutil "github.com/vasst-id/vasst-expense-api/internal/utils/errors"
)

//go:generate mockgen -source=verification_code_service.go -package=mock -destination=mock/verification_code_service_mock.go
type (
	VerificationCodeService interface {
		CreateVerificationCode(ctx context.Context, input *entities.CreateVerificationCodeRequest) (*entities.VerificationCode, error)
		VerifyCode(ctx context.Context, input *entities.VerifyVerificationCodeRequest) error
		ResendVerificationCode(ctx context.Context, phoneNumber, codeType string) error
		CleanupExpiredCodes(ctx context.Context) error
	}

	verificationCodeService struct {
		verificationCodeRepo repositories.VerificationCodeRepository
		userRepo             repositories.UserRepository
	}
)

// NewVerificationCodeService creates a new verification code service
func NewVerificationCodeService(verificationCodeRepo repositories.VerificationCodeRepository, userRepo repositories.UserRepository) VerificationCodeService {
	return &verificationCodeService{
		verificationCodeRepo: verificationCodeRepo,
		userRepo:             userRepo,
	}
}

// generateRandomCode generates a random 6-digit code
func generateRandomCode() (string, error) {
	code := ""
	for i := 0; i < 6; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code += fmt.Sprintf("%d", num.Int64())
	}
	return code, nil
}

// CreateVerificationCode creates a new verification code
func (s *verificationCodeService) CreateVerificationCode(ctx context.Context, input *entities.CreateVerificationCodeRequest) (*entities.VerificationCode, error) {
	// Validate required fields
	if input.PhoneNumber == "" {
		return nil, errors.New("phone number is required")
	}
	if input.CodeType == "" {
		return nil, errors.New("code type is required")
	}

	// Check if there's already an active code for this phone number and type
	existingCode, err := s.verificationCodeRepo.FindActiveByPhoneNumberAndType(ctx, input.PhoneNumber, input.CodeType)
	if err != nil {
		return nil, err
	}

	// If there's an existing active code, check if it's too recent (rate limiting)
	if existingCode != nil {
		timeSinceCreation := time.Since(existingCode.CreatedAt)
		if timeSinceCreation < 1*time.Minute { // 1 minute cooldown
			return nil, errorsutil.New(429, "please wait before requesting another code")
		}
	}

	// Generate a random 6-digit code
	code, err := generateRandomCode()
	if err != nil {
		return nil, err
	}

	// Set expiration time (10 minutes from now)
	expiresAt := time.Now().Add(10 * time.Minute)

	verificationCode := &entities.VerificationCode{
		VerificationCodeID: uuid.New(),
		PhoneNumber:        input.PhoneNumber,
		Code:               code,
		CodeType:           input.CodeType,
		ExpiresAt:          expiresAt,
		IsUsed:             false,
		AttemptsCount:      0,
		MaxAttempts:        3,
	}

	// Create the verification code
	createdCode, err := s.verificationCodeRepo.Create(ctx, verificationCode)
	if err != nil {
		return nil, err
	}

	// TODO: Send SMS with the code
	// For now, just return the code (in production, this should be sent via SMS)
	fmt.Printf("Verification code for %s: %s\n", input.PhoneNumber, code)

	return &createdCode, nil
}

// VerifyCode verifies a verification code
func (s *verificationCodeService) VerifyCode(ctx context.Context, input *entities.VerifyVerificationCodeRequest) error {
	// Validate required fields
	if input.PhoneNumber == "" {
		return errors.New("phone number is required")
	}
	if input.Code == "" {
		return errors.New("code is required")
	}

	// Find the active verification code
	verificationCode, err := s.verificationCodeRepo.FindActiveByPhoneNumberAndType(ctx, input.PhoneNumber, "phone_verification")
	if err != nil {
		return err
	}
	if verificationCode == nil {
		return errorsutil.New(404, "no active verification code found")
	}

	// Check if code has expired
	if time.Now().After(verificationCode.ExpiresAt) {
		return errorsutil.New(400, "verification code has expired")
	}

	// Check if maximum attempts exceeded
	if verificationCode.AttemptsCount >= verificationCode.MaxAttempts {
		return errorsutil.New(400, "maximum verification attempts exceeded")
	}

	// Increment attempts count
	err = s.verificationCodeRepo.IncrementAttempts(ctx, verificationCode.VerificationCodeID)
	if err != nil {
		return err
	}

	// Verify the code
	if verificationCode.Code != input.Code {
		return errorsutil.New(400, "invalid verification code")
	}

	// Mark the code as used
	err = s.verificationCodeRepo.MarkAsUsed(ctx, verificationCode.VerificationCodeID)
	if err != nil {
		return err
	}

	// Update user's phone verification status
	// user, err := s.userRepo.FindByPhoneNumber(ctx, input.PhoneNumber)
	// if err != nil {
	// 	return err
	// }
	// if user == nil {
	// 	return errorsutil.New(404, "user not found")
	// }

	// now := time.Now()
	// user.PhoneVerifiedAt = &now

	// _, err = s.userRepo.Update(ctx, user)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// ResendVerificationCode resends a verification code
func (s *verificationCodeService) ResendVerificationCode(ctx context.Context, phoneNumber, codeType string) error {
	// Check if there's already an active code
	existingCode, err := s.verificationCodeRepo.FindActiveByPhoneNumberAndType(ctx, phoneNumber, codeType)
	if err != nil {
		return err
	}

	// If there's an existing active code, check if it's too recent
	if existingCode != nil {
		timeSinceCreation := time.Since(existingCode.CreatedAt)
		if timeSinceCreation < 1*time.Minute { // 1 minute cooldown
			return errorsutil.New(429, "please wait before requesting another code")
		}
	}

	// Create a new verification code
	input := &entities.CreateVerificationCodeRequest{
		PhoneNumber: phoneNumber,
		CodeType:    codeType,
	}

	_, err = s.CreateVerificationCode(ctx, input)
	return err
}

// CleanupExpiredCodes removes expired verification codes
func (s *verificationCodeService) CleanupExpiredCodes(ctx context.Context) error {
	return s.verificationCodeRepo.CleanupExpiredCodes(ctx)
}
