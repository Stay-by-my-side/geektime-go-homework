package service

import (
    "context"

    pb "helloworld/api/product/v1"
)

type ProductService struct {
	pb.UnimplementedProductServer
}

func NewProductService() *ProductService {
	return &ProductService{}
}

func (s *ProductService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductReply, error) {
	return &pb.CreateProductReply{}, nil
}
func (s *ProductService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductReply, error) {
	return &pb.UpdateProductReply{}, nil
}
func (s *ProductService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductReply, error) {
	return &pb.DeleteProductReply{}, nil
}
func (s *ProductService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductReply, error) {
	return &pb.GetProductReply{}, nil
}
func (s *ProductService) ListProduct(ctx context.Context, req *pb.ListProductRequest) (*pb.ListProductReply, error) {
	return &pb.ListProductReply{}, nil
}
