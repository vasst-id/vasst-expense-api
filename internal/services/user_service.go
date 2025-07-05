package services

import (
	"context"
	"errors"
	"math/rand"

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
		CreateUser(ctx context.Context, input *entities.CreateUserInput) (*entities.User, error)
		UpdateUser(ctx context.Context, userID uuid.UUID, input *entities.UpdateUserInput) (*entities.User, error)
		DeleteUser(ctx context.Context, userID uuid.UUID) error
		ListAllUsers(ctx context.Context, limit, offset int) ([]*entities.User, error)
		GetUserByID(ctx context.Context, userID uuid.UUID) (*entities.User, error)
		GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.User, error)
		GetUserByUsername(ctx context.Context, username string) (*entities.User, error)
		ResetUserPassword(ctx context.Context, input *entities.ResetPasswordInput) error
		GenerateUserPassword(ctx context.Context, userID uuid.UUID, input *entities.GenerateUserPasswordInput) (string, error)

		// Organization-scoped methods
		ListUsersByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entities.User, error)
		GetUserByIDAndOrganization(ctx context.Context, userID, organizationID uuid.UUID) (*entities.User, error)
		GetUserByPhoneNumberAndOrganization(ctx context.Context, phoneNumber string, organizationID uuid.UUID) (*entities.User, error)
		GetUserByUsernameAndOrganization(ctx context.Context, username string, organizationID uuid.UUID) (*entities.User, error)
		Login(ctx context.Context, input *entities.LoginInput) (*entities.LoginResponse, error)
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

// Global (superadmin) methods
func (s *userService) CreateUser(ctx context.Context, input *entities.CreateUserInput) (*entities.User, error) {
	if input.PhoneNumber == "" {
		return nil, errors.New("phone number is required")
	}
	if input.UserFullName == "" {
		return nil, errors.New("full name is required")
	}
	if input.Username == "" {
		return nil, errors.New("username is required")
	}
	if input.Password == "" {
		return nil, errors.New("password is required")
	}
	if input.OrganizationID == uuid.Nil {
		return nil, errors.New("organization ID is required")
	}
	if input.RoleID == 0 {
		return nil, errors.New("role ID is required")
	}

	existingUser, err := s.userRepo.FindByPhoneNumber(ctx, input.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errorsutil.New(409, "user with this phone number already exists")
	}

	existingUsername, err := s.userRepo.FindByUsername(ctx, input.Username)
	if err != nil {
		return nil, err
	}
	if existingUsername != nil {
		return nil, errorsutil.New(409, "user with this username already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		UserID:         uuid.New(),
		OrganizationID: input.OrganizationID,
		RoleID:         input.RoleID,
		UserFullName:   input.UserFullName,
		PhoneNumber:    input.PhoneNumber,
		Username:       input.Username,
		Password:       string(hashedPassword),
		IsActive:       input.IsActive,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, userID uuid.UUID, input *entities.UpdateUserInput) (*entities.User, error) {
	existingUser, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, errorsutil.New(404, "user not found")
	}

	if input.PhoneNumber != "" && input.PhoneNumber != existingUser.PhoneNumber {
		userWithPhone, err := s.userRepo.FindByPhoneNumber(ctx, input.PhoneNumber)
		if err != nil {
			return nil, err
		}
		if userWithPhone != nil && userWithPhone.UserID != userID {
			return nil, errorsutil.New(409, "phone number already in use")
		}
	}
	if input.Username != "" && input.Username != existingUser.Username {
		userWithUsername, err := s.userRepo.FindByUsername(ctx, input.Username)
		if err != nil {
			return nil, err
		}
		if userWithUsername != nil && userWithUsername.UserID != userID {
			return nil, errorsutil.New(409, "username already in use")
		}
	}

	if input.UserFullName != "" {
		existingUser.UserFullName = input.UserFullName
	}
	if input.PhoneNumber != "" {
		existingUser.PhoneNumber = input.PhoneNumber
	}
	if input.Username != "" {
		existingUser.Username = input.Username
	}
	if input.Password != "" {
		existingUser.Password = input.Password
	}
	if input.RoleID != 0 {
		existingUser.RoleID = input.RoleID
	}
	if input.IsActive != nil {
		existingUser.IsActive = *input.IsActive
	}

	if err := s.userRepo.Update(ctx, existingUser); err != nil {
		return nil, err
	}

	return existingUser, nil
}

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

func (s *userService) ListAllUsers(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	return s.userRepo.ListAll(ctx, limit, offset)
}

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

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorsutil.New(404, "user not found")
	}
	return user, nil
}

// Organization-scoped methods
func (s *userService) ListUsersByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entities.User, error) {
	return s.userRepo.ListByOrganization(ctx, organizationID, limit, offset)
}

func (s *userService) GetUserByIDAndOrganization(ctx context.Context, userID, organizationID uuid.UUID) (*entities.User, error) {
	user, err := s.userRepo.FindByIDAndOrganization(ctx, userID, organizationID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorsutil.New(404, "user not found")
	}
	return user, nil
}

func (s *userService) GetUserByPhoneNumberAndOrganization(ctx context.Context, phoneNumber string, organizationID uuid.UUID) (*entities.User, error) {
	user, err := s.userRepo.FindByPhoneNumberAndOrganization(ctx, phoneNumber, organizationID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorsutil.New(404, "user not found")
	}
	return user, nil
}

func (s *userService) GetUserByUsernameAndOrganization(ctx context.Context, username string, organizationID uuid.UUID) (*entities.User, error) {
	user, err := s.userRepo.FindByUsernameAndOrganization(ctx, username, organizationID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorsutil.New(404, "user not found")
	}
	return user, nil
}

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

	// Get the user
	user, err := s.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return errorsutil.New(404, "user not found")
	}

	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword))
	if err != nil {
		return errorsutil.New(400, "old password is incorrect")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update the user's password
	user.Password = string(hashedPassword)

	// Save the updated user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *userService) GenerateUserPassword(ctx context.Context, userID uuid.UUID, input *entities.GenerateUserPasswordInput) (string, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errorsutil.New(404, "user not found")
	}

	if input.Password == "" {
		input.Password = GenerateRandomPassword(12)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.Password = string(hashedPassword)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return "", err
	}

	return "Success generate the password", nil
}

func GenerateRandomPassword(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (s *userService) Login(ctx context.Context, input *entities.LoginInput) (*entities.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(ctx, input.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorsutil.New(404, "user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return nil, errorsutil.New(400, "password is incorrect")
	}

	// Generate JWT token
	token, err := s.authMiddleware.GenerateToken(user.UserID, user.OrganizationID, user.Username, int64(user.RoleID))
	if err != nil {
		return nil, err
	}

	loginResponse := &entities.LoginResponse{
		AccessToken:    token,
		UserID:         user.UserID,
		OrganizationID: user.OrganizationID,
		Username:       user.Username,
		RoleID:         user.RoleID,
	}

	return loginResponse, nil
}
