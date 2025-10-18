package service

import (
	"context"
	"errors"
	"log"

	db "github.com/Anurag-S1ngh/attendance-tracker/internal/db/generated"
	"github.com/google/uuid"
)

type EventsService struct {
	db *db.Queries
}

func NewEventService(db *db.Queries) *EventsService {
	return &EventsService{
		db: db,
	}
}

func (s *EventsService) GetAllEvents(ctx context.Context, userID string) ([]db.GetUserEventsWithAttendanceAndCountsRow, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("unauthorized: %v", err)
		return []db.GetUserEventsWithAttendanceAndCountsRow{}, errors.New("unauthorized")
	}

	events, err := s.db.GetUserEventsWithAttendanceAndCounts(ctx, userUUID)
	if err != nil {
		log.Printf("failed to get events: %v", err)
		return []db.GetUserEventsWithAttendanceAndCountsRow{}, errors.New("something went wrong. please try again")
	}

	return events, nil
}

func (s *EventsService) CreateEvent(ctx context.Context, name, userID string) (db.Event, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("unauthorized: %v", err)
		return db.Event{}, errors.New("unauthorized")
	}
	event, err := s.db.CreateEvent(ctx, db.CreateEventParams{
		ID:     uuid.New(),
		UserID: userUUID,
		Name:   name,
	})
	if err != nil {
		log.Printf("failed to create event: %v", err)
		return db.Event{}, errors.New("something went wrong. please try again")
	}

	return event, nil
}

func (s *EventsService) DeleteEvent(ctx context.Context, userID string, eventUUID uuid.UUID) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("unauthorized: %v", err)
		return errors.New("unauthorized")
	}

	// ID_1 is one for eventID
	// ID_2 is one for userID
	err = s.db.DeleteEvent(ctx, db.DeleteEventParams{
		ID:   eventUUID,
		ID_2: userUUID,
	})
	if err != nil {
		log.Printf("failed to delete event: %v", err)
		return errors.New("something went wrong. please try again")
	}

	return err
}

func (s *EventsService) GetEvent(ctx context.Context, userID string, eventUUID uuid.UUID) (db.GetEventWithAttendanceAndCountsRow, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("unauthorized: %v", err)
		return db.GetEventWithAttendanceAndCountsRow{}, errors.New("unauthorized")
	}

	attendance, err := s.db.GetEventWithAttendanceAndCounts(ctx, db.GetEventWithAttendanceAndCountsParams{
		ID:     eventUUID,
		UserID: userUUID,
	})
	if err != nil {
		log.Printf("failed to get event: %v", err)
		return db.GetEventWithAttendanceAndCountsRow{}, errors.New("something went wrong. please try again")
	}

	return attendance, nil
}