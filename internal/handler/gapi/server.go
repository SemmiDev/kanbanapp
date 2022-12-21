package gapi

import (
	"github.com/SemmiDev/kanbanapp/internal/pb"
	"github.com/SemmiDev/kanbanapp/internal/service"
)

type Server struct {
	pb.UnimplementedKanbanServer
	userService service.UserService
}

func NewServer(userService service.UserService) (*Server, error) {
	return &Server{
		userService: userService,
	}, nil
}
