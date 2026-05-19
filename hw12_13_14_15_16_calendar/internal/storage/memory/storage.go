package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	events map[string]storage.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[string]storage.Event),
	}
}

func (s *Storage) Create(ctx context.Context, event *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, e := range s.events {
		if e.StartTime.Equal(event.StartTime) {
			return storage.ErrDateBusy
		}
	}

	s.events[event.ID] = *event
	return nil
}

func (s *Storage) Update(ctx context.Context, id string, event *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return storage.ErrNotFound
	}

	for _, e := range s.events {
		if e.StartTime.Equal(event.StartTime) && e.ID != id {
			return storage.ErrDateBusy
		}
	}

	s.events[id] = *event
	return nil
}

func (s *Storage) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return storage.ErrNotFound
	}

	delete(s.events, id)

	return nil
}

func (s *Storage) ListDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]storage.Event, 0)

	for _, e := range s.events {
		if e.StartTime.Year() == date.Year() &&
			e.StartTime.Month() == date.Month() &&
			e.StartTime.Day() == date.Day() {
			res = append(res, e)
		}
	}

	return res, nil
}

func (s *Storage) ListWeek(ctx context.Context, date time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]storage.Event, 0)

	end := date.AddDate(0, 0, 7)

	for _, e := range s.events {
		if (e.StartTime.Equal(date) || e.StartTime.After(date)) && e.StartTime.Before(end) {
			res = append(res, e)
		}
	}

	return res, nil
}

func (s *Storage) ListMonth(ctx context.Context, date time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]storage.Event, 0)

	end := date.AddDate(0, 1, 0)

	for _, e := range s.events {
		if (e.StartTime.Equal(date) || e.StartTime.After(date)) && e.StartTime.Before(end) {
			res = append(res, e)
		}
	}

	return res, nil
}
