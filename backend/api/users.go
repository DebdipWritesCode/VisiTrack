package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/DebdipWritesCode/VisitorManagementSystem/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required,e164"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Role        string `json:"role" binding:"omitempty,oneof=admin user"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		PhoneNumber: req.PhoneNumber,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Role:        sql.NullString{String: req.Role, Valid: req.Role != ""},
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	const timeLayout = "15:04:05"

	for day := 1; day <= 5; day++ {
		for hour := 9; hour < 18; hour++ {
			startStr := fmt.Sprintf("%02d:00:00", hour)
			endStr := fmt.Sprintf("%02d:00:00", hour+1)

			startTime, err := time.Parse(timeLayout, startStr)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to parse start time: %w", err)))
				return
			}

			endTime, err := time.Parse(timeLayout, endStr)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to parse end time: %w", err)))
				return
			}

			slot := db.CreateAvailabilitySlotParams{
				UserID:    user.ID,
				DayOfWeek: int32(day),
				StartTime: startTime,
				EndTime:   endTime,
			}

			_, err = server.store.CreateAvailabilitySlot(ctx, slot)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create availability: %w", err)))
				return
			}
		}
	}

	ctx.JSON(http.StatusOK, user)
}

type signupUserRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required,e164"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Role        string `json:"role" binding:"omitempty,oneof=admin user"`
}

func (server *Server) signupUser(ctx *gin.Context) {
	var req signupUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		PhoneNumber: req.PhoneNumber,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Role:        sql.NullString{String: req.Role, Valid: req.Role != ""},
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

type loginUserRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required,e164"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByPhone(ctx, req.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type getUserUriRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUserByID(ctx *gin.Context) {
	var req getUserUriRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByID(ctx, int32(req.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type getUserByPhoneRequest struct {
	PhoneNumber string `uri:"phone_number" binding:"required,e164"`
}

func (server *Server) getUserByPhone(ctx *gin.Context) {
	var req getUserByPhoneRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByPhone(ctx, req.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type listUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=100"`
}

func (server *Server) listUsers(ctx *gin.Context) {
	var req listUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := server.store.ListUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

type updateUserNameRequest struct {
	ID        int64  `json:"id" binding:"required,min=1"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

func (server *Server) updateUserName(ctx *gin.Context) {
	var req updateUserNameRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateUserNameParams{
		ID:        int32(req.ID),
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	user, err := server.store.UpdateUserName(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type updateUserRoleRequest struct {
	ID   int64  `json:"id" binding:"required,min=1"`
	Role string `json:"role" binding:"required,oneof=admin user"`
}

func (server *Server) updateUserRole(ctx *gin.Context) {
	var req updateUserRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateUserRoleParams{
		ID:   int32(req.ID),
		Role: sql.NullString{String: req.Role, Valid: req.Role != ""},
	}

	user, err := server.store.UpdateUserRole(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type getUsersByNameRequest struct {
	Query string `form:"query" binding:"required,min=1"`
}

func (server *Server) getUsersByName(ctx *gin.Context) {
	var req getUsersByNameRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	users, err := server.store.GetUsersByName(ctx, sql.NullString{String: req.Query, Valid: req.Query != ""})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

type deleteUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteUser(ctx *gin.Context) {
	var req deleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteUser(ctx, int32(req.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func (server *Server) getTopPopularUsers(ctx *gin.Context) {
	users, err := server.store.GetTopPopularUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (server *Server) getTotalAppointmentsVisited(ctx *gin.Context) {
	var req getUserIDUriRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	count, err := server.store.GetTotalAppointmentsVisited(ctx, int32(req.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"appointments_visited": count})
}

type getUserIDUriRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getTotalAppointmentsHosted(ctx *gin.Context) {
	var req getUserIDUriRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	count, err := server.store.GetTotalAppointmentsHosted(ctx, int32(req.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"appointments_hosted": count})
}
