package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/truckguard/core/src/models"
	"github.com/truckguard/core/src/repository"
)

// FuzzyCandidate represents a possible permit match for fuzzy search.
type FuzzyCandidate struct {
	ID         uint   `json:"id"`
	Code       string `json:"code"`
	PlateFront string `json:"plate_front"`
	PlateBack  string `json:"plate_back"`
	EntryTime  string `json:"entry_time"`
	Distance   int    `json:"distance"`
}

// Event is a notification payload sent to subscribers.
type Event struct {
	Type   string `json:"type"`
	NotificationID uint `json:"notification_id,omitempty"`
	ID     uint   `json:"id"`
	Code   string `json:"code"`
	Plate  string `json:"plate"`
	PostID uint   `json:"post_id"`
	// Fuzzy match fields
	IncomingPlate string           `json:"incoming_plate,omitempty"`
	EventID       uint             `json:"event_id,omitempty"`
	ImageKey      string           `json:"image_key,omitempty"`
	Candidates    []FuzzyCandidate `json:"candidates,omitempty"`
	Message       string           `json:"message,omitempty"`
}

// channel name for a given post
func channelName(postID uint) string {
	return fmt.Sprintf("permit_notify:post:%d", postID)
}

// PublishPermit is a convenience wrapper for publishing a new_permit event.
// postID 0 is a no-op.
func (h *Hub) PublishPermit(postID uint, permitID uint, code, plate string) {
	if postID == 0 {
		return
	}
	h.Publish(postID, Event{
		Type:   "new_permit",
		ID:     permitID,
		Code:   code,
		Plate:  plate,
		PostID: postID,
	})
}

// PublishFuzzyMatch publishes a fuzzy_match event with candidates for manual linking.
func (h *Hub) PublishFuzzyMatch(postID uint, eventID uint, incomingPlate, imageKey string, candidates []FuzzyCandidate) {
	if postID == 0 || len(candidates) == 0 {
		return
	}
	h.Publish(postID, Event{
		Type:          "fuzzy_match",
		PostID:        postID,
		EventID:       eventID,
		IncomingPlate: incomingPlate,
		ImageKey:      imageKey,
		Candidates:    candidates,
	})
}

// Hub manages SSE subscribers, backed by Valkey pub/sub.
type Hub struct {
	rdb *redis.Client

	mu   sync.RWMutex
	subs map[string][]chan Event // key = postID string
}

var Global = &Hub{subs: make(map[string][]chan Event)}

// Init wires the hub to the Valkey client and starts the global subscriber goroutine.
func Init(ctx context.Context, rdb *redis.Client) {
	Global.rdb = rdb
	go Global.listenAll(ctx)
}

// listenAll subscribes to all permit_notify channels via Valkey psubscribe.
func (h *Hub) listenAll(ctx context.Context) {
	pubsub := h.rdb.PSubscribe(ctx, "permit_notify:post:*")
	defer pubsub.Close()

	slog.Info("Notification hub: listening on Valkey permit_notify:post:*")

	msgCh := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			slog.Info("Notification hub: shutting down")
			return
		case msg, ok := <-msgCh:
			if !ok {
				return
			}
			var ev Event
			if err := json.Unmarshal([]byte(msg.Payload), &ev); err != nil {
				slog.Warn("Failed to decode notification event", "err", err, "payload", msg.Payload)
				continue
			}
			if ev.PostID == 0 {
				// Global notification: send to all subscribers in all posts
				slog.Info("Notification hub: broadcasting global event", "type", ev.Type)
				h.mu.RLock()
				for _, postSubs := range h.subs {
					for _, sub := range postSubs {
						select {
						case sub <- ev:
						default:
						}
					}
				}
				h.mu.RUnlock()
			} else {
				postKey := fmt.Sprintf("%d", ev.PostID)
				slog.Info("Notification hub: received event", "type", ev.Type, "permit_id", ev.ID, "post_id", ev.PostID)
				h.mu.RLock()
				for _, sub := range h.subs[postKey] {
					select {
					case sub <- ev:
					default:
					}
				}
				h.mu.RUnlock()
			}
		}
	}
}

// Subscribe returns a buffered channel for the given postID string.
func (h *Hub) Subscribe(postID string) chan Event {
	ch := make(chan Event, 8)
	h.mu.Lock()
	h.subs[postID] = append(h.subs[postID], ch)
	h.mu.Unlock()
	return ch
}

// Unsubscribe removes the channel and closes it.
func (h *Hub) Unsubscribe(postID string, ch chan Event) {
	h.mu.Lock()
	defer h.mu.Unlock()
	subs := h.subs[postID]
	for i, s := range subs {
		if s == ch {
			h.subs[postID] = append(subs[:i], subs[i+1:]...)
			close(ch)
			return
		}
	}
}

// Publish sends ev to all in-process subscribers AND to Valkey so other
// instances also receive it.
func (h *Hub) Publish(postID uint, ev Event) {
	b, err := json.Marshal(ev)
	if err != nil {
		slog.Error("Failed to marshal notification event", "err", err)
		return
	}

	// Create persistent notification record
	var msg string
	if ev.Message != "" {
		msg = ev.Message
	} else {
		switch ev.Type {
		case "new_permit":
			msg = fmt.Sprintf("Нова перепустка: %s", ev.Plate)
		case "fuzzy_match":
			msg = fmt.Sprintf("Нечіткий збіг: %s", ev.IncomingPlate)
		case "system_update":
			msg = "Системні налаштування змінено"
		default:
			msg = fmt.Sprintf("Сповіщення: %s", ev.Type)
		}
	}

	notif := models.Notification{
		PostID:  postID,
		Type:    ev.Type,
		Message: msg,
		Payload: b,
	}
	
	if repository.DB != nil {
		if err := repository.DB.Create(&notif).Error; err != nil {
			slog.Error("Failed to save notification to DB", "err", err)
		} else {
			// Update the Event struct and re-marshal so frontend gets NotificationID
			ev.NotificationID = notif.ID
			b, _ = json.Marshal(ev)
		}
	}

	channel := channelName(postID)
	if h.rdb != nil {
		if err := h.rdb.Publish(context.Background(), channel, string(b)).Err(); err != nil {
			slog.Error("Failed to publish to Valkey", "channel", channel, "err", err)
		} else {
			slog.Info("Published notification to Valkey", "channel", channel, "permit_id", ev.ID)
		}
	}
}
