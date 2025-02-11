package graphql

import "my_app/internal/service"

type Resolver struct {
	PostsService      service.Posts
	CommentsService   service.Comments
	CommentsObservers Observers
}
