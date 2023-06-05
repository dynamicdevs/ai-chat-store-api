package attribute

import (
	"context"

	"github.com/Abraxas-365/commerce-chat/pkg/product"
)

type Repository interface {
	Save(ctx context.Context, a *Attribute) error
	MostSimilarVectors(ctx context.Context, embedding []float32, limit int) ([]Attribute, []product.Product, error)
	GetByProduct(ctx context.Context, sku int) ([]*Attribute, error)
}
