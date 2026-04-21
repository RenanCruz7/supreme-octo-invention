package services

import (
	"fmt"

	"awesomeProject/models"
	"awesomeProject/repositories"
)

type UserService struct{}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	users, err := repositories.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuários: %v", err)
	}
	return users, nil
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	if id == 0 {
		return nil, fmt.Errorf("ID inválido")
	}

	user, err := repositories.GetUserByID(id)
	if err != nil || user == nil {
		return nil, fmt.Errorf("usuário não encontrado")
	}
	return user, nil
}

func (s *UserService) UpdateUser(id uint, req models.User) (*models.User, error) {
	if id == 0 {
		return nil, fmt.Errorf("ID inválido")
	}

	existingUser, err := repositories.GetUserByID(id)
	if err != nil || existingUser == nil {
		return nil, fmt.Errorf("usuário não encontrado")
	}

	if req.Name != "" {
		existingUser.Name = req.Name
	}
	if req.Email != "" && req.Email != existingUser.Email {
		emailExists, _ := repositories.GetUserByEmail(req.Email)
		if emailExists != nil {
			return nil, fmt.Errorf("email já cadastrado")
		}
		existingUser.Email = req.Email
	}

	err = repositories.UpdateUser(*existingUser)
	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar usuário: %v", err)
	}

	return existingUser, nil
}

func (s *UserService) DeleteUser(id uint) error {
	if id == 0 {
		return fmt.Errorf("ID inválido")
	}

	_, err := repositories.GetUserByID(id)
	if err != nil {
		return fmt.Errorf("usuário não encontrado")
	}

	err = repositories.DeleteUser(id)
	if err != nil {
		return fmt.Errorf("erro ao deletar usuário: %v", err)
	}

	return nil
}
