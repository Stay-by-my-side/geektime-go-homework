package data

import (
	"context"
	"encoding/json"
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
	Id          int64   `db:"id"`
	Name        string  `db:"name"`
	Description string  `db:"description"`
	Price       float64 `db:"price"`
}

func pToPB(product product) (*biz.Product, error) {
	pj, err := json.Marshal(product)
	if err != nil {
		return nil, errors.Wrapf(err, "product: %v\n", product)
	}
	pb := new(biz.Product)
	err = json.Unmarshal(pj, pb)
	if err != nil {
		return nil, errors.Wrapf(err, "pj: %v, pb: %v\n", pj, pb)
	}
	return pb, nil
}

func (r *productRepo) ListProduct(ctx context.Context) ([]*biz.Product, error) {
	sqlStr := "SELECT id, name, description, price FROM product"
	var products []product
	if err := r.data.db.Select(&products, sqlStr); err != nil {
		return nil, errors.Wrapf(err, "sql: %s", sqlStr)
	}

	pbs := make([]*biz.Product, len(products), len(products))
	for i := 0; i < len(products); i++ {
		pb, err := pToPB(products[i])
		if err != nil {
			return nil, err
		}
		pbs = append(pbs, pb)
	}
	return pbs, nil
}

func (r *productRepo) GetProduct(ctx context.Context, id int64) (*biz.Product, error) {
	sqlStr := "SELECT id, name, description, price FROM product WHERE id = ?"
	var p product
	if err := r.data.db.Get(&p, sqlStr, id); err != nil {
		return nil, errors.Wrapf(err, "sql: %s", sqlStr)
	}

	pb, err := pToPB(p)
	if err != nil {
		return nil, err
	}
	return pb, nil
}

func (r *productRepo) CreateProduct(ctx context.Context, p *biz.Product) (int64, error) {
	sqlStr := "INSERT INTO product(name, description, price) VALUES(?, ?, ?)"
	result, err := r.data.db.Exec(sqlStr, p.Name, p.Description, p.Price)
	if err != nil {
		return 0, errors.Wrapf(err, "sql: %s", sqlStr)
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return insertID, nil
}

func (r *productRepo) UpdateProduct(ctx context.Context, id int64, p *biz.Product) error {
	sqlStr := "UPDATE product SET description = ?, price= ? WHERE id = ?"
	_, err := r.data.db.Exec(sqlStr, id, p.Description, p.Price)
	if err != nil {
		return errors.Wrapf(err, "sql: %s", sqlStr)
	}
	return nil
}

func (r *productRepo) DeleteProduct(ctx context.Context, id int64) error {
	sqlStr := "DELETE FROM product WHERE id = ?"
	_, err := r.data.db.Exec(sqlStr, id)
	if err != nil {
		return errors.Wrapf(err, "sql: %s", sqlStr)
	}
	return nil
}
