package api

import (
	"database/sql"
	"net/http"

	db "github.com/DebdipWritesCode/VisitorManagementSystem/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createAppointmentStatsRequest struct {
	UserID            int64 `json:"user_id" binding:"required,min=1"`
	TotalAppointments int32 `json:"total_appointments" binding:"required,min=0"`
}

func (server *Server) createAppointmentStats(ctx *gin.Context) {
	var req createAppointmentStatsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAppointmentStatsParams{
		UserID:            int32(req.UserID),
		TotalAppointments: sql.NullInt32{Int32: req.TotalAppointments, Valid: true},
	}

	appointmentStats, err := server.store.CreateAppointmentStats(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, appointmentStats)
}

type getAppointmentStatsByUserIDRequest struct {
	UserID int64 `uri:"user_id" binding:"required,min=1"`
}

func (server *Server) getAppointmentStatsByUserID(ctx *gin.Context) {
	var req getAppointmentStatsByUserIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	appointmentStats, err := server.store.GetAppointmentStatsByUserID(ctx, int32(req.UserID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, appointmentStats)
}

type incrementAppointmentCountRequest struct {
	UserID int64 `json:"user_id" binding:"required,min=1"`
}

func (server *Server) incrementAppointmentCount(ctx *gin.Context) {
	var req incrementAppointmentCountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	appointmentStats, err := server.store.IncrementAppointmentCount(ctx, int32(req.UserID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, appointmentStats)
}

type decrementAppointmentCountRequest struct {
	UserID int64 `json:"user_id" binding:"required,min=1"`
}

func (server *Server) decrementAppointmentCount(ctx *gin.Context) {
	var req decrementAppointmentCountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	appointmentStats, err := server.store.DecrementAppointmentCount(ctx, int32(req.UserID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, appointmentStats)
}

type resetAppointmentCountRequest struct {
	UserID int64 `json:"user_id" binding:"required,min=1"`
}

func (server *Server) resetAppointmentCount(ctx *gin.Context) {
	var req resetAppointmentCountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	appointmentStats, err := server.store.ResetAppointmentCount(ctx, int32(req.UserID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, appointmentStats)
}

type deleteAppointmentStatsRequest struct {
	UserID int64 `uri:"user_id" binding:"required,min=1"`
}

func (server *Server) deleteAppointmentStats(ctx *gin.Context) {
	var req deleteAppointmentStatsRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteAppointmentStats(ctx, int32(req.UserID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "appointment stats deleted"})
}
