package graphql

import "my_app/internal/service"

type Resolver struct {
	service.Posts
	service.Comments
}
