package graph

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/abdulahshoaib/gql-blogs/graph/model"
)

type Resolver struct {
	mu       sync.RWMutex
	posts    []*model.Post
	comments []*model.Comment
	users    map[string]*model.User
	nextID   atomic.Int64
}

func NewResolver() *Resolver {
	return &Resolver{
		users: map[string]*model.User{
			"1": {ID: "1", Name: "Alice"},
			"2": {ID: "2", Name: "Bob"},
		},
	}
}

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }

func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*model.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[input.AuthorID]; !ok {
		return nil, fmt.Errorf("user not found: %s", input.AuthorID)
	}

	post := &model.Post{
		ID:       fmt.Sprintf("%d", r.nextID.Add(1)),
		Title:    input.Title,
		Body:     input.Body,
		Author:   r.users[input.AuthorID],
		Comments: []*model.Comment{},
	}
	r.posts = append(r.posts, post)
	return post, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, input model.NewComment) (*model.Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[input.AuthorID]; !ok {
		return nil, fmt.Errorf("user not found: %s", input.AuthorID)
	}

	var post *model.Post
	for _, p := range r.posts {
		if p.ID == input.PostID {
			post = p
			break
		}
	}
	if post == nil {
		return nil, fmt.Errorf("post not found: %s", input.PostID)
	}

	comment := &model.Comment{
		ID:     fmt.Sprintf("%d", r.nextID.Add(1)),
		Body:   input.Body,
		Author: r.users[input.AuthorID],
		Post:   post,
	}
	r.comments = append(r.comments, comment)
	post.Comments = append(post.Comments, comment)
	return comment, nil
}

func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.posts, nil
}

func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.posts {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, nil
}
