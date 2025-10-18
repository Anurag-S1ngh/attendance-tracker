package service

import (
	"context"
	"errors"
	"log"
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

func (s *AttendanceService) GetAttendance(ctx context.Context, userID string) ([]db.GetAttendanceRow, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("unauthorized: %v", err)
		return []db.GetAttendanceRow{}, errors.New("unauthorized")
	}

	attendances, err := s.db.GetAttendance(ctx, userUUID)
	if err != nil {
		log.Printf("failed to get attendance: %v", err)
		return []db.GetAttendanceRow{}, errors.New("internal server error")
	}

	return attendances, nil
}

func (s *AttendanceService) MarkAttendance(ctx context.Context, userID, attended, dateStr string, eventUUID uuid.UUID) (db.Attendance, error) {
	valid := map[string]bool{
		"present":  true,
		"absent":   true,
		"canceled": true,
	}

	if _, ok := valid[attended]; !ok {
		log.Printf("invalid input: attended=%s", attended)
		return db.Attendance{}, errors.New("invalid input")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("unauthorized: %v", err)
		return db.Attendance{}, errors.New("unauthorized")
	}

	layout := "2006-01-02T15:04:05.000Z"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		log.Printf("invalid date: %v", err)
		return db.Attendance{}, errors.New("invalid date")
	}

	response, err := s.db.MarkAttendance(ctx, db.MarkAttendanceParams{
		ID:       uuid.New(),
		UserID:   userUUID,
		ID_2:     eventUUID,
		Date:     date,
		Attended: db.AttendanceStatus(attended),
	})
	if err != nil {
		log.Printf("failed to mark attendance: %v", err)
		return db.Attendance{}, errors.New("internal server error")
	}

	return response, nil
}

func (s *AttendanceService) DeleteAttendance(ctx context.Context, userID string, attendanceUUID uuid.UUID) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("unauthorized: %v", err)
		return errors.New("unauthorized")
	}

	err = s.db.DeleteAttendance(ctx, db.DeleteAttendanceParams{
		UserID: userUUID,
		ID:     attendanceUUID,
	})
	if err != nil {
		log.Printf("failed to delete attendance: %v", err)
		return errors.New("internal server error")
	}

	return nil
}