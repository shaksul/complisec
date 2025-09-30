package domain

import (
	"context"
	"errors"
	"log"
	"time"

	"risknexus/backend/internal/repo"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo       *repo.UserRepo
	roleRepo       *repo.RoleRepo
	permissionRepo *repo.PermissionRepo
	jwtSecret      string
}

func NewAuthService(userRepo *repo.UserRepo, roleRepo *repo.RoleRepo, permissionRepo *repo.PermissionRepo) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		jwtSecret:      "your-secret-key", // TODO: get from config
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*repo.User, []string, error) {
	log.Printf("DEBUG: AuthService.Login email=%s", email)
	user, err := s.userRepo.GetByEmail(ctx, "00000000-0000-0000-0000-000000000001", email)
	if err != nil {
		log.Printf("ERROR: AuthService.Login GetByEmail failed: %v", err)
		return nil, nil, err
	}
	if user == nil {
		log.Printf("WARN: AuthService.Login user not found email=%s", email)
		return nil, nil, errors.New("invalid credentials")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		log.Printf("WARN: AuthService.Login invalid password email=%s", email)
		return nil, nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, nil, errors.New("account is disabled")
	}

	// Get user roles
	roles, err := s.userRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		log.Printf("ERROR: AuthService.Login GetUserRoles failed: %v", err)
		return nil, nil, err
	}

	log.Printf("DEBUG: AuthService.Login success user=%s roles=%v", user.ID, roles)
	return user, roles, nil
}

func (s *AuthService) GenerateTokens(userID, tenantID string, roles []string) (string, string, error) {
	// Access token (15 minutes)
	accessClaims := jwt.MapClaims{
		"user_id":   userID,
		"tenant_id": tenantID,
		"roles":     roles,
		"exp":       time.Now().Add(15 * time.Minute).Unix(),
		"type":      "access",
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", err
	}

	// Refresh token (7 days)
	refreshClaims := jwt.MapClaims{
		"user_id":   userID,
		"tenant_id": tenantID,
		"exp":       time.Now().Add(7 * 24 * time.Hour).Unix(),
		"type":      "refresh",
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

func (s *AuthService) GetUserFromToken(token *jwt.Token) (string, string, []string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", nil, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", "", nil, errors.New("invalid user_id in token")
	}

	tenantID, ok := claims["tenant_id"].(string)
	if !ok {
		return "", "", nil, errors.New("invalid tenant_id in token")
	}

	// For refresh tokens, roles are not included in the token
	// We'll need to get them from the database
	return userID, tenantID, nil, nil
}

func (s *AuthService) GetUserFromAccessToken(token *jwt.Token) (string, string, []string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", nil, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", "", nil, errors.New("invalid user_id in token")
	}

	tenantID, ok := claims["tenant_id"].(string)
	if !ok {
		return "", "", nil, errors.New("invalid tenant_id in token")
	}

	rolesInterface, ok := claims["roles"].([]interface{})
	if !ok {
		return "", "", nil, errors.New("invalid roles in token")
	}

	var roles []string
	for _, role := range rolesInterface {
		if roleStr, ok := role.(string); ok {
			roles = append(roles, roleStr)
		}
	}

	return userID, tenantID, roles, nil
}

func (s *AuthService) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	return s.userRepo.GetUserRoles(ctx, userID)
}

func (s *AuthService) CheckPermission(ctx context.Context, userID, permission string) (bool, error) {
	// Get user roles
	roles, err := s.userRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return false, err
	}

	// Check if any role has the permission
	for _, roleName := range roles {
		// Get role by name (simplified - in production, get by ID)
		// For demo, we'll check if user has admin role
		if roleName == "Admin" {
			return true, nil
		}
	}

	return false, nil
}
