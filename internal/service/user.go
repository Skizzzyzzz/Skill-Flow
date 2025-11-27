package service

import (
	"context"
	"errors"

	"github.com/vern/skillflow/internal/domain/models"
)

type UserService struct {
	deps ServicesDeps
}

func NewUserService(deps ServicesDeps) *UserService {
	return &UserService{deps: deps}
}

type UpdateUserInput struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type UpdateProfileInput struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Bio        string `json:"bio"`
	Department string `json:"department"`
	Position   string `json:"position"`
	Location   string `json:"location"`
	Phone      string `json:"phone"`
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*models.User, error) {
	return s.deps.Repos.User.GetByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, id uint, input UpdateUserInput) (*models.User, error) {
	user, err := s.deps.Repos.User.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Username != "" {
		user.Username = input.Username
	}

	if err := s.deps.Repos.User.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetProfile(ctx context.Context, userID uint) (*models.Profile, error) {
	return s.deps.Repos.Profile.GetByUserID(ctx, userID)
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uint, input UpdateProfileInput) (*models.Profile, error) {
	profile, err := s.deps.Repos.Profile.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	profile.FirstName = input.FirstName
	profile.LastName = input.LastName
	profile.Bio = input.Bio
	profile.Department = input.Department
	profile.Position = input.Position
	profile.Location = input.Location
	profile.Phone = input.Phone

	if err := s.deps.Repos.Profile.Update(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *UserService) Search(ctx context.Context, query string) ([]models.User, error) {
	return s.deps.Repos.User.Search(ctx, query)
}

type PostService struct {
	deps ServicesDeps
}

func NewPostService(deps ServicesDeps) *PostService {
	return &PostService{deps: deps}
}

type CreatePostInput struct {
	UserID     uint   `json:"-"`
	Content    string `json:"content" binding:"required"`
	MediaURLs  string `json:"media_urls"`
	Visibility string `json:"visibility"`
	GroupID    *uint  `json:"group_id"`
}

type UpdatePostInput struct {
	Content    string `json:"content"`
	Visibility string `json:"visibility"`
}

func (s *PostService) Create(ctx context.Context, input CreatePostInput) (*models.Post, error) {
	post := &models.Post{
		UserID:     input.UserID,
		Content:    input.Content,
		MediaURLs:  input.MediaURLs,
		Visibility: input.Visibility,
		GroupID:    input.GroupID,
	}

	if err := s.deps.Repos.Post.Create(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) GetByID(ctx context.Context, id uint) (*models.Post, error) {
	return s.deps.Repos.Post.GetByID(ctx, id)
}

func (s *PostService) GetFeed(ctx context.Context, userID uint, page, limit int) ([]models.Post, error) {
	return s.deps.Repos.Post.GetFeed(ctx, userID, page, limit)
}

func (s *PostService) GetUserPosts(ctx context.Context, userID uint, page, limit int) ([]models.Post, error) {
	return s.deps.Repos.Post.GetByUserID(ctx, userID, page, limit)
}

func (s *PostService) GetGroupPosts(ctx context.Context, groupID uint, page, limit int) ([]models.Post, error) {
	return s.deps.Repos.Post.GetByGroupID(ctx, groupID, page, limit)
}

func (s *PostService) Update(ctx context.Context, id, userID uint, input UpdatePostInput) (*models.Post, error) {
	post, err := s.deps.Repos.Post.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if post.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	post.Content = input.Content
	if input.Visibility != "" {
		post.Visibility = input.Visibility
	}

	if err := s.deps.Repos.Post.Update(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) Delete(ctx context.Context, id, userID uint) error {
	post, err := s.deps.Repos.Post.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if post.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.deps.Repos.Post.Delete(ctx, id)
}

// Placeholder services
type CommentService struct{ deps ServicesDeps }

func NewCommentService(deps ServicesDeps) *CommentService { return &CommentService{deps: deps} }

type ReactionService struct{ deps ServicesDeps }

func NewReactionService(deps ServicesDeps) *ReactionService { return &ReactionService{deps: deps} }

type ConnectionService struct{ deps ServicesDeps }

func NewConnectionService(deps ServicesDeps) *ConnectionService {
	return &ConnectionService{deps: deps}
}

type NotificationService struct{ deps ServicesDeps }

func NewNotificationService(deps ServicesDeps) *NotificationService {
	return &NotificationService{deps: deps}
}

type MessageService struct{ deps ServicesDeps }

func NewMessageService(deps ServicesDeps) *MessageService { return &MessageService{deps: deps} }

type GroupService struct{ deps ServicesDeps }

func NewGroupService(deps ServicesDeps) *GroupService { return &GroupService{deps: deps} }

type SkillService struct{ deps ServicesDeps }

func NewSkillService(deps ServicesDeps) *SkillService { return &SkillService{deps: deps} }

type FileService struct{ deps ServicesDeps }

func NewFileService(deps ServicesDeps) *FileService { return &FileService{deps: deps} }
