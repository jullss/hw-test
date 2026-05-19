package storage

import "time"

type Notification struct {
	EventID string
	Title   string
	Date    time.Time
	UserID  string
}
