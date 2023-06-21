// Package product provides an interface and types for managing products in an e-commerce system.
//
// The package defines a Repository interface that abstracts the storage of products.
// The Repository interface includes methods for saving products, retrieving products by SKU,
// checking if a product exists by SKU, and finding the most similar products based on embeddings.
//
// The Product struct is used to represent a product within the system, and includes fields
// such as ID, SKU, Name, and an Embedding that might be used for similarity searches.
//
// Example:
//
//	repo := // ... create or inject an implementation of the Repository interface ...
//
//	// Saving a product
//	product := &product.Product{
//	    Sku: "123",
//	    Name: "Sample Product",
//	    Embedding: []float32{...},
//	}
//	_, err := repo.Save(context.Background(), product)
//
//	// Retrieving a product by SKU
//	retrievedProduct, err := repo.GetBySku(context.Background(), "123")
//
//	// Finding similar products
//	similarProducts, err := repo.MostSimilarVectors(context.Background(), product.Embedding, 5)
package product
