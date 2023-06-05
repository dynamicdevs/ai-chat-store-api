package indexer

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

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

type productAttribute struct {
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

	for _, productAttribute := range productAttributes {
		productId, err := productdb.Save(ctx, &productAttribute.Product)
		if err != nil {
			return err
		}

		for _, attribute := range productAttribute.Attribute {
			attribute.Product = productId
			embedding, err := i.openia.GenerateEmbedding(attribute.Information)
			if err != nil {
				//TODO: should handle retries
				return err
			}
			attribute.Embedding = embedding
			if err := attributedb.Save(ctx, &attribute); err != nil {
				//TODO: should handle retries
				return err
			}
		}

	}

	return nil
}

func (i *Indexer) ReadCsv(filename string) ([]productAttribute, error) {
	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		return []productAttribute{}, err
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return []productAttribute{}, err
	}

	var productAttributes []productAttribute

	for _, line := range lines {
		product := product.Product{
			Sku:  line[1],
			Name: line[0],
		}
		fmt.Println(product)

		var attributes []attribute.Attribute
		// Parse attributes
		for _, attr := range strings.Split(line[2], ",") {
			attribute := attribute.Attribute{Information: attr}
			attributes = append(attributes, attribute)
		}

		productAttribute := productAttribute{
			Product:   product,
			Attribute: attributes,
		}
		productAttributes = append(productAttributes, productAttribute)
	}

	return productAttributes, nil
}
