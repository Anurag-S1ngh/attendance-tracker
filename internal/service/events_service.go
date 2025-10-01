package service

import (
	"context"
	"errors"

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

func (s *EventsService) GetEvents(ctx context.Context, userID string) ([]db.GetUserEventsWithAttendanceAndCountsRow, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return []db.GetUserEventsWithAttendanceAndCountsRow{}, errors.New("invalid user_id")
	}

	return s.db.GetUserEventsWithAttendanceAndCounts(ctx, userUUID)
}

func (s *EventsService) CreateEvent(ctx context.Context, name, userID string) (db.Event, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return db.Event{}, errors.New("invalid user_id")
	}
	return s.db.CreateEvent(ctx, db.CreateEventParams{
		ID:     uuid.New(),
		UserID: userUUID,
		Name:   name,
	})
}

func (s *EventsService) DeleteEvent(ctx context.Context, userID string, eventUUID uuid.UUID) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user_id")
	}

	// ID_1 is one for eventID
	// ID_2 is one for userID
	return s.db.DeleteEvent(ctx, db.DeleteEventParams{
		ID:   eventUUID,
		ID_2: userUUID,
	})
}
