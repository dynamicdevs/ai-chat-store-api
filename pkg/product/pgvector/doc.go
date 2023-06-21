// Package pgvector provides an implementation of the product.Repository interface using PostgreSQL as the backend.
//
// It leverages the pgvector-go library for working with vector data types in PostgreSQL and utilizes the
// pgx/v4/pgxpool package for efficient connection pooling and interaction with the PostgreSQL database.
//
// The productRepository struct is used internally to interact with the database, and implements the product.Repository
// interface defined in the product package.
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
//	// Initialize the productRepository
//	repo := pgvector.New(pool)
//
//	// Saving a product
//	product := &product.Product{
//	    Sku: "123",
//	    Name: "Sample Product",
//	    Embedding: []float32{...},
//	}
//	_, err = repo.Save(context.Background(), product)
//
//	// Retrieving a product by SKU
//	retrievedProduct, err := repo.GetBySku(context.Background(), "123")
//
//	// Finding similar products
//	similarProducts, err := repo.MostSimilarVectors(context.Background(), product.Embedding, 5)
package pgvector
