package app

import (
	"context"
	"time"

	"github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

type Storage interface {
	Create(ctx context.Context, event *storage.Event) error
	Update(ctx context.Context, id string, event *storage.Event) error
	Delete(ctx context.Context, id string) error
	ListDay(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListWeek(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListMonth(ctx context.Context, date time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, event *storage.Event) error {
	a.logger.Info("creating event", "id", event.ID, "title", event.Title)

	err := a.storage.Create(ctx, event)
	if err != nil {
		a.logger.Error("failed to create event", "id", event.ID, "err", err)
		return err
	}

	a.logger.Info("event created successfully", "id", event.ID)
	return nil
}

func (a *App) UpdateEvent(ctx context.Context, event *storage.Event) error {
	a.logger.Info("updating event", "id", event.ID, "title", event.Title)

	err := a.storage.Update(ctx, event.ID, event)
	if err != nil {
		a.logger.Error("failed to update event", "id", event.ID, "err", err)
		return err
	}

	a.logger.Info("event updated successfully", "id", event.ID)
	return nil
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	a.logger.Info("deleting event", "id", id)

	err := a.storage.Delete(ctx, id)
	if err != nil {
		a.logger.Error("failed to delete event", "id", id, "err", err)
		return err
	}

	a.logger.Info("event deleted successfully", "id", id)
	return nil
}

func (a *App) ListDayEvent(ctx context.Context, date time.Time) ([]storage.Event, error) {
	a.logger.Info("list day event", "date", date)

	events, err := a.storage.ListDay(ctx, date)
	if err != nil {
		a.logger.Error("failed to list day event", "date", date, "err", err)
		return []storage.Event{}, err
	}

	a.logger.Info("event list day successfully")
	return events, nil
}

func (a *App) ListWeekEvent(ctx context.Context, date time.Time) ([]storage.Event, error) {
	a.logger.Info("list week event", "date", date)

	events, err := a.storage.ListWeek(ctx, date)
	if err != nil {
		a.logger.Error("failed to list week event", "date", date, "err", err)
		return []storage.Event{}, err
	}

	a.logger.Info("event list week successfully")
	return events, nil
}

func (a *App) ListMonthEvent(ctx context.Context, date time.Time) ([]storage.Event, error) {
	a.logger.Info("list month event", "date", date)

	events, err := a.storage.ListMonth(ctx, date)
	if err != nil {
		a.logger.Error("failed to list month event", "date", date, "err", err)
		return []storage.Event{}, err
	}

	a.logger.Info("event list month successfully")
	return events, nil
}
