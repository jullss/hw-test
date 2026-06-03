package internalgrpc

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/pb"
	"github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LoggerMock struct{}

func (l *LoggerMock) Info(msg string, args ...any)  {}
func (l *LoggerMock) Error(msg string, args ...any) {}

type ApplicationMock struct {
	mock.Mock
}

func (m *ApplicationMock) CreateEvent(ctx context.Context, event *storage.Event) error {
	args := m.Called(ctx, event)
	if args.Get(0) == nil {
		event.ID = "uuid-123"
		return nil
	}
	return args.Error(0)
}

func (m *ApplicationMock) UpdateEvent(ctx context.Context, event *storage.Event) error { return nil }
func (m *ApplicationMock) DeleteEvent(ctx context.Context, id string) error            { return nil }
func (m *ApplicationMock) ListDayEvent(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return nil, nil
}

func (m *ApplicationMock) ListWeekEvent(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return nil, nil
}

func (m *ApplicationMock) ListMonthEvent(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return nil, nil
}

func dialer(lis *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(context.Context, string) (net.Conn, error) { return lis.Dial() }
}

func TestCreateEvent_Success(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	grpcServer := grpc.NewServer()

	appMock := new(ApplicationMock)
	loggerMock := new(LoggerMock)
	server := NewServer(loggerMock, appMock)

	pb.RegisterCalendarServiceServer(grpcServer, server)
	go func() { _ = grpcServer.Serve(lis) }()
	defer grpcServer.Stop()

	ctx := context.Background()
	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(dialer(lis)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err)
	defer conn.Close()

	client := pb.NewCalendarServiceClient(conn)
	appMock.On("CreateEvent", mock.Anything, mock.Anything).Return(nil)

	now := time.Now().Truncate(time.Second)
	resp, err := client.CreateEvent(ctx, &pb.CreateEventRequest{
		Event: &pb.Event{
			UserId:    "user_1",
			Title:     "Meeting",
			StartTime: timestamppb.New(now),
			EndTime:   timestamppb.New(now.Add(time.Hour)),
			NotifyIn:  durationpb.New(15 * time.Minute),
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "uuid-123", resp.Id)
}

func TestCreateEvent_NilEvent(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	grpcServer := grpc.NewServer()
	appMock := new(ApplicationMock)
	loggerMock := new(LoggerMock)
	server := NewServer(loggerMock, appMock)
	pb.RegisterCalendarServiceServer(grpcServer, server)
	go func() { _ = grpcServer.Serve(lis) }()
	defer grpcServer.Stop()

	ctx := context.Background()
	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(dialer(lis)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err)
	defer conn.Close()

	client := pb.NewCalendarServiceClient(conn)

	resp, err := client.CreateEvent(ctx, &pb.CreateEventRequest{Event: nil})

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Equal(t, "event data is required", st.Message())
}

func TestUpdateEvent_Success(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	grpcServer := grpc.NewServer()
	appMock := new(ApplicationMock)
	loggerMock := new(LoggerMock)
	server := NewServer(loggerMock, appMock)
	pb.RegisterCalendarServiceServer(grpcServer, server)
	go func() { _ = grpcServer.Serve(lis) }()
	defer grpcServer.Stop()

	ctx := context.Background()
	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(dialer(lis)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err)
	defer conn.Close()

	client := pb.NewCalendarServiceClient(conn)

	now := time.Now().Truncate(time.Second)
	resp, err := client.UpdateEvent(ctx, &pb.UpdateEventRequest{
		Id: "id-to-update",
		Event: &pb.Event{
			UserId:    "user_1",
			Title:     "New Title",
			StartTime: timestamppb.New(now),
			EndTime:   timestamppb.New(now.Add(time.Hour)),
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "id-to-update", resp.Id)
}

func TestDeleteEvent_Success(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	grpcServer := grpc.NewServer()
	appMock := new(ApplicationMock)
	loggerMock := new(LoggerMock)
	server := NewServer(loggerMock, appMock)
	pb.RegisterCalendarServiceServer(grpcServer, server)
	go func() { _ = grpcServer.Serve(lis) }()
	defer grpcServer.Stop()

	ctx := context.Background()
	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(dialer(lis)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err)
	defer conn.Close()

	client := pb.NewCalendarServiceClient(conn)

	resp, err := client.DeleteEvent(ctx, &pb.DeleteEventRequest{Id: "id-to-delete"})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "id-to-delete", resp.Id)
}

func TestDeleteEvent_EmptyID(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	grpcServer := grpc.NewServer()
	appMock := new(ApplicationMock)
	loggerMock := new(LoggerMock)
	server := NewServer(loggerMock, appMock)
	pb.RegisterCalendarServiceServer(grpcServer, server)
	go func() { _ = grpcServer.Serve(lis) }()
	defer grpcServer.Stop()

	ctx := context.Background()
	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(dialer(lis)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err)
	defer conn.Close()

	client := pb.NewCalendarServiceClient(conn)

	resp, err := client.DeleteEvent(ctx, &pb.DeleteEventRequest{Id: ""})

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestListDayEvent_Success(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	grpcServer := grpc.NewServer()
	appMock := new(ApplicationMock)
	loggerMock := new(LoggerMock)
	server := NewServer(loggerMock, appMock)
	pb.RegisterCalendarServiceServer(grpcServer, server)
	go func() { _ = grpcServer.Serve(lis) }()
	defer grpcServer.Stop()

	ctx := context.Background()
	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(dialer(lis)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err)
	defer conn.Close()

	client := pb.NewCalendarServiceClient(conn)

	targetDate := time.Now().Truncate(time.Second)
	resp, err := client.ListDayEvent(ctx, &pb.ListDayRequest{Date: timestamppb.New(targetDate)})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
