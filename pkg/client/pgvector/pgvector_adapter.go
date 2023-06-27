package pgvector

import (
	"context"

	"github.com/Abraxas-365/commerce-chat/pkg/client"
	"github.com/jackc/pgx/v4/pgxpool"
)

type clientRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) client.Repository {
	return &clientRepository{pool: pool}
}

func (r *clientRepository) Save(ctx context.Context, c *client.Client) (int, error) {
	var id int
	query := `INSERT INTO "public"."client" (name,system_prompt) VALUES ($1) RETURNING id`
	err := r.pool.QueryRow(ctx, query, c.Name).Scan(&id)
	return id, err
}

func (r *clientRepository) GetById(ctx context.Context, id int) (*client.Client, error) {
	client := &client.Client{}

	query := `SELECT id, name, system_prompt FROM product WHERE id = $1`

	err := r.pool.QueryRow(ctx, query, id).Scan(&client.Id, &client.Name, &client.SystemPromt)
	if err != nil {
		return nil, err
	}
	return client, nil
}
