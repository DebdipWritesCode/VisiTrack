package gapi

import (
	"context"
	"database/sql"

	db "github.com/DebdipWritesCode/VisitorManagementSystem/db/sqlc"
	"github.com/DebdipWritesCode/VisitorManagementSystem/pb"
	"github.com/DebdipWritesCode/VisitorManagementSystem/val"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if len(violations) > 0 {
		return nil, invalidArgumentError(violations)
	}

	arg := db.CreateUserParams{
		PhoneNumber: req.GetPhoneNumber(),
		FirstName:   req.GetFirstName(),
		LastName:    req.GetLastName(),
		Role:        sql.NullString{String: req.Role, Valid: req.Role != ""},
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return nil, status.Errorf(codes.AlreadyExists, "User with phone number %s already exists %s", req.GetPhoneNumber(), err.Error())
		}
		return nil, status.Errorf(codes.Internal, "Failed to create user: %s", err.Error())
	}

	rsp := &pb.CreateUserResponse{
		User: convertUserToProto(user),
	}

	return rsp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidatePhoneNumber(req.GetPhoneNumber()); err != nil {
		violations = append(violations, fieldViolation("phone_number", err))
	}

	if err := val.ValidateFirstName(req.GetFirstName()); err != nil {
		violations = append(violations, fieldViolation("first_name", err))
	}

	if err := val.ValidateLastName(req.GetLastName()); err != nil {
		violations = append(violations, fieldViolation("last_name", err))
	}

	if err := val.ValidateRole(req.GetRole()); err != nil {
		violations = append(violations, fieldViolation("role", err))
	}

	return violations
}
