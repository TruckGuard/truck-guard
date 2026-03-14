package data

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/truckguard/core/src/models"
	datarepo "github.com/truckguard/core/src/repository/data"
	"github.com/truckguard/core/src/utils"
	"gorm.io/gorm"
)

func HandleListPosts(c *gin.Context) {
	limit, offset, page := utils.GetPagination(c)
	filter := datarepo.PostFilter{
		Name: c.Query("name"),
	}
	posts, total, err := datarepo.ListPosts(c.Request.Context(), limit, offset, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}
	utils.SendPaginatedResponse(c, posts, total, page, limit)
}

func HandleCreatePost(c *gin.Context) {
	var input models.CustomsPost
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := datarepo.CreatePost(c.Request.Context(), &input); err != nil {
		slog.Error("Failed to create post", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, input)
}

func HandleGetPostByID(c *gin.Context) {
	id := c.Param("id")
	post, err := datarepo.GetPost(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Митний пост не знайдено"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка при отриманні митного поста"})
		}
		return
	}
	c.JSON(http.StatusOK, post)
}

func HandleUpdatePost(c *gin.Context) {
	id := c.Param("id")
	var input models.CustomsPost
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post, err := datarepo.UpdatePost(c.Request.Context(), id, &input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Митний пост не знайдено"})
		} else {
			slog.Error("Failed to update post", "error", err, "id", id)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка при оновленні митного поста: " + err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, post)
}

func HandleDeletePost(c *gin.Context) {
	id := c.Param("id")
	if err := datarepo.DeletePost(c.Request.Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Митний пост не знайдено"})
		} else {
			slog.Error("Failed to delete post", "error", err, "id", id)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка при видаленні митного поста: " + err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}
