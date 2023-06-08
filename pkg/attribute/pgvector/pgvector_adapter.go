package pgvector

import (
	"context"

	"github.com/Abraxas-365/commerce-chat/pkg/attribute"
	"github.com/Abraxas-365/commerce-chat/pkg/product"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"
	"github.com/pgvector/pgvector-go"
)

type attributeRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) attribute.Repository {
	return &attributeRepository{pool: pool}
}

func (r *attributeRepository) Save(ctx context.Context, a *attribute.Attribute) error {
	// Insert into attribute table
	query := `INSERT INTO "public"."attribute" (information, embedding) VALUES ($1, $2) RETURNING id`
	var attributeID int
	err := r.pool.QueryRow(ctx, query, a.Information, pgvector.NewVector(a.Embedding)).Scan(&attributeID)
	if err != nil {
		return err
	}

	// Associate the attribute with the product in the product_attribute table
	assocQuery := `INSERT INTO "public"."product_attribute" (product_id, attribute_id) VALUES ($1, $2)`
	_, err = r.pool.Exec(ctx, assocQuery, a.Product, attributeID)
	return err
}

func (r *attributeRepository) MostSimilarVectors(ctx context.Context, embedding []float32, limit int) ([]product.Product, error) {
	query := `
	SELECT DISTINCT ON (p.id) p.id, p.sku, p.name
	FROM "public"."attribute" a
	JOIN "public"."product_attribute" pa ON a.id = pa.attribute_id
	JOIN "public"."product" p ON pa.product_id = p.id
	ORDER BY p.id, a.embedding <-> $1
	LIMIT $2;
	`

	rows, err := r.pool.Query(ctx, query, pgvector.NewVector(embedding), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []product.Product

	for rows.Next() {
		var p product.Product

		err := rows.Scan(&p.Id, &p.Sku, &p.Name)
		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return products, nil
}

func (r *attributeRepository) GetByProducts(ctx context.Context, ids []int) (map[int][]attribute.Attribute, error) {
	query := `
	SELECT a.id, pa.product_id, a.information
	FROM "public"."attribute" a
	JOIN "public"."product_attribute" pa ON a.id = pa.attribute_id
	WHERE pa.product_id = ANY($1);
	`

	skuArray := pq.Array(ids)
	rows, err := r.pool.Query(ctx, query, skuArray)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	attributeMap := make(map[int][]attribute.Attribute)
	for rows.Next() {
		var a attribute.Attribute
		err := rows.Scan(&a.Id, &a.Product, &a.Information)
		if err != nil {
			return nil, err
		}
		attributeMap[a.Product] = append(attributeMap[a.Product], a)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return attributeMap, nil
}
