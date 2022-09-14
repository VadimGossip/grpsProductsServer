package server

import (
	"fmt"
	"net"

	"github.com/VadimGossip/grpsProductsServer/gen/products"
	"google.golang.org/grpc"
)

type Server struct {
	grpcSrv        *grpc.Server
	productsServer products.ProductsServiceServer
}

func NewServer(productsServer products.ProductsServiceServer) *Server {
	return &Server{
		grpcSrv:        grpc.NewServer(),
		productsServer: productsServer,
	}
}

func (s *Server) ListenAndServe(port int) error {
	addr := fmt.Sprintf(":%d", port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	products.RegisterProductsServiceServer(s.grpcSrv, s.productsServer)

	if err := s.grpcSrv.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop() func() {
	return s.grpcSrv.GracefulStop
}
