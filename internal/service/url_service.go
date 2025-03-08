package service

import (
	context "context"
	pb "url-server/internal/service/proto"
	"url-server/internal/storage"
)

type MyURLService struct {
	storage storage.Storage
}

func NewURLService(strg storage.Storage) *MyURLService {
	return &MyURLService{storage: strg}
}

func (s *MyURLService) GetFullURL(ctx context.Context, r *pb.Request) (*pb.Response, error) {
	return s.storage.GetFullURL(ctx, r)
}

func (s *MyURLService) CreateShortURL(ctx context.Context, r *pb.Request) (*pb.Response, error) {
	return s.storage.CreateShortURL(ctx, r)
}
