package service

import (
	"my_app/internal/gateway"
	"my_app/internal/logger"
	"my_app/internal/models"
)

type Services struct {
	Posts
	Comments
}

func NewServices(gateways *gateway.Gateways, logger *logger.Logger) *Services {
	return &Services{
		Posts:    NewPostsService(gateways.Posts, logger),
		Comments: NewCommentsService(gateways.Comments, logger),
	}
}

type Posts interface {
	CreatePost(post models.Post) (models.Post, error)
	GetPostById(id int) (models.Post, error)
	GetAllPosts(page, pageSize *int) ([]models.Post, error)
}

type Comments interface {
}
