package server

import (
	"context"
	"github.com/VadimGossip/grpsProductsServer/gen/products"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductService interface {
	List(ctx context.Context, req *products.ListRequest) (*products.ListResponse, error)
	Fetch(ctx context.Context, req *products.FetchRequest) (*emptypb.Empty, error)
}

type ProductServer struct {
	service ProductService
}

func NewProductServer(service ProductService) *ProductServer {
	return &ProductServer{
		service: service,
	}
}

func (s *ProductServer) List(ctx context.Context, req *products.ListRequest) (*products.ListResponse, error) {
	return s.service.List(ctx, req)
}

func (s *ProductServer) Fetch(ctx context.Context, req *products.FetchRequest) (*emptypb.Empty, error) {
	return s.service.Fetch(ctx, req)
}
