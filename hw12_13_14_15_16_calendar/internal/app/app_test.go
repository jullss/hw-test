package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type mockStorage struct {
	errToReturn error
}

func (m *mockStorage) Create(_ context.Context, _ *storage.Event) error {
	return m.errToReturn
}

func (m *mockStorage) Update(_ context.Context, _ string, _ *storage.Event) error { return nil }
func (m *mockStorage) Delete(_ context.Context, _ string) error                   { return nil }
func (m *mockStorage) ListDay(_ context.Context, _ time.Time) ([]storage.Event, error) {
	return nil, nil
}

func (m *mockStorage) ListWeek(_ context.Context, _ time.Time) ([]storage.Event, error) {
	return nil, nil
}

func (m *mockStorage) ListMonth(_ context.Context, _ time.Time) ([]storage.Event, error) {
	return nil, nil
}

type mockLogger struct{}

func (m *mockLogger) Info(_ string, _ ...any)  {}
func (m *mockLogger) Error(_ string, _ ...any) {}

func TestCreateEvent_Success(t *testing.T) {
	logger := &mockLogger{}
	db := &mockStorage{errToReturn: nil}

	calendar := app.New(logger, db)

	event := &storage.Event{
		ID:     "test-id",
		UserID: "user-1",
		Title:  "Trip",
	}

	err := calendar.CreateEvent(context.Background(), event)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCreateEvent_StorageError(t *testing.T) {
	logger := &mockLogger{}

	expectedErr := errors.New("db connection timeout")
	db := &mockStorage{errToReturn: expectedErr}

	calendar := app.New(logger, db)
	event := &storage.Event{ID: "1", UserID: "user-1", Title: "Trip"}

	err := calendar.CreateEvent(context.Background(), event)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

func TestUpdateEvent_Success(t *testing.T) {
	logger := &mockLogger{}
	db := &mockStorage{errToReturn: nil}

	calendar := app.New(logger, db)
	event := &storage.Event{ID: "1", UserID: "user-1", Title: "Updated Trip"}

	err := calendar.UpdateEvent(context.Background(), event)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
