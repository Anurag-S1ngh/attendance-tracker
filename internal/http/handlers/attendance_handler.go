package handlers

import (
	"net/http"

	"github.com/Anurag-S1ngh/attendance-tracker/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AttendanceHandler struct {
	attendanceService *service.AttendanceService
}

func NewAttendanceHandler(s *service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{attendanceService: s}
}

func (h *AttendanceHandler) GetAttendance(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "sign in first"})
		return
	}
	attendance, err := h.attendanceService.GetAttendance(c, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"attendance": attendance})
}

func (h *AttendanceHandler) MarkAttendance(c *gin.Context) {
	var req struct {
		Attended string `json:"attended"`
		Date     string `json:"date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	eventUUID, err := uuid.Parse(c.Param("eventID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "sign in first"})
		return
	}

	attendance, err := h.attendanceService.MarkAttendance(c, userID.(string), req.Attended, req.Date, eventUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"attendance": attendance})
}

func (h *AttendanceHandler) DeleteAttendance(c *gin.Context) {
	attendanceUUID, err := uuid.Parse(c.Param("attendanceID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "sign in first"})
		return
	}

	err = h.attendanceService.DeleteAttendance(c, userID.(string), attendanceUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "attendance deleted"})
}
