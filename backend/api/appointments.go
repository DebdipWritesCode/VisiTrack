package api

import (
	"database/sql"
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
	Date time.Time `form:"date" binding:"required"`
}

func (server *Server) listAppointmentsByDate(ctx *gin.Context) {
	var req listByDateRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	appointments, err := server.store.ListAppointmentsByDate(ctx, req.Date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, appointments)
}

type updateAppointmentStatusRequest struct {
	ID     int64  `json:"id" binding:"required,min=1"`
	Status string `json:"status" binding:"required,oneof=pending confirmed cancelled completed"`
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
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, appointment)
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
