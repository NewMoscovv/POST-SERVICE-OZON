package service

import "my_app/internal/gateway"

type CommentsService struct {
	repo gateway.Comments
}

func NewCommentsService(repo gateway.Comments) *CommentsService {
	return &CommentsService{repo: repo}
}
