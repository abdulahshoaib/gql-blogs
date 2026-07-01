package graph

import (
	"context"
	"fmt"
	"strconv"

	"github.com/abdulahshoaib/gql-blogs/graph/model"
	internalModel "github.com/abdulahshoaib/gql-blogs/internal/model"
	"gorm.io/gorm"
)

type Resolver struct {
	DB *gorm.DB
}

func NewResolver(db *gorm.DB) *Resolver {
	return &Resolver{DB: db}
}

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }

func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*model.Post, error) {
	authorID, err := strconv.ParseUint(input.AuthorID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid author ID: %s", input.AuthorID)
	}

	var author internalModel.User
	if err := r.DB.First(&author, authorID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %s", input.AuthorID)
	}

	post := internalModel.Post{
		Title:    input.Title,
		Body:     input.Body,
		AuthorID: uint(authorID),
	}
	if err := r.DB.Create(&post).Error; err != nil {
		return nil, err
	}

	r.DB.Preload("Author").First(&post, post.ID)

	return &model.Post{
		ID:       fmt.Sprintf("%d", post.ID),
		Title:    post.Title,
		Body:     post.Body,
		Author:   &model.User{ID: fmt.Sprintf("%d", post.Author.ID), Name: post.Author.Name},
		Comments: []*model.Comment{},
	}, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, input model.NewComment) (*model.Comment, error) {
	authorID, err := strconv.ParseUint(input.AuthorID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid author ID: %s", input.AuthorID)
	}
	postID, err := strconv.ParseUint(input.PostID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid post ID: %s", input.PostID)
	}

	var author internalModel.User
	if err := r.DB.First(&author, authorID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %s", input.AuthorID)
	}

	var post internalModel.Post
	if err := r.DB.First(&post, postID).Error; err != nil {
		return nil, fmt.Errorf("post not found: %s", input.PostID)
	}

	comment := internalModel.Comment{
		Body:     input.Body,
		AuthorID: uint(authorID),
		PostID:   uint(postID),
	}
	if err := r.DB.Create(&comment).Error; err != nil {
		return nil, err
	}

	r.DB.Preload("Author").Preload("Post.Author").First(&comment, comment.ID)

	gqlPost := &model.Post{
		ID:     fmt.Sprintf("%d", comment.Post.ID),
		Title:  comment.Post.Title,
		Body:   comment.Post.Body,
		Author: &model.User{ID: fmt.Sprintf("%d", comment.Post.Author.ID), Name: comment.Post.Author.Name},
	}

	return &model.Comment{
		ID:     fmt.Sprintf("%d", comment.ID),
		Body:   comment.Body,
		Author: &model.User{ID: fmt.Sprintf("%d", comment.Author.ID), Name: comment.Author.Name},
		Post:   gqlPost,
	}, nil
}

func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	var posts []internalModel.Post
	if err := r.DB.
		Preload("Author").
		Preload("Comments").
		Preload("Comments.Author").
		Order("created_at desc").
		Find(&posts).Error; err != nil {
		return nil, err
	}

	result := make([]*model.Post, len(posts))
	for i, p := range posts {
		comments := make([]*model.Comment, len(p.Comments))
		for j, c := range p.Comments {
			comments[j] = &model.Comment{
				ID:     fmt.Sprintf("%d", c.ID),
				Body:   c.Body,
				Author: &model.User{ID: fmt.Sprintf("%d", c.Author.ID), Name: c.Author.Name},
			}
		}
		result[i] = &model.Post{
			ID:       fmt.Sprintf("%d", p.ID),
			Title:    p.Title,
			Body:     p.Body,
			Author:   &model.User{ID: fmt.Sprintf("%d", p.Author.ID), Name: p.Author.Name},
			Comments: comments,
		}
	}
	return result, nil
}

func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	postID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, nil
	}

	var post internalModel.Post
	if err := r.DB.
		Preload("Author").
		Preload("Comments").
		Preload("Comments.Author").
		First(&post, postID).Error; err != nil {
		return nil, nil
	}

	comments := make([]*model.Comment, len(post.Comments))
	for j, c := range post.Comments {
		comments[j] = &model.Comment{
			ID:     fmt.Sprintf("%d", c.ID),
			Body:   c.Body,
			Author: &model.User{ID: fmt.Sprintf("%d", c.Author.ID), Name: c.Author.Name},
		}
	}

	return &model.Post{
		ID:       fmt.Sprintf("%d", post.ID),
		Title:    post.Title,
		Body:     post.Body,
		Author:   &model.User{ID: fmt.Sprintf("%d", post.Author.ID), Name: post.Author.Name},
		Comments: comments,
	}, nil
}
