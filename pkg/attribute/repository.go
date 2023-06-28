package attribute

import (
	"context"

	"github.com/Abraxas-365/commerce-chat/pkg/product"
)

type Repository interface {
	Save(ctx context.Context, a *Attribute) (int, error)
	MostSimilarVectors(ctx context.Context, embedding []float32, limit int) ([]product.Product, error)
	GetByProducts(ctx context.Context, ids []int) (map[int][]Attribute, error)
	CheckAttributeExists(ctx context.Context, information string) (int, bool, error)
	AssociateAttributeWithProduct(ctx context.Context, productID int, attributeID int) error
	GetBySKU(ctx context.Context, sku string) ([]Attribute, error)
}
