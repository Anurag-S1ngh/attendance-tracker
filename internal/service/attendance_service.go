package service

import (
	"context"
	"errors"
	"time"

	db "github.com/Anurag-S1ngh/attendance-tracker/internal/db/generated"
	"github.com/google/uuid"
)

type AttendanceService struct {
	db *db.Queries
}

func NewAttendanceService(db *db.Queries) *AttendanceService {
	return &AttendanceService{
		db: db,
	}
}

func (s *AttendanceService) GetAttendance(ctx context.Context, userID string) ([]db.Attendance, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return []db.Attendance{}, errors.New("invalid user id")
	}

	return s.db.GetAttendance(ctx, userUUID)
}

func (s *AttendanceService) MarkAttendance(ctx context.Context, userID, attended, dateStr string, eventUUID uuid.UUID) (db.Attendance, error) {
	valid := map[string]bool{
		"present":  true,
		"absent":   true,
		"canceled": true,
	}

	if _, ok := valid[attended]; !ok {
		return db.Attendance{}, errors.New("invalid input")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return db.Attendance{}, errors.New("invalid user id")
	}

	layout := "2006-01-02T15:04:05.000Z"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		return db.Attendance{}, errors.New("invalid date")
	}

	return s.db.MarkAttendance(ctx, db.MarkAttendanceParams{
		ID:       uuid.New(),
		UserID:   userUUID,
		ID_2:     eventUUID,
		Date:     date,
		Attended: db.AttendanceStatus(attended),
	})
}

func (s *AttendanceService) DeleteAttendance(ctx context.Context, userID string, attendanceUUID uuid.UUID) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	return s.db.DeleteAttendance(ctx, db.DeleteAttendanceParams{
		UserID: userUUID,
		ID:     attendanceUUID,
	})
}
