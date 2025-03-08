package service

import (
	context "context"
	pb "url-server/internal/service/proto"
	"url-server/internal/storage"
)

type URLService interface {
	GetFullURL(ctx context.Context, r *pb.Request) (*pb.Response, error)
	CreateShortURL(ctx context.Context, r *pb.Request) (*pb.Response, error)
}

type Service struct {
	URLService URLService
}

func NewService(storage storage.Storage) *Service {
	return &Service{
		URLService: NewURLService(storage),
	}
}
