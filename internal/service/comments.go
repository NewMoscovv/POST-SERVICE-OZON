package service

import (
	"database/sql"
	"errors"
	"my_app/internal/consts"
	"my_app/internal/gateway"
	"my_app/internal/logger"
	"my_app/internal/models"
	"my_app/internal/pagination"
	re "my_app/pkg/responce_errors"
)

type CommentsService struct {
	repo       gateway.Comments
	logger     *logger.Logger
	PostGetter PostGetter
}

type PostGetter interface {
	GetPostById(id int) (models.Post, error)
}

func NewCommentsService(repo gateway.Comments, logger *logger.Logger, getter PostGetter) *CommentsService {
	return &CommentsService{repo: repo, logger: logger, PostGetter: getter}
}

func (c CommentsService) CreateComment(comment models.Comment) (models.Comment, error) {
	if len(comment.Author) == 0 {
		c.logger.Err.Println(consts.EmptyAuthorError)
		return models.Comment{}, re.ResponseError{
			Message: consts.EmptyAuthorError,
			Type:    consts.BadRequestType,
		}
	}

	if len(comment.Content) >= consts.MaxContentLength {
		c.logger.Err.Println(consts.TooLongContentError, len(comment.Content))
		return models.Comment{}, re.ResponseError{
			Message: consts.TooLongContentError,
			Type:    consts.BadRequestType,
		}
	}

	if comment.Post <= 0 {
		c.logger.Err.Println(consts.WrongIdError, comment.Post)
		return models.Comment{}, re.ResponseError{
			Message: consts.WrongIdError,
			Type:    consts.BadRequestType,
		}
	}

	post, err := c.PostGetter.GetPostById(comment.Post)
	if err != nil {
		c.logger.Err.Println(consts.GettingPostError, err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return models.Comment{}, re.ResponseError{
				Message: consts.PostNotFountError,
				Type:    consts.NotFoundType,
			}
		}
	}
	if !post.CommentsAllowed {
		c.logger.Err.Println(consts.CommentsNotAllowedError)
		return models.Comment{}, re.ResponseError{
			Message: consts.CommentsNotAllowedError,
			Type:    consts.BadRequestType,
		}
	}
	newComment, err := c.repo.CreateComment(comment)
	if err != nil {
		c.logger.Err.Println(consts.CreatingCommentError, err.Error())
		return models.Comment{}, re.ResponseError{
			Message: consts.CreatingCommentError,
			Type:    consts.InternalErrorType,
		}
	}

	return newComment, nil
}

func (c CommentsService) GetCommentsByPost(postId int, page *int, pageSize *int) ([]*models.Comment, error) {
	if postId <= 0 {
		c.logger.Err.Println(consts.WrongIdError, postId)
		return nil, re.ResponseError{
			Message: consts.WrongIdError,
			Type:    consts.BadRequestType,
		}
	}

	if page != nil && *page < 0 {
		c.logger.Err.Println(consts.WrongPageError, *page)
		return nil, re.ResponseError{
			Message: consts.WrongPageError,
			Type:    consts.BadRequestType,
		}
	}

	if pageSize != nil && *pageSize < 0 {
		c.logger.Err.Println(consts.WrongPageSizeError, *pageSize)
		return nil, re.ResponseError{
			Message: consts.WrongPageSizeError,
			Type:    consts.BadRequestType,
		}
	}

	offset, limit := pagination.GetOffsetAndLimit(page, pageSize)

	comments, err := c.repo.GetCommentsByPost(postId, limit, offset)
	if err != nil {
		c.logger.Err.Println(consts.GettingCommentError, postId, err.Error())
		return nil, re.ResponseError{
			Message: consts.GettingCommentError,
			Type:    consts.InternalErrorType,
		}
	}
	return comments, nil
}

func (c CommentsService) GetRepliesOfComment(commentId int) ([]*models.Comment, error) {
	if commentId <= 0 {
		c.logger.Err.Println(consts.WrongIdError, commentId)
		return nil, re.ResponseError{
			Message: consts.WrongIdError,
			Type:    consts.BadRequestType,
		}
	}
	comments, err := c.repo.GetRepliesOfComment(commentId)
	if err != nil {
		c.logger.Err.Println(consts.GettingRepliesError, commentId, err.Error())
		return nil, re.ResponseError{
			Message: consts.GettingRepliesError,
			Type:    consts.InternalErrorType,
		}
	}
	return comments, nil
}
