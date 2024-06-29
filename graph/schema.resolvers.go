package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"errors"
	"fmt"
	"ozonTech/graph/model"
	"ozonTech/internal/models"
	"strconv"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, title string, content string, commentsAllowed bool) (*model.Post, error) {
	fmt.Println("tut")
	post, err := r.PostUsecase.CreatePost(&models.PostCreateData{
		Title:           title,
		Content:         content,
		CommentsAllowed: commentsAllowed,
		UserID:          1, // Здесь используйте реальный userID из контекста или другого источника
	})
	if err != nil {
		return nil, err
	}
	return ConvertToGraphQLPost(post), nil
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, postID string, parentID *string, content string) (*model.Comment, error) {
	intPostID, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		return nil, errors.New("invalid post id")
	}
	var intParentCommentID int64
	if parentID != nil {
		intParentCommentID, err = strconv.ParseInt(*parentID, 10, 64)
		if err != nil {
			return nil, errors.New("invalid parent comment id")
		}
	}
	comment, err := r.CommentUsecase.CreateComment(&models.CommentCreateData{
		PostID:          int(intPostID),
		ParentCommentID: int(intParentCommentID),
		Text:            content,
		UserID:          1, // Здесь используйте реальный userID из контекста или другого источника
	})
	if err != nil {
		return nil, err
	}
	return ConvertToGraphQLComment(comment), nil
}

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, name string, password string) (string, error) {
	fmt.Println("В резолвере")
	token, err := r.AuthUsecase.SignUp(name, password)
	if err != nil {
		fmt.Println("ошибка тут не равана нил")
		return "", err
	}
	return token, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, name string, password string) (string, error) {
	token, err := r.AuthUsecase.Login(name, password)
	if err != nil {
		return "", err
	}
	return token, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	posts, err := r.PostUsecase.GetAllPosts()
	if err != nil {
		return nil, err
	}

	var gqlPosts []*model.Post
	for _, post := range posts {
		gqlPosts = append(gqlPosts, ConvertToGraphQLPost(post))
	}

	return gqlPosts, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid post ID")
	}

	post, err := r.PostUsecase.GetPostByID(intID)
	if err != nil {
		return nil, err
	}

	return ConvertToGraphQLPost(post), nil
}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	// Создаем канал для передачи комментариев
	commentChan := make(chan *model.Comment)

	// Пример реализации логики подписки
	go func() {
		defer close(commentChan)

		// Здесь может быть ваша логика для подписки на новые комментарии к посту с указанным postID
		// Например, использование Pub/Sub системы или реализация логики для отслеживания новых комментариев

		// В примере мы просто создаем фиктивный комментарий и отправляем его в канал
		fakeComment := &model.Comment{
			ID:      "1", // Замените на реальный ID комментария
			Content: "New comment",
			PostID:  postID,
			UserID:  "1", // Замените на реальный ID пользователя
		}

		// Отправляем фиктивный комментарий в канал
		commentChan <- fakeComment
	}()

	return commentChan, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func ConvertToGraphQLPost(post *models.Post) *model.Post {
	return &model.Post{
		ID:              strconv.Itoa(post.ID),
		Title:           post.Title,
		Content:         post.Content,
		CommentsAllowed: post.CommentsAllowed,
		UserID:          strconv.Itoa(post.UserID),
		Comments:        ConvertToGraphQLComments(post.Comments),
	}
}
func ConvertToGraphQLComment(comment *models.Comment) *model.Comment {
	var parentCommentID *string
	if comment.ParentCommentID != 0 {
		idStr := strconv.Itoa(comment.ParentCommentID)
		parentCommentID = &idStr
	}
	var childComments []*model.Comment
	for _, childID := range comment.ChildComments {
		childComment := &model.Comment{
			ID: strconv.Itoa(childID),
		}
		childComments = append(childComments, childComment)
	}
	return &model.Comment{
		ID:              strconv.Itoa(comment.ID),
		Content:         comment.Text,
		PostID:          strconv.Itoa(comment.PostID),
		UserID:          strconv.Itoa(comment.UserID),
		ParentCommentID: parentCommentID,
		ChildComments:   childComments,
	}
}
func ConvertToGraphQLComments(comments []*models.Comment) []*model.Comment {
	var gqlComments []*model.Comment
	for _, comment := range comments {
		gqlComments = append(gqlComments, ConvertToGraphQLComment(comment))
	}
	return gqlComments
}
