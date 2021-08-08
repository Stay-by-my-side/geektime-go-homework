package service

import (
	"context"
	"fmt"
	"product/internal/biz"

	"github.com/pkg/errors"

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
	s.log.WithContext(ctx).Infof("CreateProduct Received: %v", req)
	if req.Name == "" {
		return nil, errors.New("未设置名称")
	}
	if req.Description == "" {
		return nil, errors.New("未设置描述")
	}
	if req.Price == 0.0 {
		return nil, errors.New("价格不能为0")
	}
	id, err := s.product.Create(ctx, &biz.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})
	if err != nil {
		return nil, err
	}
	fmt.Printf("insert id: %d\n", id)
	p, err := s.product.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &pb.CreateProductReply{
		Product: &pb.ProductInfo{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		},
	}, nil
}
func (s *ProductService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductReply, error) {
	s.log.WithContext(ctx).Infof("UpdateProduct Received: %v", req)
	return &pb.UpdateProductReply{}, nil
}
func (s *ProductService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductReply, error) {
	s.log.WithContext(ctx).Infof("DeleteProduct Received: %v", req.GetId())
	return &pb.DeleteProductReply{}, nil
}
func (s *ProductService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductReply, error) {
	s.log.WithContext(ctx).Infof("GetProduct Received: %v", req.GetId())
	return &pb.GetProductReply{}, nil
}
func (s *ProductService) ListProduct(ctx context.Context, req *pb.ListProductRequest) (*pb.ListProductReply, error) {
	s.log.WithContext(ctx).Infof("ListProduct Received: %v", req)
	return &pb.ListProductReply{}, nil
}
