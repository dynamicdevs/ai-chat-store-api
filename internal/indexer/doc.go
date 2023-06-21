// Package indexer provides utilities for indexing products and their attributes into
// a database from a CSV file. It utilizes an AI model to generate embeddings for the product names
// and attributes, which can be used for similarity searches.
//
// The Indexer struct is the main component in this package. It uses a database connection and an AI model
// to index products and attributes.
//
// Example Usage:
//
//	// Initialize database connection and OpenIA AI model
//	db := database.NewConnection(...)
//	openIA := openia.New(...)
//
//	// Create an Indexer
//	indexer := indexer.New(db, openIA)
//
//	// Index products and attributes from a CSV file
//	err := indexer.Index("path_to_csv_file.csv")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// The CSV file should have the following format:
//
//	ProductName,ProductSKU,Attribute1\nAttribute2\n...
//
// The Indexer reads the CSV file, generates embeddings for product names and attributes,
// and saves them to the database along with their associations.
//
// The Index method utilizes concurrency to speed up the indexing process. It employs goroutines for
// processing multiple products simultaneously, and uses a semaphore to limit the number of concurrent
// goroutines.
//
// Note: Ensure that the database schema is compatible with the product and attribute structures
// used in this package.
package indexer
