package storage

import (
	"context"
	"errors"
	"time"
)

type Repository interface {
	Create(ctx context.Context, event *Event) error
	Update(ctx context.Context, id string, event *Event) error
	Delete(ctx context.Context, id string) error
	ListDay(ctx context.Context, date time.Time) ([]Event, error)
	ListWeek(ctx context.Context, date time.Time) ([]Event, error)
	ListMonth(ctx context.Context, date time.Time) ([]Event, error)
}

var (
	ErrDateBusy = errors.New("date is already taken")
	ErrNotFound = errors.New("event is not found")
)
