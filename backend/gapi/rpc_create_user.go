package gapi

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	db "github.com/DebdipWritesCode/VisitorManagementSystem/db/sqlc"
	"github.com/DebdipWritesCode/VisitorManagementSystem/pb"
	"github.com/DebdipWritesCode/VisitorManagementSystem/val"
	"github.com/DebdipWritesCode/VisitorManagementSystem/worker"
	"github.com/hibiken/asynq"
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

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			PhoneNumber: req.GetPhoneNumber(),
			FirstName:   req.GetFirstName(),
			LastName:    req.GetLastName(),
			Role:        sql.NullString{String: req.Role, Valid: req.Role != ""},
		},
		AfterCreate: func(user db.User) error {
			// TODO: Use db transaction to ensure that we can rollback if the task distribution fails.
			taskPayload := &worker.PayloadSendVerifyEmail{
				UserID: strconv.Itoa(int(user.ID)),
			}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical), // Use a critical queue for important tasks
			}

			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
		},
	}

	txResult, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return nil, status.Errorf(codes.AlreadyExists, "User with phone number %s already exists %s", req.GetPhoneNumber(), err.Error())
		}
		return nil, status.Errorf(codes.Internal, "Failed to create user: %s", err.Error())
	}

	rsp := &pb.CreateUserResponse{
		User: convertUserToProto(txResult.User),
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
