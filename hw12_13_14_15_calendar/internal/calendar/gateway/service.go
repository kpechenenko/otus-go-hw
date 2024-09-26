package gateway

import (
	"context"
	"time"

	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/calendar/model"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/calendar/repository"
	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/calendar/service"
	desc "github.com/kpechenenko/hw12_13_14_15_calendar/calendar/pkg/api/calendar"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	desc.UnimplementedCalendarServer
	srv service.Service
}

func New(srv service.Service) *Service {
	return &Service{srv: srv}
}

func (s *Service) CreateEvent(ctx context.Context, r *desc.AddEventRequest) (*desc.AddEventResponse, error) {
	event := repository.AddEventParams{
		Title:       r.GetTitle(),
		Date:        r.GetDate().AsTime(),
		Duration:    time.Duration(r.GetDuration()),
		Description: nil,
		OwnerID:     model.UserID(r.GetOwnerId()),
		NotifyFor:   nil,
	}
	if r.Description != nil {
		d := r.GetDescription()
		event.Description = &d
	}
	if r.NotifyFor != nil {
		d := time.Duration(r.GetNotifyFor())
		event.NotifyFor = &d
	}
	eventID, err := s.srv.CreateEvent(ctx, event)
	if err != nil {
		return nil, err
	}
	return &desc.AddEventResponse{Id: eventID.String()}, nil
}
func (s *Service) UpdateEvent(ctx context.Context, r *desc.UpdateEventRequest) (*emptypb.Empty, error) {
	resp := &emptypb.Empty{}
	eventID, err := model.ParseEventID(r.GetId())
	if err != nil {
		return resp, err
	}
	event := repository.UpdateEventParams{
		ID:          eventID,
		Title:       r.GetTitle(),
		Date:        r.GetDate().AsTime(),
		Duration:    time.Duration(r.GetDuration()),
		Description: nil,
		NotifyFor:   nil,
	}
	if r.Description != nil {
		d := r.GetDescription()
		event.Description = &d
	}
	if r.NotifyFor != nil {
		d := time.Duration(r.GetNotifyFor())
		event.NotifyFor = &d
	}
	err = s.srv.UpdateEvent(ctx, event)
	return resp, err
}
func (s *Service) DeleteEvent(ctx context.Context, r *desc.DeleteEventRequest) (*emptypb.Empty, error) {
	resp := &emptypb.Empty{}
	eventID, err := model.ParseEventID(r.GetId())
	if err != nil {
		return resp, err
	}
	err = s.srv.DeleteEvent(ctx, eventID)
	return resp, err
}
func (s *Service) GetEventsForDay(ctx context.Context, r *desc.GetEventsForDayRequest) (*desc.GetEventsResponse, error) {
	events, err := s.srv.GetEventsForDay(ctx, r.GetDay().AsTime())
	if err != nil {
		return nil, err
	}
	resp := &desc.GetEventsResponse{
		Items: s.convertToPbResp(events),
	}
	return resp, nil
}

func (s *Service) GetEventsForWeek(ctx context.Context, r *desc.GetEventsForWeekRequest) (*desc.GetEventsResponse, error) {
	events, err := s.srv.GetEventsForWeek(ctx, r.GetBeginDate().AsTime())
	if err != nil {
		return nil, err
	}
	resp := &desc.GetEventsResponse{
		Items: s.convertToPbResp(events),
	}
	return resp, nil
}

func (s *Service) GetEventsForMonth(ctx context.Context, r *desc.GetEventsForMonthRequest) (*desc.GetEventsResponse, error) {
	events, err := s.srv.GetEventsForMonth(ctx, r.GetBeginDate().AsTime())
	if err != nil {
		return nil, err
	}
	resp := &desc.GetEventsResponse{
		Items: s.convertToPbResp(events),
	}
	return resp, nil
}

func (s *Service) convertToPbResp(events []model.Event) []*desc.Event {
	resp := make([]*desc.Event, len(events))
	for i, event := range events {
		resp[i] = &desc.Event{
			Id:          event.ID.String(),
			Title:       event.Title,
			Date:        timestamppb.New(event.Date),
			Duration:    int64(event.Duration),
			Description: event.Description,
			OwnerId:     int64(event.OwnerID),
			NotifyFor:   nil,
		}
		if event.NotifyFor != nil {
			nf := int64(*event.NotifyFor)
			resp[i].NotifyFor = &nf
		}
	}
	return resp
}
