package services

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/middleware"
	"github.com/vasst-id/vasst-expense-api/internal/repositories"
	errorsutil "github.com/vasst-id/vasst-expense-api/internal/utils/errors"
)

//go:generate mockgen -source=user_service.go -package=mock -destination=mock/user_service_mock.go
type (
	UserService interface {
		// Global (superadmin) methods
		CreateUser(ctx context.Context, input *entities.CreateUserInput) (*entities.LoginResponse, error)
		UpdateUser(ctx context.Context, userID uuid.UUID, input *entities.UpdateUserInput) (*entities.User, error)
		DeleteUser(ctx context.Context, userID uuid.UUID) error
		ListAllUsers(ctx context.Context, limit, offset int) ([]*entities.User, error)
		GetUserByID(ctx context.Context, userID uuid.UUID) (*entities.User, error)
		GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
		GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.User, error)
		ResetUserPassword(ctx context.Context, input *entities.ResetPasswordInput) error

		// Authentication methods
		Login(ctx context.Context, input *entities.LoginRequest) (*entities.LoginResponse, error)
		ForgotPassword(ctx context.Context, input *entities.ForgotPasswordRequest) error
		ResetPassword(ctx context.Context, input *entities.ResetPasswordRequest) error
		ChangePassword(ctx context.Context, input *entities.ChangePasswordRequest) error
		VerifyPhone(ctx context.Context, input *entities.VerifyPhoneRequest) error
		ResendVerificationCode(ctx context.Context, input *entities.ResendVerificationCodeRequest) error
		VerifyEmail(ctx context.Context, input *entities.VerifyEmailRequest) error
		ResendVerificationEmail(ctx context.Context, input *entities.ResendVerificationEmailRequest) error
	}

	userService struct {
		userRepo       repositories.UserRepository
		authMiddleware *middleware.AuthMiddleware
	}
)

// NewUserService creates a new user service
func NewUserService(userRepo repositories.UserRepository, authMiddleware *middleware.AuthMiddleware) UserService {
	return &userService{
		userRepo:       userRepo,
		authMiddleware: authMiddleware,
	}
}

// CreateUser creates a new user
func (s *userService) CreateUser(ctx context.Context, input *entities.CreateUserInput) (*entities.LoginResponse, error) {
	// Validate required fields
	// if input.Email == "" {
	// 	return nil, errors.New("email is required")
	// }
	if input.PhoneNumber == "" {
		return nil, errors.New("phone number is required")
	}
	if input.FirstName == "" {
		return nil, errors.New("first name is required")
	}
	if input.LastName == "" {
		return nil, errors.New("last name is required")
	}
	if input.Password == "" {
		return nil, errors.New("password is required")
	}
	if len(input.Password) != 6 {
		return nil, errors.New("password must be exactly 6 digits")
	}

	// Set default currency ID
	currencyID := 1

	// Set default subscription plan ID
	subscriptionPlanID := 1

	// Check if user already exists
	// existingUser, err := s.userRepo.FindByEmail(ctx, input.Email)
	// if err != nil {
	// 	return nil, err
	// }
	// if existingUser != nil {
	// 	return nil, errorsutil.New(409, "user with this email already exists")
	// }

	existingPhone, err := s.userRepo.FindByPhoneNumber(ctx, input.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if existingPhone != nil {
		return nil, errorsutil.New(409, "user with this phone number already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Set default timezone if not provided
	timezone := input.Timezone
	if timezone == "" {
		timezone = "Asia/Jakarta"
	}

	user := &entities.User{
		UserID:             uuid.New(),
		Email:              "", // This will be NULL in database since it's empty
		PhoneNumber:        input.PhoneNumber,
		PasswordHash:       string(hashedPassword),
		FirstName:          input.FirstName,
		LastName:           input.LastName,
		Timezone:           timezone,
		CurrencyID:         currencyID,
		SubscriptionPlanID: subscriptionPlanID,
		Status:             entities.UserStatusActive,
	}

	// Create the user - the repository will populate the struct with the actual data from DB
	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := s.authMiddleware.GenerateToken(&createdUser)
	if err != nil {
		return nil, err
	}

	loginResponse := &entities.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   3600, // 1 hour
		User:        &createdUser,
	}

	return loginResponse, nil
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(ctx context.Context, userID uuid.UUID, input *entities.UpdateUserInput) (*entities.User, error) {
	existingUser, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, errorsutil.New(404, "user not found")
	}

	// Check for email uniqueness
	if input.Email != "" && input.Email != existingUser.Email {
		userWithEmail, err := s.userRepo.FindByEmail(ctx, input.Email)
		if err != nil {
			return nil, err
		}
		if userWithEmail != nil && userWithEmail.UserID != userID {
			return nil, errorsutil.New(409, "email already in use")
		}
	}

	// Check for phone number uniqueness
	if input.PhoneNumber != "" && input.PhoneNumber != existingUser.PhoneNumber {
		userWithPhone, err := s.userRepo.FindByPhoneNumber(ctx, input.PhoneNumber)
		if err != nil {
			return nil, err
		}
		if userWithPhone != nil && userWithPhone.UserID != userID {
			return nil, errorsutil.New(409, "phone number already in use")
		}
	}

	// Update fields
	if input.Email != "" {
		existingUser.Email = input.Email
	}
	if input.PhoneNumber != "" {
		existingUser.PhoneNumber = input.PhoneNumber
	}
	if input.FirstName != "" {
		existingUser.FirstName = input.FirstName
	}
	if input.LastName != "" {
		existingUser.LastName = input.LastName
	}
	if input.Timezone != "" {
		existingUser.Timezone = input.Timezone
	}
	if input.CurrencyID != 0 {
		existingUser.CurrencyID = input.CurrencyID
	}
	if input.SubscriptionPlanID != 0 {
		existingUser.SubscriptionPlanID = input.SubscriptionPlanID
	}
	if input.Status != 0 {
		existingUser.Status = input.Status
	}

	// Update the user - the repository will populate the struct with the actual data from DB
	updatedUser, err := s.userRepo.Update(ctx, existingUser)
	if err != nil {
		return nil, err
	}

	// Return the user with data populated from the database
	return &updatedUser, nil
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	existingUser, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errorsutil.New(404, "user not found")
	}
	return s.userRepo.Delete(ctx, userID)
}

// ListAllUsers returns all users with pagination
func (s *userService) ListAllUsers(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	return s.userRepo.ListAll(ctx, limit, offset)
}

// GetUserByID returns a user by ID
func (s *userService) GetUserByID(ctx context.Context, userID uuid.UUID) (*entities.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorsutil.New(404, "user not found")
	}
	return user, nil
}

// GetUserByEmail returns a user by email
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorsutil.New(404, "user not found")
	}
	return user, nil
}

// GetUserByPhoneNumber returns a user by phone number
func (s *userService) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.User, error) {
	user, err := s.userRepo.FindByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorsutil.New(404, "user not found")
	}
	return user, nil
}

// ResetUserPassword resets a user's password
func (s *userService) ResetUserPassword(ctx context.Context, input *entities.ResetPasswordInput) error {
	// Validate input
	if input.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}
	if input.OldPassword == "" {
		return errors.New("old password is required")
	}
	if input.NewPassword == "" {
		return errors.New("new password is required")
	}
	if len(input.NewPassword) != 6 {
		return errors.New("new password must be exactly 6 digits")
	}

	// Get the user
	user, err := s.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return errorsutil.New(404, "user not found")
	}

	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.OldPassword))
	if err != nil {
		return errorsutil.New(400, "old password is incorrect")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update the user's password
	user.PasswordHash = string(hashedPassword)

	// Save the updated user
	_, err = s.userRepo.Update(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// Login authenticates a user and returns a login response
func (s *userService) Login(ctx context.Context, input *entities.LoginRequest) (*entities.LoginResponse, error) {
	var user *entities.User
	var err error

	// Find user by email or phone
	if input.PhoneNumber != "" {
		user, err = s.userRepo.FindByPhoneNumber(ctx, input.PhoneNumber)
	} else {
		return nil, errorsutil.New(400, "phone number is required")
	}

	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorsutil.New(404, "user not found")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return nil, errorsutil.New(400, "password is incorrect")
	}

	// Check if user is active
	if user.Status != entities.UserStatusActive {
		return nil, errorsutil.New(403, "user account is inactive")
	}

	// Generate JWT token
	token, err := s.authMiddleware.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	loginResponse := &entities.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   3600, // 1 hour
		User:        user,
	}

	return loginResponse, nil
}

// ForgotPassword initiates password reset process
func (s *userService) ForgotPassword(ctx context.Context, input *entities.ForgotPasswordRequest) error {
	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return err
	}
	if user == nil {
		return errorsutil.New(404, "user not found")
	}

	// TODO: Implement password reset token generation and email sending
	// For now, just return success
	return nil
}

// ResetPassword resets password using a token
func (s *userService) ResetPassword(ctx context.Context, input *entities.ResetPasswordRequest) error {
	// TODO: Implement token validation and password reset
	// For now, just return success
	return nil
}

// ChangePassword changes user's password
func (s *userService) ChangePassword(ctx context.Context, input *entities.ChangePasswordRequest) error {
	// TODO: Get current user from context
	// For now, return not implemented error
	return errorsutil.New(501, "change password not implemented")
}

// VerifyPhone verifies a phone number with a code
func (s *userService) VerifyPhone(ctx context.Context, input *entities.VerifyPhoneRequest) error {
	user, err := s.userRepo.FindByPhoneNumber(ctx, input.PhoneNumber)
	if err != nil {
		return err
	}
	if user == nil {
		return errorsutil.New(404, "user not found")
	}

	// TODO: Implement phone verification logic
	// For now, just mark as verified
	now := time.Now()
	user.PhoneVerifiedAt = &now

	_, err = s.userRepo.Update(ctx, user)
	return err
}

// ResendVerificationCode resends verification code for phone
func (s *userService) ResendVerificationCode(ctx context.Context, input *entities.ResendVerificationCodeRequest) error {
	user, err := s.userRepo.FindByPhoneNumber(ctx, input.PhoneNumber)
	if err != nil {
		return err
	}
	if user == nil {
		return errorsutil.New(404, "user not found")
	}

	// TODO: Implement SMS sending logic
	// For now, just return success
	return nil
}

// VerifyEmail verifies an email with a token
func (s *userService) VerifyEmail(ctx context.Context, input *entities.VerifyEmailRequest) error {
	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return err
	}
	if user == nil {
		return errorsutil.New(404, "user not found")
	}

	// TODO: Implement email verification logic
	// For now, just mark as verified
	now := time.Now()
	user.EmailVerifiedAt = &now

	_, err = s.userRepo.Update(ctx, user)
	return err
}

// ResendVerificationEmail resends verification email
func (s *userService) ResendVerificationEmail(ctx context.Context, input *entities.ResendVerificationEmailRequest) error {
	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return err
	}
	if user == nil {
		return errorsutil.New(404, "user not found")
	}

	// TODO: Implement email sending logic
	// For now, just return success
	return nil
}
