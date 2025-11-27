package router

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vern/skillflow/internal/api/handlers"
	"github.com/vern/skillflow/internal/config"
	"github.com/vern/skillflow/internal/service"
	"github.com/vern/skillflow/pkg/logger"
	"github.com/vern/skillflow/pkg/middleware"
)

func NewRouter(services *service.Services, cfg *config.Config, log *logger.Logger) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware(log))
	router.Use(middleware.CORSMiddleware(cfg.CORS))
	router.Use(middleware.RequestIDMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Metrics endpoint (Prometheus)
	if cfg.Monitoring.Prometheus.Enabled {
		router.GET(cfg.Monitoring.Prometheus.Path, gin.WrapH(promhttp.Handler()))
	}

	// Initialize handlers
	h := handlers.NewHandlers(services, cfg, log)

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", h.Auth.Register)
			auth.POST("/login", h.Auth.Login)
			auth.POST("/refresh", h.Auth.RefreshToken)
			auth.GET("/oidc/login", h.Auth.OIDCLogin)
			auth.GET("/oidc/callback", h.Auth.OIDCCallback)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.Auth.JWTSecret))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", h.User.GetCurrentUser)
				users.PUT("/me", h.User.UpdateCurrentUser)
				users.GET("/:id", h.User.GetUserByID)
				users.GET("/:id/profile", h.User.GetUserProfile)
				users.PUT("/:id/profile", h.User.UpdateUserProfile)
				users.GET("/search", h.User.SearchUsers)
			}

			// Post routes
			posts := protected.Group("/posts")
			{
				posts.POST("", h.Post.CreatePost)
				posts.GET("", h.Post.GetFeed)
				posts.GET("/:id", h.Post.GetPostByID)
				posts.PUT("/:id", h.Post.UpdatePost)
				posts.DELETE("/:id", h.Post.DeletePost)
				posts.GET("/user/:user_id", h.Post.GetUserPosts)

				// Comments
				posts.POST("/:id/comments", h.Comment.CreateComment)
				posts.GET("/:id/comments", h.Comment.GetPostComments)

				// Reactions
				posts.POST("/:id/reactions", h.Reaction.AddReaction)
				posts.DELETE("/:id/reactions", h.Reaction.RemoveReaction)
				posts.GET("/:id/reactions", h.Reaction.GetReactions)
			}

			// Comment routes
			comments := protected.Group("/comments")
			{
				comments.PUT("/:id", h.Comment.UpdateComment)
				comments.DELETE("/:id", h.Comment.DeleteComment)
				comments.POST("/:id/reactions", h.Reaction.AddCommentReaction)
				comments.DELETE("/:id/reactions", h.Reaction.RemoveCommentReaction)
			}

			// Connection routes
			connections := protected.Group("/connections")
			{
				connections.POST("", h.Connection.SendConnectionRequest)
				connections.GET("", h.Connection.GetConnections)
				connections.GET("/pending", h.Connection.GetPendingRequests)
				connections.PUT("/:id/accept", h.Connection.AcceptConnection)
				connections.PUT("/:id/reject", h.Connection.RejectConnection)
				connections.DELETE("/:id", h.Connection.RemoveConnection)
			}

			// Notification routes
			notifications := protected.Group("/notifications")
			{
				notifications.GET("", h.Notification.GetNotifications)
				notifications.PUT("/:id/read", h.Notification.MarkAsRead)
				notifications.PUT("/read-all", h.Notification.MarkAllAsRead)
				notifications.GET("/unread/count", h.Notification.GetUnreadCount)
			}

			// Message routes
			messages := protected.Group("/messages")
			{
				messages.POST("", h.Message.SendMessage)
				messages.GET("/conversations", h.Message.GetConversations)
				messages.GET("/conversation/:user_id", h.Message.GetConversation)
				messages.PUT("/:id/read", h.Message.MarkAsRead)
			}

			// Group routes
			groups := protected.Group("/groups")
			{
				groups.POST("", h.Group.CreateGroup)
				groups.GET("", h.Group.GetGroups)
				groups.GET("/:id", h.Group.GetGroupByID)
				groups.PUT("/:id", h.Group.UpdateGroup)
				groups.DELETE("/:id", h.Group.DeleteGroup)
				groups.POST("/:id/join", h.Group.JoinGroup)
				groups.POST("/:id/leave", h.Group.LeaveGroup)
				groups.GET("/:id/members", h.Group.GetGroupMembers)
				groups.GET("/:id/posts", h.Post.GetGroupPosts)
			}

			// Skill routes
			skills := protected.Group("/skills")
			{
				skills.GET("", h.Skill.GetSkills)
				skills.POST("", h.Skill.CreateSkill)
				skills.POST("/user", h.Skill.AddUserSkill)
				skills.DELETE("/user/:id", h.Skill.RemoveUserSkill)
				skills.PUT("/user/:id", h.Skill.UpdateUserSkill)
				skills.POST("/endorse/:user_skill_id", h.Skill.EndorseSkill)
			}

			// File routes
			files := protected.Group("/files")
			{
				files.POST("/upload", h.File.Upload)
				files.GET("/:id", h.File.GetFile)
				files.DELETE("/:id", h.File.DeleteFile)
			}

			// WebSocket for real-time features
			protected.GET("/ws", h.WebSocket.HandleConnection)
		}

		// Admin routes
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg.Auth.JWTSecret))
		admin.Use(middleware.AdminMiddleware())
		{
			admin.GET("/users", h.Admin.GetAllUsers)
			admin.PUT("/users/:id/activate", h.Admin.ActivateUser)
			admin.PUT("/users/:id/deactivate", h.Admin.DeactivateUser)
			admin.DELETE("/posts/:id", h.Admin.DeletePost)
			admin.DELETE("/comments/:id", h.Admin.DeleteComment)
			admin.GET("/stats", h.Admin.GetStats)
		}
	}

	return router
}
