package pgvector

import (
	"context"

	"github.com/Abraxas-365/commerce-chat/pkg/product"
	"github.com/jackc/pgx/v4/pgxpool"
)

type productRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) product.Repository {
	return &productRepository{pool: pool}
}

func (r *productRepository) Save(ctx context.Context, p *product.Product) (int, error) {
	var id int
	query := `INSERT INTO "public"."product" (name, sku) VALUES ($1, $2) RETURNING id`
	err := r.pool.QueryRow(ctx, query, p.Name, p.Sku).Scan(&id)
	return id, err
}

func (r *productRepository) GetBySku(ctx context.Context, sku string) (*product.Product, error) {
	product := &product.Product{}

	query := `SELECT id, sku, name FROM product WHERE sku = $1`

	err := r.pool.QueryRow(ctx, query, sku).Scan(&product.Id, &product.Sku, &product.Name)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *productRepository) ProductExistsBySku(ctx context.Context, sku string) (int, bool, error) {
	checkQuery := `SELECT id FROM "public"."product" WHERE sku = $1`
	var id int
	err := r.pool.QueryRow(ctx, checkQuery, sku).Scan(&id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			// Product does not exist
			return 0, false, nil
		}
		// Some other error occurred
		return 0, false, err
	}
	return id, true, nil
}
