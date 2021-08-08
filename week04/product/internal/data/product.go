package data

import (
	"context"
	"product/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
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

type product struct {
	id          int64   `db:"id"`
	name        string  `db:"name"`
	description string  `db:"description"`
	price       float64 `db:"price"`
}

func (r *productRepo) ListProduct(ctx context.Context) ([]*biz.Product, error) {
	return nil, nil
}

func (r *productRepo) GetProduct(ctx context.Context, id int64) (*biz.Product, error) {
	sqlStr := "SELECT id, name, description, price FROM product WHERE id = ?"
	var p product
	if err := r.data.db.Get(&p, sqlStr, id); err != nil {
		return nil, errors.Wrapf(err, "sql: %s", sqlStr)
	}
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
