package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/DebdipWritesCode/VisitorManagementSystem/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createAppointmentRequest struct {
	VisitorID       int64     `json:"visitor_id" binding:"required"`
	HostID          int64     `json:"host_id" binding:"required"`
	AppointmentDate time.Time `json:"appointment_date" binding:"required"`
	StartTime       time.Time `json:"start_time" binding:"required"`
	EndTime         time.Time `json:"end_time" binding:"required"`
	Status          *string   `json:"status"` // Optional — will default to 'pending'
	QRCode          string    `json:"qr_code" binding:"required"`
}

func (server *Server) createAppointment(ctx *gin.Context) {
	var req createAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAppointmentParams{
		VisitorID:       int32(req.VisitorID),
		HostID:          int32(req.HostID),
		AppointmentDate: req.AppointmentDate,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		Status:          sql.NullString{String: "pending", Valid: true},
		QrCode:          sql.NullString{String: req.QRCode, Valid: true},
	}

	if req.Status != nil {
		arg.Status = sql.NullString{String: *req.Status, Valid: true}
	}

	appointment, err := server.store.CreateAppointment(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, appointment)
}

type getAppointmentUriRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAppointmentByID(ctx *gin.Context) {
	var req getAppointmentUriRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	appointment, err := server.store.GetAppointmentByID(ctx, int32(req.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, appointment)
}

type listByIDUri struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) listAppointmentsByVisitor(ctx *gin.Context) {
	var req listByIDUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	appointments, err := server.store.ListAppointmentsByVisitor(ctx, int32(req.ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, appointments)
}

func (server *Server) listAppointmentsByHost(ctx *gin.Context) {
	var req listByIDUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	appointments, err := server.store.ListAppointmentsByHost(ctx, int32(req.ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, appointments)
}

type listByDateRequest struct {
	Date string `form:"date" binding:"required"`
}

// List appointments by date
func (server *Server) listAppointmentsByDate(ctx *gin.Context) {
	var req listByDateRequest
	// Bind the query parameter as date
	if err := ctx.ShouldBindQuery(&req); err != nil {
		// If there's an error, respond with BadRequest
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Parse the date string to time.Time
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid date format, use YYYY-MM-DD")))
		return
	}

	// Use the parsed time.Time object for the database query
	appointments, err := server.store.ListAppointmentsByDate(ctx, parsedDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, appointments)
}

type getAppointmentByQRCodeRequest struct {
	QRCode string `uri:"qr_code" binding:"required"`
}

func (server *Server) getAppointmentByQRCode(ctx *gin.Context) {
	var req getAppointmentByQRCodeRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	appointment, err := server.store.GetAppointmentByQRCode(ctx, sql.NullString{String: req.QRCode, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("no appointment found for this QR code")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, appointment)
}

type updateAppointmentStatusRequest struct {
	ID     int64  `json:"id" binding:"required,min=1"`
	Status string `json:"status" binding:"required,oneof=pending ongoing cancelled completed"`
}

func (server *Server) updateAppointmentStatus(ctx *gin.Context) {
	var req updateAppointmentStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateAppointmentStatusParams{
		ID:     int32(req.ID),
		Status: sql.NullString{String: req.Status, Valid: true},
	}

	appointment, err := server.store.UpdateAppointmentStatus(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("no appointment found with this ID")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, appointment)
}

type getUserAppointmentStatsRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUserAppointmentStats(ctx *gin.Context) {
	var req getUserAppointmentStatsRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	stats, err := server.store.GetUserAppointmentStats(ctx, int32(req.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("no user found with this ID")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

func (server *Server) cancelAppointment(ctx *gin.Context) {
	var req getAppointmentUriRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	appointment, err := server.store.CancelAppointment(ctx, int32(req.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, appointment)
}

func (server *Server) deleteAppointment(ctx *gin.Context) {
	var req getAppointmentUriRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteAppointment(ctx, int32(req.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "appointment deleted"})
}
