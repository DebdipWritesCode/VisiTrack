package api

import (
	"fmt"

	db "github.com/DebdipWritesCode/VisitorManagementSystem/db/sqlc"
	"github.com/DebdipWritesCode/VisitorManagementSystem/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config util.Config
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	fmt.Println("Initializing Twilio client...")
	util.InitTwilio()

	server := &Server{
		config: config,
		store:  store,
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) GetRouter() *gin.Engine {
	return server.router
}

// setupRouter initializes the Gin router with all routes.
func (server *Server) setupRouter() {
	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	// Appointment routes
	router.POST("/appointments", server.createAppointment)
	router.GET("/appointments/:id", server.getAppointmentByID)
	router.GET("/appointments/visitor/:id", server.listAppointmentsByVisitor)
	router.GET("/appointments/host/:id", server.listAppointmentsByHost)
	router.GET("/appointments/date", server.listAppointmentsByDate)
	router.GET("/appointments/qr/:qr_code", server.getAppointmentByQRCode)
	router.PUT("/appointments/status", server.updateAppointmentStatus)
	router.GET("/users/:id/stats", server.getUserAppointmentStats)
	router.DELETE("/appointments/:id", server.deleteAppointment)
	router.POST("/appointments/:id/cancel", server.cancelAppointment)

	// User routes
	router.POST("/users", server.createUser)
	router.POST("/auth/signup", server.signupUser)
	router.POST("/auth/login", server.loginUser)
	router.GET("/users/:id", server.getUserByID)
	router.GET("/users/phone/:phone_number", server.getUserByPhone)
	router.GET("/users", server.listUsers)
	router.PUT("/users/name", server.updateUserName)
	router.PUT("/users/role", server.updateUserRole)
	router.DELETE("/users/:id", server.deleteUser)
	router.GET("/users/search", server.getUsersByName)

	// User appointment stats
	router.GET("/users/:id/appointments/hosted", server.getTotalAppointmentsHosted)
	router.GET("/users/:id/appointments/visited", server.getTotalAppointmentsVisited)
	router.GET("/users/popular", server.getTopPopularUsers)

	// Appointment Stats routes
	router.POST("/appointment_stats", server.createAppointmentStats)
	router.GET("/appointment_stats/:user_id", server.getAppointmentStatsByUserID)
	router.PUT("/appointment_stats/increment", server.incrementAppointmentCount)
	router.PUT("/appointment_stats/decrement", server.decrementAppointmentCount)
	router.PUT("/appointment_stats/reset", server.resetAppointmentCount)
	router.GET("/appointment_stats/popular", server.getTopPopularUsers)
	router.DELETE("/appointment_stats/:user_id", server.deleteAppointmentStats)

	// Availability routes
	router.POST("/availability", server.createAvailabilitySlot)
	router.GET("/availability/:user_id", server.getAvailabilityByUser)
	router.PUT("/availability/status", server.updateAvailabilityStatus)
	router.DELETE("/availability", server.deleteAvailabilitySlot)
	router.DELETE("/availability/:user_id", server.deleteAvailabilityByUser)

	// OTP routes using Twilio
	router.POST("/otp/send", server.sendOTP)
	router.POST("/otp/verify", server.verifyOTP)

	// Appointment Log routes
	router.POST("/appointment_logs", server.createAppointmentLog)
	router.GET("/appointment_logs/:appointment_id", server.getAppointmentLogByAppointmentID)
	router.PUT("/appointment_logs/check_in", server.updateCheckInTime)
	router.PUT("/appointment_logs/check_out", server.updateCheckOutTime)
	router.DELETE("/appointment_logs", server.deleteAppointmentLog)

	server.router = router
}

// Start runs the HTTP server.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse standardizes error responses.
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
