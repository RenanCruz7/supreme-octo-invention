package services

import (
	"awesomeProject/errors"
	"awesomeProject/models"
	"awesomeProject/repositories"
)

type PostService struct {
	UserRepo repositories.IUserRepository
	PostRepo repositories.IPostRepository
}

func NewPostService(userRepo repositories.IUserRepository, postRepo repositories.IPostRepository) *PostService {
	return &PostService{UserRepo: userRepo, PostRepo: postRepo}
}

func (s *PostService) CreatePost(post models.Post, userID uint) (*models.Post, error) {
	if userID == 0 {
		return nil, errors.ErrUnauthorized("usuário não autenticado")
	}

	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil || user == nil {
		return nil, errors.ErrNotFound("usuário não encontrado")
	}

	post.UserID = userID
	id, err := s.PostRepo.CreatePost(post)
	if err != nil {
		return nil, errors.ErrInternalWithErr("erro ao criar post", err)
	}

	post.ID = id
	return &post, nil
}

func (s *PostService) GetAllPosts(page, limit int) (*models.PaginatedResponse[models.Post], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	posts, err := s.PostRepo.GetAllPosts(page, limit)
	if err != nil {
		return nil, errors.ErrInternalWithErr("erro ao buscar posts", err)
	}
	total, err := s.PostRepo.CountAllPosts()
	if err != nil {
		return nil, errors.ErrInternalWithErr("erro ao contar posts", err)
	}
	resp := models.NewPaginatedResponse(posts, page, limit, total)
	return &resp, nil
}

func (s *PostService) GetPostByID(id uint) (*models.Post, error) {
	if id == 0 {
		return nil, errors.ErrBadRequest("ID inválido")
	}

	post, err := s.PostRepo.GetPostByID(id)
	if err != nil || post == nil {
		return nil, errors.ErrNotFound("post não encontrado")
	}
	return post, nil
}

func (s *PostService) GetUserPosts(userID uint, page, limit int) (*models.PaginatedResponse[models.Post], error) {
	if userID == 0 {
		return nil, errors.ErrBadRequest("ID de usuário inválido")
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil || user == nil {
		return nil, errors.ErrNotFound("usuário não encontrado")
	}

	posts, err := s.PostRepo.GetPostsByUserID(userID, page, limit)
	if err != nil {
		return nil, errors.ErrInternalWithErr("erro ao buscar posts", err)
	}
	total, err := s.PostRepo.CountPostsByUserID(userID)
	if err != nil {
		return nil, errors.ErrInternalWithErr("erro ao contar posts", err)
	}
	resp := models.NewPaginatedResponse(posts, page, limit, total)
	return &resp, nil
}

func (s *PostService) UpdatePost(id uint, post models.Post, userID uint) (*models.Post, error) {
	if id == 0 {
		return nil, errors.ErrBadRequest("ID inválido")
	}

	existingPost, err := s.PostRepo.GetPostByID(id)
	if err != nil || existingPost == nil {
		return nil, errors.ErrNotFound("post não encontrado")
	}

	if existingPost.UserID != userID {
		return nil, errors.ErrForbidden("você não tem permissão para atualizar este post")
	}

	if post.Title != "" {
		existingPost.Title = post.Title
	}
	if post.Body != "" {
		existingPost.Body = post.Body
	}

	err = s.PostRepo.UpdatePost(*existingPost)
	if err != nil {
		return nil, errors.ErrInternalWithErr("erro ao atualizar post", err)
	}

	return existingPost, nil
}

func (s *PostService) DeletePost(id uint, userID uint) error {
	if id == 0 {
		return errors.ErrBadRequest("ID inválido")
	}

	post, err := s.PostRepo.GetPostByID(id)
	if err != nil || post == nil {
		return errors.ErrNotFound("post não encontrado")
	}

	if post.UserID != userID {
		return errors.ErrForbidden("você não tem permissão para deletar este post")
	}

	err = s.PostRepo.DeletePost(id)
	if err != nil {
		return errors.ErrInternalWithErr("erro ao deletar post", err)
	}

	return nil
}
