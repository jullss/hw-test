package storage

import "time"

type Event struct {
	ID        string        `db:"id"`
	UserID    string        `db:"user_id"`
	Title     string        `db:"title"`
	Desc      string        `db:"description"`
	StartTime time.Time     `db:"start_time"`
	EndTime   time.Time     `db:"end_time"`
	NotifyIn  time.Duration `db:"notify_in"`
}
