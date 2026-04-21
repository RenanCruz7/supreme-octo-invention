package services

import (
	"fmt"
	"time"

	"awesomeProject/config"
	"awesomeProject/models"
	"awesomeProject/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (s *AuthService) Register(req models.RegisterRequest) (*models.AuthResponse, error) {
	existingUser, _ := repositories.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("email já cadastrado")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("erro ao processar senha: %v", err)
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	id, err := repositories.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar usuário: %v", err)
	}

	user.ID = id

	token, err := s.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar token: %v", err)
	}

	return &models.AuthResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}, nil
}

func (s *AuthService) Login(req models.LoginRequest) (*models.AuthResponse, error) {
	user, err := repositories.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		return nil, fmt.Errorf("email ou senha inválidos")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("email ou senha inválidos")
	}

	token, err := s.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar token: %v", err)
	}

	return &models.AuthResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}, nil
}

func (s *AuthService) GenerateToken(userID uint, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inválido")
		}
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("erro ao validar token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token inválido")
	}

	return claims, nil
}
