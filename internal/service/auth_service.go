package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/veezify/interview-ai-backend/internal/api/response"
	"github.com/veezify/interview-ai-backend/internal/domain/model"
	"github.com/veezify/interview-ai-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepository *repository.UserRepository
	jwtSecret      []byte
	tokenDuration  time.Duration
}

func NewAuthService(userRepository *repository.UserRepository) *AuthService {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-default-secret-key-for-development-only"
	}

	return &AuthService{
		userRepository: userRepository,
		jwtSecret:      []byte(jwtSecret),
		tokenDuration:  time.Hour * 24,
	}
}

type Claims struct {
	UserID        string `json:"user_id"`
	Email         string `json:"email"`
	AssociationId string `json:association_id`
	jwt.StandardClaims
}

func (s *AuthService) Login(email, password string) (*response.LoginResponse, error) {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !user.Active {
		return nil, errors.New("user account is inactive")
	}

	expiresAt := time.Now().Add(s.tokenDuration)
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	userWithDetails, err := s.userRepository.FindUserWithAssociationsAndRoles(user.ID)
	if err != nil {
		return nil, err
	}

	return &response.LoginResponse{
		Token:     signedToken,
		ExpiresAt: expiresAt,
		User:      userWithDetails,
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}

	return "", errors.New("invalid token")
}

func (s *AuthService) GetUserFromToken(tokenString string) (*model.User, error) {
	userID, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	return s.userRepository.FindUserWithAssociationsAndRoles(userID)
}
