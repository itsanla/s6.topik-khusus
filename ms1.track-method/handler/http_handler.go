package handler

import (
	"net/http"

	"track-method/domain"

	"github.com/gin-gonic/gin"
)

type TrackHandler struct {
	usecase domain.TrackUsecase
}

func NewTrackHandler(r *gin.Engine, uc domain.TrackUsecase) {
	h := &TrackHandler{usecase: uc}

	api := r.Group("/api/v1")
	{
		api.POST("/track", h.RecordEvent)
		api.GET("/track", h.GetAllStats)
		api.GET("/track/:event", h.GetEventStats)
		api.DELETE("/track/:event", h.ResetEvent)
	}
}

type recordEventRequest struct {
	Event string `json:"event" binding:"required"`
}

func (h *TrackHandler) RecordEvent(c *gin.Context) {
	var req recordEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "request tidak valid: field 'event' wajib diisi",
		})
		return
	}

	result, err := h.usecase.RecordEvent(c.Request.Context(), req.Event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "event berhasil dicatat",
		"data":    result,
	})
}

func (h *TrackHandler) GetEventStats(c *gin.Context) {
	eventName := c.Param("event")

	result, err := h.usecase.GetEventStats(c.Request.Context(), eventName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

func (h *TrackHandler) GetAllStats(c *gin.Context) {
	results, err := h.usecase.GetAllStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"total":   len(results),
		"data":    results,
	})
}

func (h *TrackHandler) ResetEvent(c *gin.Context) {
	eventName := c.Param("event")

	err := h.usecase.ResetEventStats(c.Request.Context(), eventName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "statistik event berhasil direset",
	})
}
