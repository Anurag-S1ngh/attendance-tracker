package handlers

import (
	"net/http"

	"github.com/Anurag-S1ngh/attendance-tracker/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventsHandler struct {
	eventsService *service.EventsService
}

func NewEventsHandler(s *service.EventsService) *EventsHandler {
	return &EventsHandler{
		eventsService: s,
	}
}

func (h *EventsHandler) GetEvents(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "sign in first"})
		return
	}
	events, err := h.eventsService.GetEvents(c, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

func (h *EventsHandler) CreateEvent(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "sign in first"})
		return
	}

	event, err := h.eventsService.CreateEvent(c, req.Name, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": event})
}

func (h *EventsHandler) DeleteEvent(c *gin.Context) {
	eventUUID, err := uuid.Parse(c.Param("eventID"))
	if err != nil || eventUUID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "sign in first"})
		return
	}

	err = h.eventsService.DeleteEvent(c, userID.(string), eventUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "event deleted"})
}
