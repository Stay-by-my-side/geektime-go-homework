package service

import (
	pb "product/api/product/v1"
	"product/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewProductService)

type ProductService struct {
	pb.UnimplementedProductServer

	log *log.Helper

	product *biz.ProductUsecase
}
