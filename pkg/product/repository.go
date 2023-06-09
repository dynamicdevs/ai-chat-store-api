package product

import (
	"context"
)

type Repository interface {
	Save(ctx context.Context, p *Product) (int, error)
	GetBySku(ctx context.Context, sku string) (*Product, error)
	ProductExistsBySku(ctx context.Context, sku string) (int, bool, error)
	MostSimilarVectors(ctx context.Context, embedding []float32, limit int) ([]Product, error)
}
