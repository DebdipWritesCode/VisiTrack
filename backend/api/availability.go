package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/DebdipWritesCode/VisitorManagementSystem/db/sqlc"
	"github.com/gin-gonic/gin"
)

// CreateAvailabilitySlotRequest struct to bind request for creating availability slot
type createAvailabilitySlotRequest struct {
	UserID    int64  `json:"user_id" binding:"required,min=1"`
	DayOfWeek int32  `json:"day_of_week" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

// Helper function to parse time from string
func parseTime(timeStr string) time.Time {
	parsedTime, err := time.Parse("15:04", timeStr) // Adjust format as needed
	if err != nil {
		panic("invalid time format, expected HH:mm")
	}
	return parsedTime
}

// Function to create an availability slot
func (server *Server) createAvailabilitySlot(ctx *gin.Context) {
	var req createAvailabilitySlotRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAvailabilitySlotParams{
		UserID:    int32(req.UserID),
		DayOfWeek: req.DayOfWeek,
		StartTime: parseTime(req.StartTime),
		EndTime:   parseTime(req.EndTime),
	}

	availabilitySlot, err := server.store.CreateAvailabilitySlot(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, availabilitySlot)
}

// GetAvailabilityByUserRequest struct to bind request for getting availability by user ID
type getAvailabilityByUserRequest struct {
	UserID int64 `uri:"user_id" binding:"required,min=1"`
}

// Function to get availability slots for a user
func (server *Server) getAvailabilityByUser(ctx *gin.Context) {
	var req getAvailabilityByUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	availabilitySlots, err := server.store.GetAvailabilityByUser(ctx, int32(req.UserID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, availabilitySlots)
}

// DeleteAvailabilitySlotRequest struct to bind request for deleting an availability slot
type deleteAvailabilitySlotRequest struct {
	UserID    int64  `json:"user_id" binding:"required,min=1"`
	DayOfWeek int32  `json:"day_of_week" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

// Function to delete an availability slot
func (server *Server) deleteAvailabilitySlot(ctx *gin.Context) {
	var req deleteAvailabilitySlotRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.DeleteAvailabilitySlotParams{
		UserID:    int32(req.UserID),
		DayOfWeek: req.DayOfWeek,
		StartTime: parseTime(req.StartTime),
		EndTime:   parseTime(req.EndTime),
	}
	err := server.store.DeleteAvailabilitySlot(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "availability slot deleted"})
}

// DeleteAvailabilityByUserRequest struct to bind request for deleting all availability slots by user ID
type deleteAvailabilityByUserRequest struct {
	UserID int64 `uri:"user_id" binding:"required,min=1"`
}

// Function to delete all availability slots by user ID
func (server *Server) deleteAvailabilityByUser(ctx *gin.Context) {
	var req deleteAvailabilityByUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteAvailabilityByUser(ctx, int32(req.UserID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "all availability slots deleted for user"})
}

// UpdateAvailabilityStatusRequest struct to bind request for updating availability status
type updateAvailabilityStatusRequest struct {
	UserID    int64  `json:"user_id" binding:"required,min=1"`
	DayOfWeek int32  `json:"day_of_week" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
	Status    string `json:"status" binding:"required,oneof=available not_available"`
}

// Function to update availability status
func (server *Server) updateAvailabilityStatus(ctx *gin.Context) {
	var req updateAvailabilityStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateAvailabilityStatusParams{
		UserID:    int32(req.UserID),
		DayOfWeek: req.DayOfWeek,
		StartTime: parseTime(req.StartTime),
		EndTime:   parseTime(req.EndTime),
		Status:    sql.NullString{String: req.Status, Valid: req.Status != ""},
	}

	err := server.store.UpdateAvailabilityStatus(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "availability status updated"})
}
