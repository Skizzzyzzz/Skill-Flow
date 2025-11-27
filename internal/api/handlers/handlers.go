package handlers

import (
	"github.com/vern/skillflow/internal/config"
	"github.com/vern/skillflow/internal/service"
	"github.com/vern/skillflow/pkg/logger"
)

type Handlers struct {
	Auth         *AuthHandler
	User         *UserHandler
	Post         *PostHandler
	Comment      *CommentHandler
	Reaction     *ReactionHandler
	Connection   *ConnectionHandler
	Notification *NotificationHandler
	Message      *MessageHandler
	Group        *GroupHandler
	Skill        *SkillHandler
	File         *FileHandler
	WebSocket    *WebSocketHandler
	Admin        *AdminHandler
}

func NewHandlers(services *service.Services, cfg *config.Config, log *logger.Logger) *Handlers {
	return &Handlers{
		Auth:         NewAuthHandler(services, cfg, log),
		User:         NewUserHandler(services, log),
		Post:         NewPostHandler(services, log),
		Comment:      NewCommentHandler(services, log),
		Reaction:     NewReactionHandler(services, log),
		Connection:   NewConnectionHandler(services, log),
		Notification: NewNotificationHandler(services, log),
		Message:      NewMessageHandler(services, log),
		Group:        NewGroupHandler(services, log),
		Skill:        NewSkillHandler(services, log),
		File:         NewFileHandler(services, log),
		WebSocket:    NewWebSocketHandler(services, log),
		Admin:        NewAdminHandler(services, log),
	}
}
