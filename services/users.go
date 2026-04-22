package services

import (
	"awesomeProject/errors"
	"awesomeProject/models"
	"awesomeProject/repositories"
)

type UserService struct{}

func (s *UserService) GetAllUsers(page, limit int) ([]models.User, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	users, err := repositories.GetAllUsers(page, limit)
	if err != nil {
		return nil, errors.ErrInternalWithErr("erro ao buscar usuários", err)
	}
	return users, nil
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	if id == 0 {
		return nil, errors.ErrBadRequest("ID inválido")
	}

	user, err := repositories.GetUserByID(id)
	if err != nil || user == nil {
		return nil, errors.ErrNotFound("usuário não encontrado")
	}
	return user, nil
}

func (s *UserService) UpdateUser(id uint, req models.User) (*models.User, error) {
	if id == 0 {
		return nil, errors.ErrBadRequest("ID inválido")
	}

	existingUser, err := repositories.GetUserByID(id)
	if err != nil || existingUser == nil {
		return nil, errors.ErrNotFound("usuário não encontrado")
	}

	if req.Name != "" {
		existingUser.Name = req.Name
	}
	if req.Email != "" && req.Email != existingUser.Email {
		emailExists, _ := repositories.GetUserByEmail(req.Email)
		if emailExists != nil {
			return nil, errors.ErrConflict("email já cadastrado")
		}
		existingUser.Email = req.Email
	}

	err = repositories.UpdateUser(*existingUser)
	if err != nil {
		return nil, errors.ErrInternalWithErr("erro ao atualizar usuário", err)
	}

	return existingUser, nil
}

func (s *UserService) DeleteUser(id uint) error {
	if id == 0 {
		return errors.ErrBadRequest("ID inválido")
	}

	_, err := repositories.GetUserByID(id)
	if err != nil {
		return errors.ErrNotFound("usuário não encontrado")
	}

	err = repositories.DeleteUser(id)
	if err != nil {
		return errors.ErrInternalWithErr("erro ao deletar usuário", err)
	}

	return nil
}
