package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/truckguard/core/src/pkg/notify"
)

// HandleSSENotifications streams permit events to the connected operator
// via Server-Sent Events. The operator's post is read from the request context
// (injected by AuthContextMiddleware as "user_post_id").
func HandleSSENotifications(c *gin.Context) {
	postID, _ := c.Request.Context().Value("user_post_id").(string)
	if postID == "" || postID == "nil" {
		// Operators without a post receive no events
		c.JSON(http.StatusForbidden, gin.H{"error": "No post assigned"})
		return
	}

	// Subscribe to hub
	ch := notify.Global.Subscribe(postID)
	defer notify.Global.Unsubscribe(postID, ch)

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no") // disable nginx buffering

	// Send an initial ping so the client knows the connection is live
	fmt.Fprintf(c.Writer, "data: {\"type\":\"connected\"}\n\n")
	c.Writer.Flush()

	// Keepalive ticker to prevent proxy timeouts
	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()

	clientGone := c.Request.Context().Done()

	for {
		select {
		case <-clientGone:
			return
		case <-ticker.C:
			fmt.Fprintf(c.Writer, ": keepalive\n\n")
			c.Writer.Flush()
		case ev, ok := <-ch:
			if !ok {
				return
			}
			b, _ := json.Marshal(ev)
			fmt.Fprintf(c.Writer, "data: %s\n\n", b)
			c.Writer.Flush()
		}
	}
}
