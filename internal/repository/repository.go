package repository

import (
	"context"

	"github.com/vern/skillflow/internal/domain/models"
	"gorm.io/gorm"
)

type Repositories struct {
	User         UserRepositoryInterface
	Profile      ProfileRepositoryInterface
	Post         PostRepositoryInterface
	Comment      CommentRepositoryInterface
	Reaction     ReactionRepositoryInterface
	Connection   ConnectionRepositoryInterface
	Notification NotificationRepositoryInterface
	Message      MessageRepositoryInterface
	Group        GroupRepositoryInterface
	GroupMember  GroupMemberRepositoryInterface
	Skill        SkillRepositoryInterface
	UserSkill    UserSkillRepositoryInterface
	Endorsement  EndorsementRepositoryInterface
	File         FileRepositoryInterface
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:         &UserRepository{db: db},
		Profile:      &ProfileRepository{db: db},
		Post:         &PostRepository{db: db},
		Comment:      &CommentRepository{db: db},
		Reaction:     &ReactionRepository{db: db},
		Connection:   &ConnectionRepository{db: db},
		Notification: &NotificationRepository{db: db},
		Message:      &MessageRepository{db: db},
		Group:        &GroupRepository{db: db},
		GroupMember:  &GroupMemberRepository{db: db},
		Skill:        &SkillRepository{db: db},
		UserSkill:    &UserSkillRepository{db: db},
		Endorsement:  &EndorsementRepository{db: db},
		File:         &FileRepository{db: db},
	}
}

// Interfaces
type UserRepositoryInterface interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	Search(ctx context.Context, query string) ([]models.User, error)
}

type ProfileRepositoryInterface interface {
	Create(ctx context.Context, profile *models.Profile) error
	GetByUserID(ctx context.Context, userID uint) (*models.Profile, error)
	Update(ctx context.Context, profile *models.Profile) error
}

type PostRepositoryInterface interface {
	Create(ctx context.Context, post *models.Post) error
	GetByID(ctx context.Context, id uint) (*models.Post, error)
	GetFeed(ctx context.Context, userID uint, page, limit int) ([]models.Post, error)
	GetByUserID(ctx context.Context, userID uint, page, limit int) ([]models.Post, error)
	GetByGroupID(ctx context.Context, groupID uint, page, limit int) ([]models.Post, error)
	Update(ctx context.Context, post *models.Post) error
	Delete(ctx context.Context, id uint) error
}

type CommentRepositoryInterface interface{}
type ReactionRepositoryInterface interface{}
type ConnectionRepositoryInterface interface{}
type NotificationRepositoryInterface interface{}
type MessageRepositoryInterface interface{}
type GroupRepositoryInterface interface{}
type GroupMemberRepositoryInterface interface{}
type SkillRepositoryInterface interface{}
type UserSkillRepositoryInterface interface{}
type EndorsementRepositoryInterface interface{}
type FileRepositoryInterface interface{}

// Implementations
type UserRepository struct{ db *gorm.DB }
type ProfileRepository struct{ db *gorm.DB }
type PostRepository struct{ db *gorm.DB }
type CommentRepository struct{ db *gorm.DB }
type ReactionRepository struct{ db *gorm.DB }
type ConnectionRepository struct{ db *gorm.DB }
type NotificationRepository struct{ db *gorm.DB }
type MessageRepository struct{ db *gorm.DB }
type GroupRepository struct{ db *gorm.DB }
type GroupMemberRepository struct{ db *gorm.DB }
type SkillRepository struct{ db *gorm.DB }
type UserSkillRepository struct{ db *gorm.DB }
type EndorsementRepository struct{ db *gorm.DB }
type FileRepository struct{ db *gorm.DB }

// User repository methods
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Preload("Profile").First(&user, id).Error
	return &user, err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (r *UserRepository) Search(ctx context.Context, query string) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).
		Preload("Profile").
		Where("username LIKE ? OR email LIKE ?", "%"+query+"%", "%"+query+"%").
		Find(&users).Error
	return users, err
}

// Profile repository methods
func (r *ProfileRepository) Create(ctx context.Context, profile *models.Profile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

func (r *ProfileRepository) GetByUserID(ctx context.Context, userID uint) (*models.Profile, error) {
	var profile models.Profile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error
	return &profile, err
}

func (r *ProfileRepository) Update(ctx context.Context, profile *models.Profile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

// Post repository methods
func (r *PostRepository) Create(ctx context.Context, post *models.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *PostRepository) GetByID(ctx context.Context, id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.WithContext(ctx).
		Preload("User.Profile").
		Preload("Comments").
		Preload("Reactions").
		First(&post, id).Error
	return &post, err
}

func (r *PostRepository) GetFeed(ctx context.Context, userID uint, page, limit int) ([]models.Post, error) {
	var posts []models.Post
	offset := (page - 1) * limit
	err := r.db.WithContext(ctx).
		Preload("User.Profile").
		Preload("Reactions").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error
	return posts, err
}

func (r *PostRepository) GetByUserID(ctx context.Context, userID uint, page, limit int) ([]models.Post, error) {
	var posts []models.Post
	offset := (page - 1) * limit
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("User.Profile").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error
	return posts, err
}

func (r *PostRepository) GetByGroupID(ctx context.Context, groupID uint, page, limit int) ([]models.Post, error) {
	var posts []models.Post
	offset := (page - 1) * limit
	err := r.db.WithContext(ctx).
		Where("group_id = ?", groupID).
		Preload("User.Profile").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error
	return posts, err
}

func (r *PostRepository) Update(ctx context.Context, post *models.Post) error {
	return r.db.WithContext(ctx).Save(post).Error
}

func (r *PostRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Post{}, id).Error
}
