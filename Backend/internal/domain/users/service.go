package users

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/config"
	sharedErrors "github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/shared/errors"
	"github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/shared/id"
	"github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/shared/validate"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo   Repository
	config config.Config
}

func NewService(repo Repository, cfg config.Config) *Service {
	return &Service{
		repo:   repo,
		config: cfg,
	}
}

func (s *Service) Register(ctx context.Context, input CreateUserInput) (*User, error) {
	if !validate.Email(input.Email) {
		return nil, sharedErrors.ErrInvalidInput
	}
	if !validate.Password(input.Password) {
		return nil, sharedErrors.ErrInvalidInput
	}

	existing, _ := s.repo.FindByEmail(ctx, input.Email)
	if existing != nil {
		return nil, sharedErrors.ErrAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), s.config.BcryptCost)
	if err != nil {
		return nil, sharedErrors.ErrInternal
	}

	now := time.Now()
	user := &User{
		ID:        id.New(),
		Email:     input.Email,
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	user.Password = "" // Don't return password
	return user, nil
}

func (s *Service) Login(ctx context.Context, input LoginInput) (*TokenPair, error) {
	if !validate.Email(input.Email) {
		return nil, sharedErrors.ErrInvalidInput
	}

	user, err := s.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, sharedErrors.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, sharedErrors.ErrUnauthorized
	}

	accessToken, err := s.generateAccessToken(user.ID)
	if err != nil {
		return nil, sharedErrors.ErrInternal
	}

	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return nil, sharedErrors.ErrInternal
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) GetUser(ctx context.Context, userID string) (*User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, sharedErrors.ErrNotFound
	}
	user.Password = "" // Don't return password
	return user, nil
}

func (s *Service) generateAccessToken(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		Issuer:    s.config.JWTIssuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.AccessTokenTTL)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

func (s *Service) generateRefreshToken(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		Issuer:    s.config.JWTIssuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.RefreshTokenTTL)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

func (s *Service) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return "", sharedErrors.ErrUnauthorized
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		if claims.Subject == "" {
			return "", sharedErrors.ErrUnauthorized
		}
		return claims.Subject, nil
	}

	return "", sharedErrors.ErrUnauthorized
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	userID, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, sharedErrors.ErrUnauthorized
	}

	// Verify user still exists
	_, err = s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, sharedErrors.ErrUnauthorized
	}

	// Generate new token pair
	accessToken, err := s.generateAccessToken(userID)
	if err != nil {
		return nil, sharedErrors.ErrInternal
	}

	newRefreshToken, err := s.generateRefreshToken(userID)
	if err != nil {
		return nil, sharedErrors.ErrInternal
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

