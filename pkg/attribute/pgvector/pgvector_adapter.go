package pgvector

import (
	"context"
	"fmt"

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

func (r *attributeRepository) Save(ctx context.Context, a *attribute.Attribute) (int, error) {
	query := `INSERT INTO "public"."attribute" (information, embedding) VALUES ($1, $2) RETURNING id`
	var attributeID int
	err := r.pool.QueryRow(ctx, query, a.Information, pgvector.NewVector(a.Embedding)).Scan(&attributeID)
	if err != nil {
		return 0, err
	}

	return attributeID, nil
}

func (r *attributeRepository) AssociateAttributeWithProduct(ctx context.Context, productID int, attributeID int) error {
	query := `INSERT INTO "public"."product_attribute" (product_id, attribute_id) VALUES ($1, $2)`
	_, err := r.pool.Exec(ctx, query, productID, attributeID)
	return err
}

func (r *attributeRepository) CheckAttributeExists(ctx context.Context, information string) (int, bool, error) {
	checkQuery := `SELECT id FROM "public"."attribute" WHERE information = $1`
	var existingAttributeID int
	err := r.pool.QueryRow(ctx, checkQuery, information).Scan(&existingAttributeID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return 0, false, nil
		}
		fmt.Println("ERRR", err)
		return 0, false, err
	}

	return existingAttributeID, true, nil
}

func (r *attributeRepository) MostSimilarVectors(ctx context.Context, embedding []float32, limit int) ([]product.Product, error) {
	query := `
    SELECT p.id, p.sku, p.name
    FROM "public"."product" p
    JOIN (
        SELECT pa.product_id, MIN(a.embedding <-> $1) as distance
        FROM "public"."attribute" a
        JOIN "public"."product_attribute" pa ON a.id = pa.attribute_id
        GROUP BY pa.product_id
    ) subquery ON p.id = subquery.product_id
    ORDER BY subquery.distance
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
