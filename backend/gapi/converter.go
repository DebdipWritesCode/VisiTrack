package gapi

import (
	db "github.com/DebdipWritesCode/VisitorManagementSystem/db/sqlc"
	"github.com/DebdipWritesCode/VisitorManagementSystem/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUserToProto(user db.User) *pb.User {
	return &pb.User{
		Id:          int64(user.ID),
		PhoneNumber: user.PhoneNumber,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Role:        user.Role.String,
		CreatedAt: func() *timestamppb.Timestamp {
			if user.CreatedAt.Valid {
				return timestamppb.New(user.CreatedAt.Time)
			}
			return nil
		}(),
		AppointmentsBooked: int32(user.AppointmentsVisited.Int32),
		AppointmentsHosted: int32(user.AppointmentsHosted.Int32),
	}
}
