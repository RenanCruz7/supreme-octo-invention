package mappers

import (
	"awesomeProject/models"
	"golang.org/x/crypto/bcrypt"
)

func RegisterRequestToUser(req models.RegisterRequest) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}, nil
}
