// Package attribute provides the definitions and interface required for managing attributes
// associated with products in an e-commerce context.
//
// The package defines the Attribute struct which represents an attribute of a product.
// Attributes have an ID, information (e.g. "Brand: APPLE"), and an embedding for similarity
// search.
//
// The Repository interface specifies the methods required for attribute management, such
// as saving attributes, finding attributes with similar embeddings, and associating attributes
// with products.
//
// Example:
//
//	attr := &attribute.Attribute{
//	    Information: "Brand: APPLE",
//	    Embedding:   []float32{...},
//	}
//
//	repo := ... // Implementation of attribute.Repository
//	_, err = repo.Save(context.Background(), attr)
//
// Consumers of this package should implement the Repository interface to provide concrete
// implementations for data storage and retrieval.
package attribute
