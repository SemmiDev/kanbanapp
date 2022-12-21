package gapi

import (
	"context"

	"github.com/SemmiDev/kanbanapp/internal/entity"
	"github.com/SemmiDev/kanbanapp/internal/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	userArg := entity.User{
		Fullname: request.GetFullname(),
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	}

	if userArg.Fullname == "" || userArg.Email == "" || userArg.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "register data is empty")
	}

	userResult, err := s.userService.Register(ctx, &userArg)
	if err != nil {
		if err.Error() == "email already exists" {
			return nil, status.Errorf(codes.Internal, "email already exists")
		}
		return nil, status.Errorf(codes.Internal, "error internal server")
	}

	response := pb.CreateUserResponse{
		UserId:  float64(userResult.ID), // because json encode int to float
		Message: "register success",
	}

	return &response, nil
}
