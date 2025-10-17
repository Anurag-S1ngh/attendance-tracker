package service

import (
	"context"
	"errors"
	"time"

	db "github.com/Anurag-S1ngh/attendance-tracker/internal/db/generated"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AttendanceService struct {
	db     *db.Queries
	logger *zap.Logger
}

func NewAttendanceService(db *db.Queries, logger *zap.Logger) *AttendanceService {
	return &AttendanceService{
		db:     db,
		logger: logger,
	}
}

func (s *AttendanceService) GetAttendance(ctx context.Context, userID string) ([]db.GetAttendanceRow, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.Error("unauthorized", zap.Error(err))
		return []db.GetAttendanceRow{}, errors.New("unauthorized")
	}

	attendances, err := s.db.GetAttendance(ctx, userUUID)
	if err != nil {
		s.logger.Error("failed to get attendance", zap.Error(err))
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
		s.logger.Error("invalid input", zap.String("attended", attended))
		return db.Attendance{}, errors.New("invalid input")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.Error("unauthorized", zap.Error(err))
		return db.Attendance{}, errors.New("unauthorized")
	}

	layout := "2006-01-02T15:04:05.000Z"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		s.logger.Error("invalid date", zap.Error(err))
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
		s.logger.Error("failed to mark attendance", zap.Error(err))
		return db.Attendance{}, errors.New("internal server error")
	}

	return response, nil
}

func (s *AttendanceService) DeleteAttendance(ctx context.Context, userID string, attendanceUUID uuid.UUID) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.Error("unauthorized", zap.Error(err))
		return errors.New("unauthorized")
	}

	err = s.db.DeleteAttendance(ctx, db.DeleteAttendanceParams{
		UserID: userUUID,
		ID:     attendanceUUID,
	})
	if err != nil {
		s.logger.Error("failed to delete attendance", zap.Error(err))
		return errors.New("internal server error")
	}

	return nil
}
