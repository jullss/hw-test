package memorystorage

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	s := New()
	ctx := context.Background()
	wg := sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			id := fmt.Sprintf("%d", n)
			_ = s.Create(ctx, &storage.Event{
				ID:        id,
				StartTime: time.Now().Add(time.Duration(n) * time.Second),
			})
			_, _ = s.ListDay(ctx, time.Now())
		}(i)
	}
	wg.Wait()
}

func TestStorage_Create(t *testing.T) {
	s := New()
	ctx := context.Background()
	startTime := time.Now()

	e1 := &storage.Event{ID: "1", StartTime: startTime, Title: "Event 1"}

	err := s.Create(ctx, e1)
	assert.NoError(t, err)

	e2 := &storage.Event{ID: "2", StartTime: startTime, Title: "Event 2"}
	err = s.Create(ctx, e2)
	assert.ErrorIs(t, err, storage.ErrDateBusy)
}

func TestStorage_UpdateDelete_Errors(t *testing.T) {
	s := New()
	ctx := context.Background()

	err := s.Update(ctx, "not_existed", &storage.Event{Title: "New"})
	assert.ErrorIs(t, err, storage.ErrNotFound)

	err = s.Delete(ctx, "not_existed")
	assert.ErrorIs(t, err, storage.ErrNotFound)
}

func TestStorage_ConcurrencyAndBusinessErrors(t *testing.T) {
	s := New()
	ctx := context.Background()
	baseTime := time.Date(2026, 5, 17, 12, 0, 0, 0, time.UTC)

	var wg sync.WaitGroup
	workers := 100

	errChan := make(chan error, workers*2)

	for i := 0; i < workers; i++ {
		wg.Add(2)

		go func(n int) {
			defer wg.Done()
			id := fmt.Sprintf("unique-%d", n)
			eventTime := baseTime.Add(time.Duration(n) * time.Hour)

			err := s.Create(ctx, &storage.Event{
				ID:        id,
				StartTime: eventTime,
			})
			if err != nil {
				errChan <- fmt.Errorf("unexpected error for unique event %d: %w", n, err)
			}
		}(i)

		go func(n int) {
			defer wg.Done()
			id := fmt.Sprintf("conflict-%d", n)

			err := s.Create(ctx, &storage.Event{
				ID:        id,
				StartTime: baseTime,
			})
			if err != nil && !errors.Is(err, storage.ErrDateBusy) {
				errChan <- fmt.Errorf("unexpected error for conflicting event %d: %w", n, err)
			}
		}(i)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		t.Errorf("concurrency error: %v", err)
	}

	events, err := s.ListDay(ctx, baseTime)
	if err != nil {
		t.Fatalf("failed to list day: %v", err)
	}

	var conflictCount int
	for _, e := range events {
		if e.StartTime.Equal(baseTime) {
			conflictCount++
		}
	}

	if conflictCount != 1 {
		t.Errorf("expected exactly 1 event for baseTime due to ErrDateBusy, found %d", conflictCount)
	}
}
