package gapi

import (
	"context"
	"database/sql"

	"github.com/DebdipWritesCode/VisitorManagementSystem/pb"
	"github.com/DebdipWritesCode/VisitorManagementSystem/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := validateLoginUserRequest(req)
	if len(violations) > 0 {
		return nil, invalidArgumentError(violations)
	}

	user, err := server.store.GetUserByPhone(ctx, req.GetPhoneNumber())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "User with phone number %s not found %s", req.GetPhoneNumber(), err.Error())
		}
		return nil, status.Errorf(codes.Internal, "Failed to retrieve user: %s", err.Error())
	}

	mtdt := server.extractMetadata(ctx)
	rsq := &pb.LoginUserResponse{
		User:     convertUserToProto(user),
		Metadata: convertMetadataToProto(mtdt),
	}

	return rsq, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidatePhoneNumber(req.GetPhoneNumber()); err != nil {
		violations = append(violations, fieldViolation("phone_number", err))
	}

	return violations
}
