package data

import (
	"context"
	"helloworld/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type productRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewProductRepo(data *Data, logger log.Logger) biz.ProductRepo {
	return &productRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *productRepo) ListProduct(ctx context.Context) ([]*biz.Product, error) {
	return nil, nil
}

func (r *productRepo) GetProduct(ctx context.Context, id int64) (*biz.Product, error) {
	return nil, nil
}

func (r *productRepo) CreateProduct(ctx context.Context, p *biz.Product) error {
	return nil
}

func (r *productRepo) UpdateProduct(ctx context.Context, id int64, p *biz.Product) error {
	return nil
}

func (r *productRepo) DeleteProduct(ctx context.Context, id int64) error {
	return nil
}
