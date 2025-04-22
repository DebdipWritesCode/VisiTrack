package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/DebdipWritesCode/VisitorManagementSystem/db/sqlc"
	"github.com/gin-gonic/gin"
)

// CreateOTPRequest struct to bind request for creating OTP
type createOTPRequest struct {
	PhoneNumber string    `json:"phone_number" binding:"required"`
	OTPCode     string    `json:"otp_code" binding:"required"`
	ExpiresAt   time.Time `json:"expires_at" binding:"required"`
}

// Function to create an OTP for a phone number
func (server *Server) createOTP(ctx *gin.Context) {
	var req createOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateOTPParams{
		PhoneNumber: sql.NullString{String: req.PhoneNumber, Valid: req.PhoneNumber != ""},
		OtpCode:     sql.NullString{String: req.OTPCode, Valid: req.OTPCode != ""},
		ExpiresAt:   sql.NullTime{Time: req.ExpiresAt, Valid: !req.ExpiresAt.IsZero()},
	}

	otp, err := server.store.CreateOTP(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, otp)
}

// GetOTPByPhoneRequest struct to bind request for getting OTP by phone number
type getOTPByPhoneRequest struct {
	PhoneNumber string `uri:"phone_number" binding:"required"`
}

// Function to get the most recent OTP for a phone number
func (server *Server) getOTPByPhone(ctx *gin.Context) {
	var req getOTPByPhoneRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	otp, err := server.store.GetOTPByPhone(ctx, sql.NullString{String: req.PhoneNumber, Valid: req.PhoneNumber != ""})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, otp)
}

// DeleteOTPByPhoneRequest struct to bind request for deleting OTP by phone number
type deleteOTPByPhoneRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

// Function to delete OTP for a specific phone number
func (server *Server) deleteOTPByPhone(ctx *gin.Context) {
	var req deleteOTPByPhoneRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteOTPByPhone(ctx, sql.NullString{String: req.PhoneNumber, Valid: req.PhoneNumber != ""})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "OTP deleted"})
}

// Function to delete expired OTPs
func (server *Server) deleteExpiredOTPs(ctx *gin.Context) {
	err := server.store.DeleteExpiredOTPs(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Expired OTPs deleted"})
}
