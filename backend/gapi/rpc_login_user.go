package gapi

import (
	"context"
	"database/sql"

	"github.com/DebdipWritesCode/VisitorManagementSystem/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUserByPhone(ctx, req.GetPhoneNumber())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "User with phone number %s not found %s", req.GetPhoneNumber(), err.Error())
		}
		return nil, status.Errorf(codes.Internal, "Failed to retrieve user: %s", err.Error())
	}

	rsq := &pb.LoginUserResponse{
		User: convertUserToProto(user),
	}

	return rsq, nil
}
