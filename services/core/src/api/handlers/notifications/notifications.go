package notifications

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/truckguard/core/src/repository"
)

func getPostID(c *gin.Context) (uint, bool) {
	postIDStr, _ := c.Request.Context().Value("user_post_id").(string)
	if postIDStr == "" || postIDStr == "nil" {
		return 0, false
	}
	id, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		return 0, false
	}
	return uint(id), true
}

// List returns a list of notifications for the current post.
func List(c *gin.Context) {
	postID, ok := getPostID(c)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)
	offsetStr := c.DefaultQuery("offset", "0")
	offset, _ := strconv.Atoi(offsetStr)
	unreadOnly := c.Query("unread") == "true"

	notifs, err := repository.GetNotifications(c.Request.Context(), postID, limit, offset, unreadOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": notifs})
}

// MarkRead marks a single notification as read.
func MarkRead(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification ID"})
		return
	}

	if err := repository.MarkNotificationRead(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// MarkAllRead marks all notifications as read for the current post.
func MarkAllRead(c *gin.Context) {
	postID, ok := getPostID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "post_id not found in context"})
		return
	}

	if err := repository.MarkAllNotificationsRead(c.Request.Context(), postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
