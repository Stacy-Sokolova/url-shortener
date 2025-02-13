package app

import (
	"fmt"
	"net"

	s "url-server/internal/service"
	pb "url-server/internal/service/proto"

	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
}

func (s *Server) Run(port string, service *s.MyURLServer) error {
	lis, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		return err
		//log.Fatalln("Cant listen port", err)
	}

	s.grpcServer = grpc.NewServer()
	pb.RegisterURLServer(s.grpcServer, service)
	fmt.Printf("Server listening at %v", lis.Addr())

	if err := s.grpcServer.Serve(lis); err != nil {
		return err
		//log.Fatalf("failed to serve: %v", err)
	}

	return nil
}

func (s *Server) Shutdown() {
	s.grpcServer.GracefulStop()
}
