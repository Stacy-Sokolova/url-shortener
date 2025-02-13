package storage

import (
	"context"
	pb "url-server/internal/service/proto"
)

type Storage interface {
	GetFullURL(ctx context.Context, r *pb.Request) (*pb.Response, error)
	CreateShortURL(ctx context.Context, r *pb.Request) (*pb.Response, error)
}
