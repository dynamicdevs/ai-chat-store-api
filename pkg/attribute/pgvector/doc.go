// Package pgvector provides an implementation of the attribute.Repository interface
// using PostgreSQL as the backend. It leverages the pgvector-go library for working with
// vector data types in PostgreSQL and utilizes the pgx/v4/pgxpool package for efficient
// connection pooling and interaction with the PostgreSQL database.
//
// This package is primarily meant to handle product attributes within an e-commerce context
// by storing and retrieving attribute data in/from a PostgreSQL database.
//
// The attributeRepository struct provides the actual implementation by holding a connection pool
// to the PostgreSQL database.
//
// Example:
//
//	// Create a new connection pool
//	connString := "postgresql://username:password@localhost/dbname"
//	pool, err := pgxpool.Connect(context.Background(), connString)
//	if err != nil {
//	    log.Fatal("Unable to connect to database:", err)
//	}
//	defer pool.Close()
//
//	// Initialize the attributeRepository
//	repo := pgvector.New(pool)
//
//	// Saving an attribute
//	attr := &attribute.Attribute{
//	    Information: "Brand: APPLE",
//	    Embedding:   []float32{...},
//	}
//	_, err = repo.Save(context.Background(), attr)
//
//	// Retrieving attributes by SKU
//	attributes, err := repo.GetBySKU(context.Background(), "SKU123")
package pgvector
