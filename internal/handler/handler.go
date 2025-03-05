package handler

import (
	"context"
	"url-server/internal/service"
	pb "url-server/internal/service/proto"

	"google.golang.org/grpc"
)

type Handler struct {
	service *service.MyURLServer
	pb.UnimplementedURLServer
}

func NewHandler(service *service.MyURLServer) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(grpc *grpc.Server) {
	pb.RegisterURLServer(grpc, h)
}

func (h *Handler) GetFullURL(ctx context.Context, r *pb.Request) (*pb.Response, error) {
	return h.service.GetFullURL(ctx, r)
}

func (h *Handler) CreateShortURL(ctx context.Context, r *pb.Request) (*pb.Response, error) {
	return h.service.CreateShortURL(ctx, r)
}
