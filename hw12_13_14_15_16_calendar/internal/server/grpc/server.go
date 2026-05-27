package internalgrpc

import (
	"context"
	"time"

	"github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/pb"
	"github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedCalendarServiceServer
	log Logger
	app Application
}

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

func NewServer(logger Logger, app Application) *Server {
	return &Server{
		log: logger,
		app: app,
	}
}

type Application interface {
	CreateEvent(ctx context.Context, event *storage.Event) error
	UpdateEvent(ctx context.Context, event *storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	ListDayEvent(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListWeekEvent(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListMonthEvent(ctx context.Context, date time.Time) ([]storage.Event, error)
}

func (s *Server) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	s.log.Info("gRPC request: CreateEvent", "user_id", req.Event.UserId, "title", req.Event.Title)

	if req.Event == nil {
		return nil, status.Error(codes.InvalidArgument, "event data is required")
	}

	startTime := req.Event.StartTime.AsTime()
	endTime := req.Event.EndTime.AsTime()
	notifyIn := req.Event.NotifyIn.AsDuration()

	event := storage.Event{
		UserID:    req.Event.UserId,
		Title:     req.Event.Title,
		Desc:      req.Event.Desc,
		StartTime: startTime,
		EndTime:   endTime,
		NotifyIn:  notifyIn,
	}

	if err := s.app.CreateEvent(ctx, &event); err != nil {
		s.log.Error("gRPC CreateEvent: business logic failed", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.CreateEventResponse{
		Id: event.ID,
	}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	s.log.Info("gRPC request: UpdateEvent", "user_id", req.Event.UserId, "title", req.Event.Title)

	if req.Event == nil {
		return nil, status.Error(codes.InvalidArgument, "event data is required")
	}

	startTime := req.Event.StartTime.AsTime()
	endTime := req.Event.EndTime.AsTime()
	notifyIn := req.Event.NotifyIn.AsDuration()

	event := storage.Event{
		ID:        req.Id,
		UserID:    req.Event.UserId,
		Title:     req.Event.Title,
		Desc:      req.Event.Desc,
		StartTime: startTime,
		EndTime:   endTime,
		NotifyIn:  notifyIn,
	}

	if err := s.app.UpdateEvent(ctx, &event); err != nil {
		s.log.Error("gRPC UpdateEvent: business logic failed", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.UpdateEventResponse{
		Id: event.ID,
	}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	s.log.Info("gRPC request: DeleteEvent", "id", req.Id)

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.app.DeleteEvent(ctx, req.Id); err != nil {
		s.log.Error("gRPC DeleteEvent: business logic failed", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.DeleteEventResponse{
		Id: req.Id,
	}, nil
}

func (s *Server) ListDayEvent(ctx context.Context, req *pb.ListDayRequest) (*pb.ListDayResponse, error) {
	s.log.Info("gRPC request: ListDayEvent", "date", req.Date)

	date := req.Date.AsTime()
	var pbEvents []*pb.Event

	events, err := s.app.ListDayEvent(ctx, date)
	if err != nil {
		s.log.Error("gRPC ListDayEvent: business logic failed", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	for _, e := range events {
		pbEvents = append(pbEvents, &pb.Event{
			Id:        e.ID,
			UserId:    e.UserID,
			Title:     e.Title,
			Desc:      e.Desc,
			StartTime: timestamppb.New(e.StartTime),
			EndTime:   timestamppb.New(e.EndTime),
			NotifyIn:  durationpb.New(e.NotifyIn),
		})
	}

	return &pb.ListDayResponse{
		Events: pbEvents,
	}, nil
}

func (s *Server) ListWeekEvent(ctx context.Context, req *pb.ListWeekRequest) (*pb.ListWeekResponse, error) {
	s.log.Info("gRPC request: ListWeekEvent", "date", req.Date)

	date := req.Date.AsTime()
	var pbEvents []*pb.Event

	events, err := s.app.ListWeekEvent(ctx, date)
	if err != nil {
		s.log.Error("gRPC ListWeekEvent: business logic failed", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	for _, e := range events {
		pbEvents = append(pbEvents, &pb.Event{
			Id:        e.ID,
			UserId:    e.UserID,
			Title:     e.Title,
			Desc:      e.Desc,
			StartTime: timestamppb.New(e.StartTime),
			EndTime:   timestamppb.New(e.EndTime),
			NotifyIn:  durationpb.New(e.NotifyIn),
		})
	}

	return &pb.ListWeekResponse{
		Events: pbEvents,
	}, nil
}

func (s *Server) ListMonthEvent(ctx context.Context, req *pb.ListMonthRequest) (*pb.ListMonthResponse, error) {
	s.log.Info("gRPC request: ListMonthEvent", "date", req.Date)

	date := req.Date.AsTime()
	var pbEvents []*pb.Event

	events, err := s.app.ListMonthEvent(ctx, date)
	if err != nil {
		s.log.Error("gRPC ListMonthEvent: business logic failed", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	for _, e := range events {
		pbEvents = append(pbEvents, &pb.Event{
			Id:        e.ID,
			UserId:    e.UserID,
			Title:     e.Title,
			Desc:      e.Desc,
			StartTime: timestamppb.New(e.StartTime),
			EndTime:   timestamppb.New(e.EndTime),
			NotifyIn:  durationpb.New(e.NotifyIn),
		})
	}

	return &pb.ListMonthResponse{
		Events: pbEvents,
	}, nil
}
