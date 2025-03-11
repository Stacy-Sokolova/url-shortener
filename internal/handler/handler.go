package handler

import (
	"context"
	"url-server/internal/service"
	pb "url-server/internal/service/proto"

	"github.com/sirupsen/logrus"
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
	url := r.GetUrl()
	got, err := h.services.URLService.GetFullURL(ctx, url)
	if err != nil {
		logrus.Debugf("error while getting original url: %s", err.Error())
		return nil, err
	}

	logrus.Info("OK GetFullURL")
	return &pb.Response{Url: got}, nil
}

func (h *Handler) CreateShortURL(ctx context.Context, r *pb.Request) (*pb.Response, error) {
	url := r.GetUrl()
	got, err := h.services.URLService.CreateShortURL(ctx, url)
	if err != nil {
		logrus.Debugf("error while generating short URL: %s", err.Error())
		return nil, err
	}

	logrus.Info("OK CreateShortURL")
	return &pb.Response{Url: got}, nil
}
