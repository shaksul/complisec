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

func NewAuthService(userRepo *repo.UserRepo, roleRepo *repo.RoleRepo, permissionRepo *repo.PermissionRepo, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		permissionRepo: permissionRepo,
		jwtSecret:      jwtSecret,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password, tenantID string) (*repo.User, []string, error) {
	log.Printf("DEBUG: AuthService.Login email=%s tenantID=%s", email, tenantID)
	user, err := s.userRepo.GetByEmail(ctx, tenantID, email)
	if err != nil {
		log.Printf("ERROR: AuthService.Login GetByEmail failed: %v", err)
		return nil, nil, err
	}
	if user == nil {
		log.Printf("WARN: AuthService.Login user not found email=%s tenantID=%s", email, tenantID)
		return nil, nil, errors.New("invalid credentials")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		log.Printf("WARN: AuthService.Login invalid password email=%s tenantID=%s", email, tenantID)
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

func (s *AuthService) GenerateTokens(user *repo.User, roles []string) (string, string, error) {
	// Generate access token
	accessClaims := jwt.MapClaims{
		"user_id":   user.ID,
		"email":     user.Email,
		"tenant_id": user.TenantID,
		"roles":     roles,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
		"iat":       time.Now().Unix(),
		"type":      "access",
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshClaims := jwt.MapClaims{
		"user_id":   user.ID,
		"tenant_id": user.TenantID,
		"exp":       time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iat":       time.Now().Unix(),
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

func (s *AuthService) ExtractUserFromToken(token *jwt.Token) (string, string, []string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", nil, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", "", nil, errors.New("user_id not found in token")
	}

	tenantID, ok := claims["tenant_id"].(string)
	if !ok {
		return "", "", nil, errors.New("tenant_id not found in token")
	}

	var roles []string
	if rolesClaim, ok := claims["roles"].([]interface{}); ok {
		for _, role := range rolesClaim {
			if roleStr, ok := role.(string); ok {
				roles = append(roles, roleStr)
			}
		}
	}

	return userID, tenantID, roles, nil
}

func (s *AuthService) RefreshToken(refreshTokenString string) (string, string, error) {
	token, err := s.ValidateToken(refreshTokenString)
	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return "", "", errors.New("invalid token type")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", "", errors.New("user_id not found in token")
	}

	_, ok = claims["tenant_id"].(string)
	if !ok {
		return "", "", errors.New("tenant_id not found in token")
	}

	// Get user from database
	user, err := s.userRepo.GetByID(context.Background(), userID)
	if err != nil {
		return "", "", err
	}

	if user == nil {
		return "", "", errors.New("user not found")
	}

	if !user.IsActive {
		return "", "", errors.New("account is disabled")
	}

	// Get user roles
	roles, err := s.userRepo.GetUserRoles(context.Background(), user.ID)
	if err != nil {
		return "", "", err
	}

	// Generate new tokens
	return s.GenerateTokens(user, roles)
}

func (s *AuthService) GetUserFromToken(tokenString string) (*repo.User, []string, error) {
	token, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, nil, err
	}

	userID, tenantID, roles, err := s.ExtractUserFromToken(token)
	if err != nil {
		return nil, nil, err
	}

	// Get user from database
	user, err := s.userRepo.GetByID(context.Background(), userID)
	if err != nil {
		return nil, nil, err
	}

	if user == nil {
		return nil, nil, errors.New("user not found")
	}

	if user.TenantID != tenantID {
		return nil, nil, errors.New("tenant mismatch")
	}

	return user, roles, nil
}

func (s *AuthService) GetUserFromAccessToken(tokenString string) (*repo.User, []string, error) {
	return s.GetUserFromToken(tokenString)
}

func (s *AuthService) GetUser(ctx context.Context, userID string) (*repo.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *AuthService) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	return s.userRepo.GetUserRoles(ctx, userID)
}

func (s *AuthService) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	return s.userRepo.GetUserPermissions(ctx, userID)
}
