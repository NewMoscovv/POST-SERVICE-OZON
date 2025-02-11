package graphql

import "my_app/internal/service"

type Resolver struct {
	PostService    service.Posts
	CommentService service.Comments
}
