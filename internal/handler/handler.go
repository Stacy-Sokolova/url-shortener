package handler

import (
	"context"
	"url-server/internal/service"
	pb "url-server/internal/service/proto"

	"google.golang.org/grpc"
)

type Handler struct {
	services *service.Service
	pb.UnimplementedURLServer
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{services: service}
}

func (h *Handler) Register(grpc *grpc.Server) {
	pb.RegisterURLServer(grpc, h)
}

func (h *Handler) GetFullURL(ctx context.Context, r *pb.Request) (*pb.Response, error) {
	return h.services.URLService.GetFullURL(ctx, r)
}

func (h *Handler) CreateShortURL(ctx context.Context, r *pb.Request) (*pb.Response, error) {
	return h.services.URLService.CreateShortURL(ctx, r)
}
