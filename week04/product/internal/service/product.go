package service

import (
	"context"
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
	// TODO 修改验证与错误处理
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
		s.log.WithContext(ctx).Errorf("product.internal.service.product.CreateProduct error: \n%+v\n", err)
		return nil, errors.New("保存商品失败")
	}
	p, err := s.product.Get(ctx, id)
	if err != nil {
		s.log.WithContext(ctx).Errorf("product.internal.service.product.CreateProduct error: \n%+v\n", err)
		return nil, errors.New("查询商品失败")
	}
	return &pb.CreateProductReply{
		Product: s.DTO(p),
	}, nil
}
func (s *ProductService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductReply, error) {
	s.log.WithContext(ctx).Infof("UpdateProduct Received: %v", req)
	// TODO 修改验证与错误处理
	if req.Id <= 0 {
		return nil, errors.New("商品Id不能为空")
	}
	if req.Description == "" {
		return nil, errors.New("未设置描述")
	}
	if req.Price == 0.0 {
		return nil, errors.New("价格不能为0")
	}
	err := s.product.Update(ctx, req.Id, &biz.Product{
		Description: req.Description,
		Price:       req.Price,
	})
	if err != nil {
		s.log.WithContext(ctx).Errorf("product.internal.service.product.UpdateProduct error: \n%+v\n", err)
		return nil, errors.New("更新商品失败")
	}
	return &pb.UpdateProductReply{}, nil
}
func (s *ProductService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductReply, error) {
	s.log.WithContext(ctx).Infof("DeleteProduct Received: %v", req.GetId())
	// TODO 修改验证与错误处理
	if req.Id <= 0 {
		return nil, errors.New("商品Id不能为空")
	}
	err := s.product.Delete(ctx, req.Id)
	if err != nil {
		s.log.WithContext(ctx).Errorf("product.internal.service.product.DeleteProduct error: \n%+v\n", err)
		return nil, errors.New("删除商品失败")
	}
	return &pb.DeleteProductReply{}, nil
}
func (s *ProductService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductReply, error) {
	s.log.WithContext(ctx).Infof("GetProduct Received: %v", req.GetId())
	// TODO 修改验证与错误处理
	if req.Id <= 0 {
		return nil, errors.New("商品Id不能为空")
	}
	p, err := s.product.Get(ctx, req.Id)
	if err != nil {
		s.log.WithContext(ctx).Errorf("product.internal.service.product.GetProduct error: \n%+v\n", err)
		return nil, errors.New("查询商品失败")
	}
	return &pb.GetProductReply{
		Product: s.DTO(p),
	}, nil
}
func (s *ProductService) ListProduct(ctx context.Context, req *pb.ListProductRequest) (*pb.ListProductReply, error) {
	s.log.WithContext(ctx).Infof("ListProduct Received: %v", req)
	ps, err := s.product.List(ctx)
	if err != nil {
		s.log.WithContext(ctx).Errorf("product.internal.service.product.ListProduct error: \n%+v\n", err)
		return nil, errors.New("查询列表出错")
	}
	return &pb.ListProductReply{
		Results: s.DTOs(ps),
	}, nil
}

func (s ProductService) DTO(p *biz.Product) *pb.ProductInfo {
	if p == nil {
		return nil
	}
	return &pb.ProductInfo{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}
}

func (s ProductService) DTOs(ps []*biz.Product) []*pb.ProductInfo {
	if ps == nil {
		return nil
	}
	pbs := make([]*pb.ProductInfo, len(ps), len(ps))
	for i := 0; i < len(ps); i++ {
		pb := s.DTO(ps[i])
		if pb == nil {
			continue
		}
		pbs[i] = pb
	}
	return pbs
}
