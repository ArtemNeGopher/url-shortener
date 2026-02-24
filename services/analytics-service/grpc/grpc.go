package grpc

import pb "github.com/ArtemNeGopher/url-shortener/pkg/genproto/analytics"

type service struct {
	pb.UnimplementedAnalyticsServiceServer
}

func NewService() *service {
	return &service{}
}
