package internalhttp

type CreateEventRequest struct {
	UserID    string `json:"user_id" required:"true"`
	Title     string `json:"title" required:"true"`
	Desc      string `json:"desc" required:"true"`
	StartTime string `json:"start_time" required:"true"`
	EndTime   string `json:"end_time" required:"true"`
	NotifyIn  string `json:"notify_in" required:"true"`
}

type UpdateEventRequest struct {
	ID        string `json:"id" required:"true"`
	UserID    string `json:"user_id" required:"true"`
	Title     string `json:"title" required:"true"`
	Desc      string `json:"desc" required:"true"`
	StartTime string `json:"start_time" required:"true"`
	EndTime   string `json:"end_time" required:"true"`
	NotifyIn  string `json:"notify_in" required:"true"`
}

type DeleteEventRequest struct {
	ID string `json:"id"`
}
