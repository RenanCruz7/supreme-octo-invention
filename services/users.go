package services

import (
	"awesomeProject/errors"
	"awesomeProject/models"
	"awesomeProject/repositories"
)

type UserService struct {
	UserRepo repositories.IUserRepository
}

func NewUserService(userRepo repositories.IUserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (s *UserService) GetAllUsers(page, limit int) (*models.PaginatedResponse[models.User], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	users, err := s.UserRepo.GetAllUsers(page, limit)
	if err != nil {
		return nil, errors.ErrInternalWithErr("erro ao buscar usuários", err)
	}
	total, err := s.UserRepo.CountAllUsers()
	if err != nil {
		return nil, errors.ErrInternalWithErr("erro ao contar usuários", err)
	}
	resp := models.NewPaginatedResponse(users, page, limit, total)
	return &resp, nil
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	if id == 0 {
		return nil, errors.ErrBadRequest("ID inválido")
	}

	user, err := s.UserRepo.GetUserByID(id)
	if err != nil || user == nil {
		return nil, errors.ErrNotFound("usuário não encontrado")
	}
	return user, nil
}

func (s *UserService) UpdateUser(id uint, req models.UpdateUserRequest) (*models.User, error) {
	if id == 0 {
		return nil, errors.ErrBadRequest("ID inválido")
	}

	existingUser, err := s.UserRepo.GetUserByID(id)
	if err != nil || existingUser == nil {
		return nil, errors.ErrNotFound("usuário não encontrado")
	}

	if req.Name != nil && *req.Name != "" {
		existingUser.Name = *req.Name
	}
	if req.Email != nil && *req.Email != "" && *req.Email != existingUser.Email {
		emailExists, _ := s.UserRepo.GetUserByEmail(*req.Email)
		if emailExists != nil {
			return nil, errors.ErrConflict("email já cadastrado")
		}
		existingUser.Email = *req.Email
	}

	err = s.UserRepo.UpdateUser(*existingUser)
	if err != nil {
		return nil, errors.ErrInternalWithErr("erro ao atualizar usuário", err)
	}

	return existingUser, nil
}

func (s *UserService) DeleteUser(id uint) error {
	if id == 0 {
		return errors.ErrBadRequest("ID inválido")
	}

	_, err := s.UserRepo.GetUserByID(id)
	if err != nil {
		return errors.ErrNotFound("usuário não encontrado")
	}

	err = s.UserRepo.DeleteUser(id)
	if err != nil {
		return errors.ErrInternalWithErr("erro ao deletar usuário", err)
	}

	return nil
}
