package product

import (
	"context"
)

type Repository interface {
	Save(ctx context.Context, p *Product) (int, error)
	GetBySku(ctx context.Context, sku int) (*Product, error)
}
