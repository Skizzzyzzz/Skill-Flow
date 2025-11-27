package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vern/skillflow/internal/service"
	"github.com/vern/skillflow/pkg/logger"
)

type PostHandler struct {
	services *service.Services
	logger   *logger.Logger
}

func NewPostHandler(services *service.Services, log *logger.Logger) *PostHandler {
	return &PostHandler{services: services, logger: log}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	userID := c.GetUint("user_id")
	var input service.CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.UserID = userID

	post, err := h.services.Post.Create(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}
	c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) GetFeed(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	posts, err := h.services.Post.GetFeed(c.Request.Context(), userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feed"})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) GetPostByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	post, err := h.services.Post.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var input service.UpdatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := h.services.Post.Update(c.Request.Context(), uint(id), userID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}
	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.services.Post.Delete(c.Request.Context(), uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}

func (h *PostHandler) GetUserPosts(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	posts, err := h.services.Post.GetUserPosts(c.Request.Context(), uint(userID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts"})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) GetGroupPosts(c *gin.Context) {
	groupID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	posts, err := h.services.Post.GetGroupPosts(c.Request.Context(), uint(groupID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get group posts"})
		return
	}
	c.JSON(http.StatusOK, posts)
}
