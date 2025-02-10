package service

import "my_app/internal/gateway"

type PostsService struct {
	repo gateway.Posts
}

func NewPostsService(repo gateway.Posts) *PostsService {
	return &PostsService{repo: repo}
}
