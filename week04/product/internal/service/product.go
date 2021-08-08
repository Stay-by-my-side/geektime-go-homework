package service

import (
	"context"
	"product/internal/biz"

	pb "product/api/product/v1"

	"github.com/go-kratos/kratos/v2/log"
)

func NewProductService(product *biz.ProductUsecase, logger log.Logger) *ProductService {
	return &ProductService{
		product: product,
		log:     log.NewHelper(logger),
	}
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
	s.log.WithContext(ctx).Infof("GetProduct Received: %v", req.GetId())
	return &pb.GetProductReply{}, nil
}
func (s *ProductService) ListProduct(ctx context.Context, req *pb.ListProductRequest) (*pb.ListProductReply, error) {
	return &pb.ListProductReply{}, nil
}
