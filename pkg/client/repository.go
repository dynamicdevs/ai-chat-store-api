package client

import "context"

type Repository interface {
	Save(ctx context.Context, c *Client) (int, error)
	GetById(ctx context.Context, id int) (*Client, error)
}
