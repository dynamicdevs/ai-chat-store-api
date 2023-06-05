package pgvector

import (
	"context"

	"github.com/Abraxas-365/commerce-chat/pkg/attribute"
	"github.com/Abraxas-365/commerce-chat/pkg/product"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pgvector/pgvector-go"
)

type attributeRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) attribute.Repository {
	return &attributeRepository{pool: pool}
}

func (r *attributeRepository) Save(ctx context.Context, a *attribute.Attribute) error {
	query := `INSERT INTO "public"."attribute" (product, information, embedding) VALUES ($1, $2, $3)`
	_, err := r.pool.Exec(ctx, query, a.Product, a.Information, pgvector.NewVector(a.Embedding))
	return err
}

func (r *attributeRepository) MostSimilarVectors(ctx context.Context, embedding []float32, limit int) ([]attribute.Attribute, []product.Product, error) {
	query := `
	SELECT a.id, a.product, a.information, p.id, p.sku, p.name
	FROM "public"."attribute" a
	JOIN "public"."product" p ON a.product = p.id
	ORDER BY a.embedding <-> $1
	LIMIT $2;
	`

	rows, err := r.pool.Query(ctx, query, pgvector.NewVector(embedding), limit)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var attributes []attribute.Attribute
	var products []product.Product

	for rows.Next() {
		var a attribute.Attribute
		var p product.Product

		err := rows.Scan(&a.Id, &a.Product, &a.Information, &p.Id, &p.Sku, &p.Name)
		if err != nil {
			return nil, nil, err
		}

		attributes = append(attributes, a)
		products = append(products, p)
	}

	if rows.Err() != nil {
		return nil, nil, rows.Err()
	}

	return attributes, products, nil
}

func (r *attributeRepository) GetByProduct(ctx context.Context, sku int) ([]*attribute.Attribute, error) {
	query := `
	SELECT a.id, a.product, a.information, a.embedding
	FROM "public"."attribute" a
	JOIN "public"."product" p ON a.product = p.id
	WHERE p.sku = $1;
	`

	rows, err := r.pool.Query(ctx, query, sku)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attributes []*attribute.Attribute

	for rows.Next() {
		var a attribute.Attribute

		err := rows.Scan(&a.Id, &a.Product, &a.Information, &a.Embedding)
		if err != nil {
			return nil, err
		}

		attributes = append(attributes, &a)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return attributes, nil
}
