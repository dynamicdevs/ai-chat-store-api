package attribute

import (
	"context"
)

type Repository interface {
	Save(ctx context.Context, a *Attribute) (int, error)
	GetByProducts(ctx context.Context, ids []int) (map[int][]Attribute, error)
	CheckAttributeExists(ctx context.Context, information string) (int, bool, error)
	AssociateAttributeWithProduct(ctx context.Context, productID int, attributeID int) error
	GetBySKU(ctx context.Context, sku string) ([]Attribute, error)
}
