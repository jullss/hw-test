package internalhttp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type LoggerMock struct{}

func (l *LoggerMock) Info(msg string, args ...any)  {}
func (l *LoggerMock) Error(msg string, args ...any) {}

type ApplicationMock struct {
	mock.Mock
}

func (m *ApplicationMock) CreateEvent(ctx context.Context, event *storage.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *ApplicationMock) UpdateEvent(ctx context.Context, event *storage.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *ApplicationMock) DeleteEvent(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ApplicationMock) ListDayEvent(ctx context.Context, date time.Time) ([]storage.Event, error) {
	args := m.Called(ctx, date)
	return args.Get(0).([]storage.Event), args.Error(1)
}

func (m *ApplicationMock) ListWeekEvent(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return nil, nil
}

func (m *ApplicationMock) ListMonthEvent(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return nil, nil
}

func TestHTTP_CreateEvent_Success(t *testing.T) {
	appMock := new(ApplicationMock)
	loggerMock := new(LoggerMock)
	server := NewServer(loggerMock, appMock)

	appMock.On("CreateEvent", mock.Anything, mock.Anything).Return(nil)

	jsonBody := `{
		"user_id": "http_user_123",
		"title": "Standup",
		"desc": "Daily sync",
		"start_time": "2026-05-30 15:00:00",
		"end_time": "2026-05-30 15:30:00",
		"notify_in": "15m"
	}`

	req, err := http.NewRequest("POST", "/create", strings.NewReader(jsonBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	http.HandlerFunc(server.createEventHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.JSONEq(t, `{"status": "created"}`, rr.Body.String())
	appMock.AssertExpectations(t)
}

func TestHTTP_CreateEvent_InvalidTime(t *testing.T) {
	appMock := new(ApplicationMock)
	loggerMock := new(LoggerMock)
	server := NewServer(loggerMock, appMock)

	jsonBody := `{
		"user_id": "http_user_123",
		"title": "Standup",
		"start_time": "30.05.2026 15:00:00",
		"end_time": "2026-05-30 15:30:00",
		"notify_in": "15m"
	}`

	req, err := http.NewRequest("POST", "/create", strings.NewReader(jsonBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	http.HandlerFunc(server.createEventHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "invalid start_time format")
}

func TestHTTP_DeleteEvent_Success(t *testing.T) {
	appMock := new(ApplicationMock)
	loggerMock := new(LoggerMock)
	server := NewServer(loggerMock, appMock)

	appMock.On("DeleteEvent", mock.Anything, "event-uuid-to-delete").Return(nil)

	jsonBody := `{"id": "event-uuid-to-delete"}`

	req, err := http.NewRequest("DELETE", "/destroy", strings.NewReader(jsonBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	http.HandlerFunc(server.deleteEventHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `{"status": "deleted"}`, rr.Body.String())
	appMock.AssertExpectations(t)
}
