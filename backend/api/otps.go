package api

import (
	"log"
	"net/http"

	"github.com/DebdipWritesCode/VisitorManagementSystem/util"
	"github.com/gin-gonic/gin"
)

type sendOTPRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type verifyOTPRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	OTPCode     string `json:"otp_code" binding:"required"`
}

func (server *Server) sendOTP(ctx *gin.Context) {
	var req sendOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := util.StartPhoneVerification(req.PhoneNumber)
	if err != nil {
		// Print actual error on server log
		log.Printf("❌ Failed to send OTP to %s: %v\n", req.PhoneNumber, err)

		// Send generic error response to frontend
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

func (server *Server) verifyOTP(ctx *gin.Context) {
	var req verifyOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ok, err := util.CheckPhoneVerification(req.PhoneNumber, req.OTPCode)
	if err != nil || !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Phone number verified successfully"})
}
