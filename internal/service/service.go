package service

import (
	"github.com/redis/go-redis/v9"
	"github.com/vern/skillflow/internal/config"
	"github.com/vern/skillflow/internal/repository"
	"github.com/vern/skillflow/pkg/logger"
)

type Services struct {
	Auth         *AuthService
	User         *UserService
	Post         *PostService
	Comment      *CommentService
	Reaction     *ReactionService
	Connection   *ConnectionService
	Notification *NotificationService
	Message      *MessageService
	Group        *GroupService
	Skill        *SkillService
	File         *FileService
}

type ServicesDeps struct {
	Repos  *repository.Repositories
	Cache  *redis.Client
	Config *config.Config
	Logger *logger.Logger
}

func NewServices(deps ServicesDeps) *Services {
	return &Services{
		Auth:         NewAuthService(deps),
		User:         NewUserService(deps),
		Post:         NewPostService(deps),
		Comment:      NewCommentService(deps),
		Reaction:     NewReactionService(deps),
		Connection:   NewConnectionService(deps),
		Notification: NewNotificationService(deps),
		Message:      NewMessageService(deps),
		Group:        NewGroupService(deps),
		Skill:        NewSkillService(deps),
		File:         NewFileService(deps),
	}
}
