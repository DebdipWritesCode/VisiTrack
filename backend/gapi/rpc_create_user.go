package gapi

import (
	"context"
	"database/sql"

	db "github.com/DebdipWritesCode/VisitorManagementSystem/db/sqlc"
	"github.com/DebdipWritesCode/VisitorManagementSystem/pb"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
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
