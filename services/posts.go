package services

import (
	"fmt"

	"awesomeProject/models"
	"awesomeProject/repositories"
)

type PostService struct{}

func (s *PostService) CreatePost(post models.Post, userID uint) (*models.Post, error) {
	if userID == 0 {
		return nil, fmt.Errorf("usuário não autenticado")
	}

	user, err := repositories.GetUserByID(userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("usuário não encontrado")
	}

	post.UserID = userID
	id, err := repositories.CreatePost(post)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar post: %v", err)
	}

	post.ID = id
	return &post, nil
}

func (s *PostService) GetAllPosts() ([]models.Post, error) {
	posts, err := repositories.GetAllPosts()
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar posts: %v", err)
	}
	return posts, nil
}

func (s *PostService) GetPostByID(id uint) (*models.Post, error) {
	if id == 0 {
		return nil, fmt.Errorf("ID inválido")
	}

	post, err := repositories.GetPostByID(id)
	if err != nil || post == nil {
		return nil, fmt.Errorf("post não encontrado")
	}
	return post, nil
}

func (s *PostService) GetUserPosts(userID uint) ([]models.Post, error) {
	if userID == 0 {
		return nil, fmt.Errorf("ID de usuário inválido")
	}

	user, err := repositories.GetUserByID(userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("usuário não encontrado")
	}

	posts, err := repositories.GetPostsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar posts: %v", err)
	}
	return posts, nil
}

func (s *PostService) UpdatePost(id uint, post models.Post, userID uint) (*models.Post, error) {
	if id == 0 {
		return nil, fmt.Errorf("ID inválido")
	}

	existingPost, err := repositories.GetPostByID(id)
	if err != nil || existingPost == nil {
		return nil, fmt.Errorf("post não encontrado")
	}

	if existingPost.UserID != userID {
		return nil, fmt.Errorf("você não tem permissão para atualizar este post")
	}

	if post.Title != "" {
		existingPost.Title = post.Title
	}
	if post.Body != "" {
		existingPost.Body = post.Body
	}

	err = repositories.UpdatePost(*existingPost)
	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar post: %v", err)
	}

	return existingPost, nil
}

func (s *PostService) DeletePost(id uint, userID uint) error {
	if id == 0 {
		return fmt.Errorf("ID inválido")
	}

	post, err := repositories.GetPostByID(id)
	if err != nil || post == nil {
		return fmt.Errorf("post não encontrado")
	}

	if post.UserID != userID {
		return fmt.Errorf("você não tem permissão para deletar este post")
	}

	err = repositories.DeletePost(id)
	if err != nil {
		return fmt.Errorf("erro ao deletar post: %v", err)
	}

	return nil
}
