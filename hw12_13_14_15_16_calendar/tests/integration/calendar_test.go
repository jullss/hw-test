//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/queue"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CalendarSuite struct {
	suite.Suite
	calendarURL string
	rabbitURL   string
	queueName   string
	sentQueue   string
}

func TestCalendarIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode")
	}
	suite.Run(t, new(CalendarSuite))
}

func (s *CalendarSuite) SetupSuite() {
	s.calendarURL = getEnv("CALENDAR_URL", "http://calendar:8888")
	s.rabbitURL = getEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/")
	s.queueName = getEnv("RABBITMQ_QUEUE", "calendar_notifications")
	s.sentQueue = getEnv("RABBITMQ_SENT_QUEUE", "calendar_notifications_sent")

	s.waitForCalendar()
}

func (s *CalendarSuite) waitForCalendar() {
	deadline := time.Now().Add(60 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := http.Get(s.calendarURL + "/list_day?date=2026-01-01")
		if err == nil {
			resp.Body.Close()
			return
		}
		time.Sleep(time.Second)
	}
	s.T().Fatal("calendar service did not become ready in time")
}

func (s *CalendarSuite) createEvent(title, userID, startTime, endTime, notifyIn string) {
	s.T().Helper()
	body := map[string]string{
		"title":      title,
		"user_id":    userID,
		"start_time": startTime,
		"end_time":   endTime,
		"notify_in":  notifyIn,
		"desc":       "integration test event",
	}
	resp := s.doRequest("POST", "/create", body)
	defer resp.Body.Close()
	require.Equal(s.T(), http.StatusCreated, resp.StatusCode, "create event should return 201")
}

func (s *CalendarSuite) doRequest(method, path string, body interface{}) *http.Response {
	s.T().Helper()
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		require.NoError(s.T(), err)
		reqBody = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, s.calendarURL+path, reqBody)
	require.NoError(s.T(), err)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(s.T(), err)
	return resp
}

func (s *CalendarSuite) listEvents(period, date string) []map[string]interface{} {
	s.T().Helper()
	resp := s.doRequest("GET", fmt.Sprintf("/list_%s?date=%s", period, date), nil)
	defer resp.Body.Close()
	require.Equal(s.T(), http.StatusOK, resp.StatusCode)
	var events []map[string]interface{}
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&events))
	return events
}

func (s *CalendarSuite) TestCreateEvent() {
	userID := "00000000-0000-0000-0000-000000000001"
	date := "2030-03-01"
	s.createEvent("Test Create Event", userID, date+" 10:00:00", date+" 11:00:00", "1h")

	events := s.listEvents("day", date)
	found := false
	for _, e := range events {
		if e["title"] == "Test Create Event" {
			found = true
			break
		}
	}
	require.True(s.T(), found, "created event should appear in day listing")
}

func (s *CalendarSuite) TestCreateEventValidationErrors() {
	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{
			name: "invalid start_time format",
			body: map[string]string{
				"title":      "Bad Event",
				"user_id":    "00000000-0000-0000-0000-000000000002",
				"start_time": "not-a-date",
				"end_time":   "2030-04-01 11:00:00",
				"notify_in":  "1h",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid end_time format",
			body: map[string]string{
				"title":      "Bad Event",
				"user_id":    "00000000-0000-0000-0000-000000000002",
				"start_time": "2030-04-01 10:00:00",
				"end_time":   "not-a-date",
				"notify_in":  "1h",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid notify_in format",
			body: map[string]string{
				"title":      "Bad Event",
				"user_id":    "00000000-0000-0000-0000-000000000002",
				"start_time": "2030-04-01 10:00:00",
				"end_time":   "2030-04-01 11:00:00",
				"notify_in":  "not-a-duration",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json body",
			body:       nil,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			var resp *http.Response
			if tc.body == nil {
				req, err := http.NewRequest("POST", s.calendarURL+"/create", bytes.NewBufferString("{invalid json}"))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				resp, err = http.DefaultClient.Do(req)
				require.NoError(t, err)
			} else {
				data, err := json.Marshal(tc.body)
				require.NoError(t, err)
				req, err := http.NewRequest("POST", s.calendarURL+"/create", bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				resp, err = http.DefaultClient.Do(req)
				require.NoError(t, err)
			}
			defer resp.Body.Close()
			require.Equal(t, tc.wantStatus, resp.StatusCode, "expected status %d for %s", tc.wantStatus, tc.name)
		})
	}
}

func (s *CalendarSuite) TestListDayEvents() {
	userID := "00000000-0000-0000-0000-000000000003"
	date := "2031-05-15"

	s.createEvent("Day Event 1", userID, date+" 09:00:00", date+" 10:00:00", "30m")
	s.createEvent("Day Event 2", userID, date+" 14:00:00", date+" 15:00:00", "30m")
	s.createEvent("Other Day Event", userID, "2031-05-16 09:00:00", "2031-05-16 10:00:00", "30m")

	events := s.listEvents("day", date)
	titles := extractTitles(events)
	require.Contains(s.T(), titles, "Day Event 1")
	require.Contains(s.T(), titles, "Day Event 2")
	require.NotContains(s.T(), titles, "Other Day Event")
}

func (s *CalendarSuite) TestListWeekEvents() {
	userID := "00000000-0000-0000-0000-000000000004"
	weekStart := "2032-06-01"

	s.createEvent("Week Event Mon", userID, "2032-06-01 10:00:00", "2032-06-01 11:00:00", "1h")
	s.createEvent("Week Event Fri", userID, "2032-06-05 10:00:00", "2032-06-05 11:00:00", "1h")
	s.createEvent("Next Week Event", userID, "2032-06-08 10:00:00", "2032-06-08 11:00:00", "1h")

	events := s.listEvents("week", weekStart)
	titles := extractTitles(events)
	require.Contains(s.T(), titles, "Week Event Mon")
	require.Contains(s.T(), titles, "Week Event Fri")
	require.NotContains(s.T(), titles, "Next Week Event")
}

func (s *CalendarSuite) TestListMonthEvents() {
	userID := "00000000-0000-0000-0000-000000000005"
	monthStart := "2033-07-01"

	s.createEvent("Month Event July 1", userID, "2033-07-01 10:00:00", "2033-07-01 11:00:00", "1h")
	s.createEvent("Month Event July 31", userID, "2033-07-31 10:00:00", "2033-07-31 11:00:00", "1h")
	s.createEvent("August Event", userID, "2033-08-01 10:00:00", "2033-08-01 11:00:00", "1h")

	events := s.listEvents("month", monthStart)
	titles := extractTitles(events)
	require.Contains(s.T(), titles, "Month Event July 1")
	require.Contains(s.T(), titles, "Month Event July 31")
	require.NotContains(s.T(), titles, "August Event")
}

func (s *CalendarSuite) TestNotificationSent() {
	now := time.Now().UTC()
	startTime := now.Add(10 * time.Second)
	endTime := startTime.Add(time.Hour)

	userID := "00000000-0000-0000-0000-000000000099"
	body := map[string]string{
		"title":      "Notification Test Event",
		"user_id":    userID,
		"start_time": startTime.Format("2006-01-02 15:04:05"),
		"end_time":   endTime.Format("2006-01-02 15:04:05"),
		"notify_in":  "1m",
		"desc":       "notification integration test",
	}

	resp := s.doRequest("POST", "/create", body)
	defer resp.Body.Close()
	require.Equal(s.T(), http.StatusCreated, resp.StatusCode)

	sentClient, err := queue.NewRabbitClient(s.rabbitURL, s.sentQueue)
	require.NoError(s.T(), err)
	defer sentClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	msgs, err := sentClient.Consume(ctx)
	require.NoError(s.T(), err)

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				s.T().Fatal("sent queue channel closed before receiving notification")
			}
			if msg.Title == "Notification Test Event" && msg.UserID == userID {
				return
			}
		case <-ctx.Done():
			s.T().Fatal("timed out waiting for notification in sent queue")
		}
	}
}

func (s *CalendarSuite) TestUpdateEvent() {
	userID := "00000000-0000-0000-0000-000000000006"
	date := "2034-08-10"

	body := map[string]string{
		"title":      "Original Title",
		"user_id":    userID,
		"start_time": date + " 10:00:00",
		"end_time":   date + " 11:00:00",
		"notify_in":  "30m",
		"desc":       "original description",
	}
	createResp := s.doRequest("POST", "/create", body)
	defer createResp.Body.Close()
	require.Equal(s.T(), http.StatusCreated, createResp.StatusCode)

	events := s.listEvents("day", date)
	var eventID string
	for _, e := range events {
		if e["title"] == "Original Title" {
			eventID = e["id"].(string)
			break
		}
	}
	require.NotEmpty(s.T(), eventID, "should find the created event")

	updateBody := map[string]string{
		"id":         eventID,
		"title":      "Updated Title",
		"user_id":    userID,
		"start_time": date + " 12:00:00",
		"end_time":   date + " 13:00:00",
		"notify_in":  "1h",
		"desc":       "updated description",
	}
	updateResp := s.doRequest("PATCH", "/update", updateBody)
	defer updateResp.Body.Close()
	require.Equal(s.T(), http.StatusOK, updateResp.StatusCode)

	updatedEvents := s.listEvents("day", date)
	found := false
	for _, e := range updatedEvents {
		if e["id"] == eventID {
			require.Equal(s.T(), "Updated Title", e["title"])
			found = true
			break
		}
	}
	require.True(s.T(), found, "updated event should still appear in day listing")
}

func (s *CalendarSuite) TestDeleteEvent() {
	userID := "00000000-0000-0000-0000-000000000007"
	date := "2035-09-20"

	body := map[string]string{
		"title":      "Event To Delete",
		"user_id":    userID,
		"start_time": date + " 10:00:00",
		"end_time":   date + " 11:00:00",
		"notify_in":  "30m",
		"desc":       "will be deleted",
	}
	createResp := s.doRequest("POST", "/create", body)
	defer createResp.Body.Close()
	require.Equal(s.T(), http.StatusCreated, createResp.StatusCode)

	events := s.listEvents("day", date)
	var eventID string
	for _, e := range events {
		if e["title"] == "Event To Delete" {
			eventID = e["id"].(string)
			break
		}
	}
	require.NotEmpty(s.T(), eventID, "should find the created event")

	deleteBody := map[string]string{"id": eventID}
	deleteResp := s.doRequest("DELETE", "/destroy", deleteBody)
	defer deleteResp.Body.Close()
	require.Equal(s.T(), http.StatusOK, deleteResp.StatusCode)

	remainingEvents := s.listEvents("day", date)
	for _, e := range remainingEvents {
		require.NotEqual(s.T(), eventID, e["id"], "deleted event should not appear in listing")
	}
}

func extractTitles(events []map[string]interface{}) []string {
	titles := make([]string, 0, len(events))
	for _, e := range events {
		if t, ok := e["title"].(string); ok {
			titles = append(titles, t)
		}
	}
	return titles
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
