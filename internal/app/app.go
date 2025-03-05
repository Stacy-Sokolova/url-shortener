package app

import (
	"fmt"
	"net"

	"url-server/internal/handler"

	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
}

func (s *Server) Run(port string, handler *handler.Handler) error {
	lis, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		return err
	}

	s.grpcServer = grpc.NewServer()
	handler.Register(s.grpcServer)
	fmt.Printf("Server listening at %v", lis.Addr())

	if err := s.grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (s *Server) Shutdown() {
	s.grpcServer.GracefulStop()
}
