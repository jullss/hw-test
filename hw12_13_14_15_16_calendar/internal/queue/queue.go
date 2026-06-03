package queue

import (
	"context"
	"time"
)

type Notification struct {
	EventID string    `json:"event_id"`
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
	UserID  string    `json:"user_id"`
}

type Publisher interface {
	Publish(ctx context.Context, msg Notification) error
	Close() error
}

type Consumer interface {
	Consume(ctx context.Context) (<-chan Notification, error)
	Close() error
}
