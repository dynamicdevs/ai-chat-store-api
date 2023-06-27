package admin

import (
	"github.com/Abraxas-365/commerce-chat/internal/database"
	"github.com/Abraxas-365/commerce-chat/pkg/client"
	clientpg "github.com/Abraxas-365/commerce-chat/pkg/client/pgvector"
)

type Config struct {
	Db *database.Connection
}

type Admin struct {
	client client.Repository
}

func New(c Config) *Admin {
	clientdb := clientpg.New(c.Db.Pool)
	return &Admin{
		clientdb,
	}
}
