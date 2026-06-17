package storage

import "time"

type Event struct {
	ID        string        `db:"id"          json:"id"`
	UserID    string        `db:"user_id"      json:"user_id"`
	Title     string        `db:"title"        json:"title"`
	Desc      string        `db:"description"  json:"desc"`
	StartTime time.Time     `db:"start_time"   json:"start_time"`
	EndTime   time.Time     `db:"end_time"     json:"end_time"`
	NotifyIn  time.Duration `db:"notify_in"    json:"notify_in"`
}
