package indexer

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/Abraxas-365/commerce-chat/internal/database"
	"github.com/Abraxas-365/commerce-chat/pkg/attribute"
	attributepg "github.com/Abraxas-365/commerce-chat/pkg/attribute/pgvector"
	"github.com/Abraxas-365/commerce-chat/pkg/openia"
	"github.com/Abraxas-365/commerce-chat/pkg/product"
	productpg "github.com/Abraxas-365/commerce-chat/pkg/product/pgvector"
)

type Indexer struct {
	db     *database.Connection
	openia *openia.Openia
}

type ProductAttribute struct {
	Product   product.Product       `json:"product"`
	Attribute []attribute.Attribute `json:"attribute"`
}

func New(db *database.Connection, openia *openia.Openia) *Indexer {
	return &Indexer{
		db,
		openia,
	}
}

func (i *Indexer) Index(csv string) error {
	ctx := context.Background()
	productdb := productpg.New(i.db.Pool)
	attributedb := attributepg.New(i.db.Pool)

	productAttributes, err := i.ReadCsv(csv)
	if err != nil {
		return err
	}
	fmt.Println(productAttributes[0].Product.Name)
	fmt.Println(productAttributes[0].Product.Sku)
	fmt.Println(productAttributes[0].Product.UrlPath)
	fmt.Println(productAttributes[0].Product.Price)

	var wg sync.WaitGroup
	sem := make(chan struct{}, 5)

	for _, productAttribute := range productAttributes {
		wg.Add(1)
		sem <- struct{}{}

		go func(productAttribute ProductAttribute) {
			defer wg.Done()
			defer func() { <-sem }()

			id, exist, err := productdb.ProductExistsBySku(ctx, productAttribute.Product.Sku)
			if err != nil {
				fmt.Println("Error saving checking product:", err)
				return
			}
			var productId int
			if exist {
				productId = id
			} else {
				nameEmbedding, err := i.openia.GenerateEmbedding(productAttribute.Product.Name)
				if err != nil {
					fmt.Println("Error generating product:", err)
					return
				}
				productAttribute.Product.Embedding = nameEmbedding
				productId, err = productdb.Save(ctx, &productAttribute.Product)
				if err != nil {
					fmt.Println("Error saving product:", err)
					return
				}
			}

			for _, attribute := range productAttribute.Attribute {
				attribute.Product = productId

				id, exist, err := attributedb.CheckAttributeExists(ctx, attribute.Information)
				if err != nil {
					fmt.Println("Error checking attribute existence:", err)
					return
				}

				if exist {
					err := attributedb.AssociateAttributeWithProduct(ctx, productId, id)
					if err != nil {
						fmt.Println("Error associating attribute with product:", err)
					}
					continue
				}

				embedding, err := i.openia.GenerateEmbedding(attribute.Information)
				if err != nil {
					fmt.Println("Error generating embedding:", err, productId)
					continue
				}
				attribute.Embedding = embedding

				id, err = attributedb.Save(ctx, &attribute)
				if err != nil {
					fmt.Println("Error saving attribute:", err)
					return
				}

				err = attributedb.AssociateAttributeWithProduct(ctx, productId, id)
				if err != nil {
					fmt.Println("Error associating attribute with product:", err)
				}
			}
		}(productAttribute)
	}

	wg.Wait()

	return nil
}

func (i *Indexer) ReadCsv(filename string) ([]ProductAttribute, error) {
	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		return []ProductAttribute{}, err
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return []ProductAttribute{}, err
	}

	var productAttributes []ProductAttribute

	// Skipping the header line
	for idx, line := range lines {
		if idx == 0 {
			continue
		}

		product := product.Product{
			Sku:     line[1],
			Name:    line[0],
			UrlPath: line[3],
			Price:   line[4],
		}

		var attributes []attribute.Attribute
		for _, attr := range strings.Split(line[2], "\n") {
			attribute := attribute.Attribute{Information: attr}
			attributes = append(attributes, attribute)
		}
		attributes = append(attributes, attribute.Attribute{Information: product.Name})

		productAttribute := ProductAttribute{
			Product:   product,
			Attribute: attributes,
		}
		productAttributes = append(productAttributes, productAttribute)
	}

	return productAttributes, nil
}
