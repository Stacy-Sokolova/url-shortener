package service

import (
	context "context"
	pb "url-server/internal/service/proto"
	"url-server/internal/storage"
)

type MyURLServer struct {
	storage storage.Storage
	pb.UnimplementedURLServer
}

func NewURLServer(strg storage.Storage) *MyURLServer {
	return &MyURLServer{storage: strg}
}

func (s *MyURLServer) GetFullURL(ctx context.Context, r *pb.Request) (*pb.Response, error) {
	return s.storage.GetFullURL(ctx, r)
}

func (s *MyURLServer) CreateShortURL(ctx context.Context, r *pb.Request) (*pb.Response, error) {
	return s.storage.CreateShortURL(ctx, r)
}
