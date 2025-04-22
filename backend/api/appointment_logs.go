package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/DebdipWritesCode/VisitorManagementSystem/db/sqlc"
	"github.com/gin-gonic/gin"
)

// CreateAppointmentLogRequest struct to bind request for creating appointment log
type createAppointmentLogRequest struct {
	AppointmentID int       `json:"appointment_id" binding:"required"`
	CheckInTime   time.Time `json:"check_in_time"`
	CheckOutTime  time.Time `json:"check_out_time"`
}

// Function to create an appointment log
func (server *Server) createAppointmentLog(ctx *gin.Context) {
	var req createAppointmentLogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAppointmentLogParams{
		AppointmentID: int32(req.AppointmentID),
		CheckInTime:   sql.NullTime{Time: req.CheckInTime, Valid: !req.CheckInTime.IsZero()},
		CheckOutTime:  sql.NullTime{Time: req.CheckOutTime, Valid: !req.CheckOutTime.IsZero()},
	}

	log, err := server.store.CreateAppointmentLog(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, log)
}

// GetAppointmentLogRequest struct to bind request for getting appointment log by appointment ID
type getAppointmentLogRequest struct {
	AppointmentID int `uri:"appointment_id" binding:"required"`
}

// Function to get the appointment log for a specific appointment ID
func (server *Server) getAppointmentLogByAppointmentID(ctx *gin.Context) {
	var req getAppointmentLogRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	log, err := server.store.GetAppointmentLogByAppointmentID(ctx, int32(req.AppointmentID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, log)
}

// UpdateCheckInTimeRequest struct to bind request for updating check-in time
type updateCheckInTimeRequest struct {
	AppointmentID int       `json:"appointment_id" binding:"required"`
	CheckInTime   time.Time `json:"check_in_time" binding:"required"`
}

// Function to update the check-in time for a specific appointment log
func (server *Server) updateCheckInTime(ctx *gin.Context) {
	var req updateCheckInTimeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	log, err := server.store.UpdateCheckInTime(ctx, db.UpdateCheckInTimeParams{
		AppointmentID: int32(req.AppointmentID),
		CheckInTime:   sql.NullTime{Time: req.CheckInTime, Valid: !req.CheckInTime.IsZero()},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, log)
}

// UpdateCheckOutTimeRequest struct to bind request for updating check-out time
type updateCheckOutTimeRequest struct {
	AppointmentID int       `json:"appointment_id" binding:"required"`
	CheckOutTime  time.Time `json:"check_out_time" binding:"required"`
}

// Function to update the check-out time for a specific appointment log
func (server *Server) updateCheckOutTime(ctx *gin.Context) {
	var req updateCheckOutTimeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	log, err := server.store.UpdateCheckOutTime(ctx, db.UpdateCheckOutTimeParams{
		AppointmentID: int32(req.AppointmentID),
		CheckOutTime:  sql.NullTime{Time: req.CheckOutTime, Valid: !req.CheckOutTime.IsZero()},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, log)
}

// DeleteAppointmentLogRequest struct to bind request for deleting appointment log
type deleteAppointmentLogRequest struct {
	AppointmentID int `json:"appointment_id" binding:"required"`
}

// Function to delete an appointment log by appointment ID
func (server *Server) deleteAppointmentLog(ctx *gin.Context) {
	var req deleteAppointmentLogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteAppointmentLog(ctx, int32(req.AppointmentID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Appointment log deleted"})
}
