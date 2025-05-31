package gapi

import (
	"context"
	"database/sql"

	db "github.com/DebdipWritesCode/VisitorManagementSystem/db/sqlc"
	"github.com/DebdipWritesCode/VisitorManagementSystem/pb"
	"github.com/DebdipWritesCode/VisitorManagementSystem/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUserName(ctx context.Context, req *pb.UpdateUserNameRequest) (*pb.UpdateUserNameResponse, error) {
	violations := validateUpdateUserNameRequest(req)
	if len(violations) > 0 {
		return nil, invalidArgumentError(violations)
	}

	arg := db.UpdateUserNameParams{
		ID: req.GetUserId(),
		FirstName: sql.NullString{
			String: req.GetFirstName(),
			Valid:  req.GetFirstName() != "",
		},
		LastName: sql.NullString{
			String: req.GetLastName(),
			Valid:  req.GetLastName() != "",
		},
	}

	user, err := server.store.UpdateUserName(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "User with ID %d not found", req.GetUserId())
		}
		return nil, status.Errorf(codes.Internal, "Failed to create user: %s", err.Error())
	}

	rsp := &pb.UpdateUserNameResponse{
		User: convertUserToProto(user),
	}

	return rsp, nil
}

func validateUpdateUserNameRequest(req *pb.UpdateUserNameRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateID(int64(req.GetUserId())); err != nil {
		violations = append(violations, fieldViolation("user_id", err))
	}

	if req.GetFirstName() != "" {
		if err := val.ValidateFirstName(req.GetFirstName()); err != nil {
			violations = append(violations, fieldViolation("first_name", err))
		}
	}

	if req.GetLastName() != "" {
		if err := val.ValidateLastName(req.GetLastName()); err != nil {
			violations = append(violations, fieldViolation("last_name", err))
		}
	}

	return violations
}
