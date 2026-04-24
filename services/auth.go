package services

import (
	"time"

	"awesomeProject/config"
	"awesomeProject/errors"
	"awesomeProject/mappers"
	"awesomeProject/models"
	"awesomeProject/repositories"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
		return nil, errors.ErrConflict("email já cadastrado")
	}

	user, err := mappers.RegisterRequestToUser(req)
	if err != nil {
		return nil, errors.ErrBadRequestWithErr("erro ao processar senha", err)
	}

	id, err := repositories.CreateUser(*user)
	if err != nil {
		return nil, errors.ErrInternalWithErr("erro ao criar usuário", err)
	}

	user.ID = id

	token, err := s.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.ErrInternalWithErr("erro ao gerar token", err)
	}

	refreshToken, err := s.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errors.ErrInternalWithErr("erro ao gerar refresh token", err)
	}

	return &models.AuthResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Login(req models.LoginRequest) (*models.AuthResponse, error) {
	user, err := repositories.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		return nil, errors.ErrUnauthorized("email ou senha inválidos")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.ErrUnauthorized("email ou senha inválidos")
	}

	token, err := s.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.ErrInternalWithErr("erro ao gerar token", err)
	}

	refreshToken, err := s.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errors.ErrInternalWithErr("erro ao gerar refresh token", err)
	}

	return &models.AuthResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

// GenerateToken gera um access token JWT com expiração de 15 minutos
func (s *AuthService) GenerateToken(userID uint, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", errors.ErrInternalWithErr("erro ao gerar token", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken cria um refresh token UUID, persiste no banco e o retorna
func (s *AuthService) GenerateRefreshToken(userID uint) (string, error) {
	tokenString := uuid.New().String()

	rt := models.RefreshToken{
		UserID:    userID,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
	}

	if err := repositories.CreateRefreshToken(rt); err != nil {
		return "", errors.ErrInternalWithErr("erro ao salvar refresh token", err)
	}

	return tokenString, nil
}

// RefreshTokens valida o refresh token, revoga o antigo (rotação) e emite novo par
func (s *AuthService) RefreshTokens(refreshTokenStr string) (*models.AuthResponse, error) {
	rt, err := repositories.GetRefreshToken(refreshTokenStr)
	if err != nil {
		return nil, errors.ErrUnauthorized("refresh token inválido ou expirado")
	}

	// Rotação: revoga o token atual antes de emitir novo
	if err := repositories.RevokeRefreshToken(refreshTokenStr); err != nil {
		return nil, errors.ErrInternalWithErr("erro ao revogar token", err)
	}

	user, err := repositories.GetUserByID(rt.UserID)
	if err != nil || user == nil {
		return nil, errors.ErrUnauthorized("usuário não encontrado")
	}

	newToken, err := s.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.ErrInternalWithErr("erro ao gerar token", err)
	}

	newRefreshToken, err := s.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errors.ErrInternalWithErr("erro ao gerar refresh token", err)
	}

	return &models.AuthResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Token:        newToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// Logout revoga todos os refresh tokens ativos do usuário
func (s *AuthService) Logout(userID uint) error {
	return repositories.RevokeAllUserTokens(userID)
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrUnauthorized("método de assinatura inválido")
		}
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil {
		return nil, errors.ErrUnauthorizedWithErr("erro ao validar token", err)
	}

	if !token.Valid {
		return nil, errors.ErrUnauthorized("token inválido")
	}

	return claims, nil
}
