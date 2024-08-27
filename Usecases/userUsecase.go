package Usecases

import (
	"Loan_manager/Domain"
	"Loan_manager/Repository"
	"Loan_manager/infrastructure"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserUsecase defines the contract for user-related use cases
type UserUsecase interface {
	Register(input Domain.RegisterInput) (*Domain.User, error)
	UpdateUser(username string, updatedUser *Domain.UpdateUserInput) error
	DeleteUser(username string) error
	Login(c *gin.Context, LoginUser *Domain.LoginInput) (string, error)
	Logout(tokenString string) error
	ForgotPassword(username string) (string, error)
	Reset(token string) (string, error)
	UpdatePassword(username string, newPassword string) error
	Verify(token string) error
}

// userUsecase implements the UserUsecase interface
type userUsecase struct {
	userRepo        Repository.UserRepository
	emailService    *infrastructure.EmailService
	passwordService *infrastructure.PasswordService
}

// NewUserUsecase creates a new instance of UserUsecase
func NewUserUsecase(userRepo Repository.UserRepository, emailService *infrastructure.EmailService) UserUsecase {
	return &userUsecase{
		userRepo:        userRepo,
		emailService:    emailService,
		passwordService: infrastructure.NewPasswordService(),
	}
}

const (
	passwordMinLength = 8
	passwordMaxLength = 20
)

// Register handles user registration logic
func (u *userUsecase) Register(input Domain.RegisterInput) (*Domain.User, error) {
	// Validate username
	if strings.Contains(input.Username, "@") {
		return nil, errors.New("username must not contain '@'")
	}

	// Check if username already exists
	if _, err := u.userRepo.FindByUsername(input.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	// Validate email format
	if !isValidEmail(input.Email) {
		return nil, errors.New("invalid email format")
	}

	// Check if email already registered
	if _, err := u.userRepo.FindByEmail(input.Email); err == nil {
		return nil, errors.New("email already registered")
	}

	// Validate password strength
	if err := validatePasswordStrength(input.Password); err != nil {
		return nil, err
	}

	// Hash the password
	hashedPassword, err := u.passwordService.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Create new user
	user := &Domain.User{
		Id:             primitive.NewObjectID(),
		Name:           input.Name,
		Username:       input.Username,
		Email:          input.Email,
		Password:       hashedPassword,
		ProfilePicture: input.ProfilePicture,
		Bio:            input.Bio,
		Gender:         input.Gender,
		Address:        input.Address,
		IsActive:       false, // Initially inactive
		PostsIDs:       []string{},
	}

	// Set user role based on database state
	if ok, err := u.userRepo.IsDbEmpty(); ok && err == nil {
		user.Role = "admin"
	} else {
		user.Role = "user"
	}

	// Save user to repository
	if err := u.userRepo.Save(user); err != nil {
		return nil, fmt.Errorf("failed to save user: %v", err)
	}

	// Generate a verification token
	newToken, err := infrastructure.GenerateResetToken(user.Username, []byte("BlogManagerSecretKey"))
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification token: %v", err)
	}

	// Construct the email body
	subject := "Welcome to Our Service!"
	body := fmt.Sprintf("Hi %s,\n\nWelcome to our platform! Please verify your account by clicking the link below:\n\nhttp://localhost:8080/verify/%s\n\nThank you!", input.Name, newToken)

	// Send verification email
	if err := u.emailService.SendEmail(input.Email, subject, body); err != nil {
		return nil, fmt.Errorf("failed to send welcome email: %v", err)
	}

	return user, nil
}

// UpdateUser handles the user update logic
func (u *userUsecase) UpdateUser(username string, updatedUser *Domain.UpdateUserInput) error {
	_, err := u.userRepo.FindByUsername(username)
	if err != nil {
		return errors.New("user not found")
	}

	updateFields := bson.M{}

	if updatedUser.Username != "" {
		if strings.Contains(updatedUser.Username, "@") {
			return errors.New("username must not contain '@'")
		}
		updateFields["username"] = updatedUser.Username
	}
	if updatedUser.Password != "" {
		hashedPassword, err := u.passwordService.HashPassword(updatedUser.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %v", err)
		}
		updateFields["password"] = hashedPassword
	}
	if updatedUser.ProfilePicture != "" {
		updateFields["profile_picture"] = updatedUser.ProfilePicture
	}
	if updatedUser.Bio != "" {
		updateFields["bio"] = updatedUser.Bio
	}
	if updatedUser.Address != "" {
		updateFields["address"] = updatedUser.Address
	}

	if err := u.userRepo.Update(username, updateFields); err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}

// DeleteUser handles the user deletion logic
func (u *userUsecase) DeleteUser(username string) error {
	_, err := u.userRepo.FindByUsername(username)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	if err := u.userRepo.Delete(username); err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}

// Login handles the user login logic
func (u *userUsecase) Login(c *gin.Context, loginUser *Domain.LoginInput) (string, error) {
	user, err := u.userRepo.FindByUsername(loginUser.Username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	if err := u.passwordService.ComparePasswords(user.Password, loginUser.Password); err != nil {
		return "", errors.New("invalid username or password")
	}

	accessToken, err := infrastructure.GenerateJWT(user.Username, user.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, err := infrastructure.GenerateRefreshToken(user.Username)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %v", err)
	}

	c.SetCookie("refresh_token", refreshToken, 3600, "/", "", false, true)

	if err := u.userRepo.InsertToken(user.Username, accessToken, refreshToken); err != nil {
		return "", fmt.Errorf("failed to store tokens: %v", err)
	}

	if !user.IsActive {
		return "", fmt.Errorf("user not verified")
	}

	return accessToken, nil
}

// Logout handles the user logout logic
func (u *userUsecase) Logout(tokenString string) error {
	if err := u.userRepo.ExpireToken(tokenString); err != nil {
		return err
	}
	return nil
}

// ForgotPassword handles the forgot password logic
func (u *userUsecase) ForgotPassword(username string) (string, error) {
	user, err := u.userRepo.FindByUsername(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	resetToken, err := infrastructure.GenerateResetToken(user.Username, []byte("BlogManagerSecretKey"))
	if err != nil {
		return "", err
	}
	subject := "Password Reset Request"
	body := fmt.Sprintf(`
	Hi %s,

	It seems like you requested a password reset. No worries, it happens to the best of us! You can reset your password by clicking the link below:

	<a href="http://localhost:8080/reset/%s">Reset Your Password</a>

	If you did not request a password reset, please ignore this email.

	Best regards,
	Your Support Team
	`, user.Name, resetToken)

	if err := u.emailService.SendEmail(user.Email, subject, body); err != nil {
		return "", fmt.Errorf("failed to send reset email: %v", err)
	}

	return resetToken, nil
}

// Reset handles the password reset logic
func (u *userUsecase) Reset(token string) (string, error) {
	claims, err := infrastructure.ParseResetToken(token, []byte("BlogManagerSecretKey"))
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return "", err
	}

	user, err := u.userRepo.FindByUsername(claims.Username)
	if err != nil {
		return "", errors.New("user not found")
	}
	accessToken, err := infrastructure.GenerateJWT(user.Username, user.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %v", err)
	}

	return accessToken, nil
}

// UpdatePassword handles the update password logic
func (u *userUsecase) UpdatePassword(username string, newPassword string) error {
	hashedPassword, err := u.passwordService.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	if err := u.userRepo.Update(username, bson.M{"password": hashedPassword}); err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}

	return nil
}

// Verify handles the verification of users via token
func (u *userUsecase) Verify(token string) error {
	claims, err := infrastructure.ParseResetToken(token, []byte("BlogManagerSecretKey"))
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return err
	}

	user, err := u.userRepo.FindByUsername(claims.Username)
	if err != nil {
		return errors.New("user not found")
	}

	user.IsActive = true
	updateFields := bson.M{"is_active": true}
	if err := u.userRepo.Update(user.Username, updateFields); err != nil {
		return fmt.Errorf("failed to update user status: %v", err)
	}

	return nil
}

// validatePasswordStrength ensures the password meets certain criteria
func validatePasswordStrength(password string) error {
	if len(password) < passwordMinLength || len(password) > passwordMaxLength {
		return fmt.Errorf("password must be between %d and %d characters", passwordMinLength, passwordMaxLength)
	}
	return nil
}

// isValidEmail validates the email format
func isValidEmail(email string) bool {
	// Simplified email validation
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
