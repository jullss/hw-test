package sqlstorage

import (
	"context"
	"embed"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type Storage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Connect(ctx context.Context, dbUrl string) error {
	db, err := sqlx.ConnectContext(ctx, "postgres", dbUrl)
	if err != nil {
		return err
	}

	s.db = db

	if err := goose.Up(s.db.DB, "migrations"); err != nil {
		return err
	}

	return nil
}

var embedMigrations embed.FS

func (s *Storage) Migrate(ctx context.Context) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.UpContext(ctx, s.db.DB, "migrations")
}

func (s *Storage) Close(ctx context.Context) error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Storage) Create(ctx context.Context, event *storage.Event) error {
	query := `INSERT INTO events (id, title, description, start_time, end_time, user_id, notify_in)
	          VALUES (:id, :title, :description, :start_time, :end_time, :user_id, :notify_in)`

	_, err := s.db.NamedExecContext(ctx, query, event)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Update(ctx context.Context, id string, event *storage.Event) error {
	query := `UPDATE events SET title = :title, description = :description,
			  start_time = :start_time, end_time = :end_time, user_id = :user_id,
			  notify_in = :notify_in WHERE id = :id`

	res, err := s.db.NamedExecContext(ctx, query, event)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM events WHERE id = $1`

	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func (s *Storage) ListDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	dayStart := startOfDay(date)
	dayEnd := dayStart.AddDate(0, 0, 1)

	query := `SELECT id, title, description, start_time, end_time, user_id, notify_in
			  FROM events WHERE start_time >= $1 AND start_time < $2`
	events := make([]storage.Event, 0)

	err := s.db.SelectContext(ctx, &events, query, dayStart, dayEnd)
	if err != nil {
		return events, err
	}

	return events, nil
}

func (s *Storage) ListWeek(ctx context.Context, date time.Time) ([]storage.Event, error) {
	dayStart := startOfDay(date)
	dayEnd := dayStart.AddDate(0, 0, 7)

	query := `SELECT id, title, description, start_time, end_time, user_id, notify_in
			  FROM events WHERE start_time >= $1 AND start_time < $2`
	events := make([]storage.Event, 0)

	err := s.db.SelectContext(ctx, &events, query, dayStart, dayEnd)
	if err != nil {
		return events, err
	}

	return events, nil
}

func (s *Storage) ListMonth(ctx context.Context, date time.Time) ([]storage.Event, error) {
	dayStart := startOfDay(date)
	dayEnd := dayStart.AddDate(0, 1, 0)

	query := `SELECT id, title, description, start_time, end_time, user_id, notify_in
			  FROM events WHERE start_time >= $1 AND start_time < $2`
	events := make([]storage.Event, 0)

	err := s.db.SelectContext(ctx, &events, query, dayStart, dayEnd)
	if err != nil {
		return events, err
	}

	return events, nil
}
