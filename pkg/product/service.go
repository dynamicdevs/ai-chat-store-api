package product

import (
	"context"
	"fmt"

	"github.com/Abraxas-365/commerce-chat/pkg/attribute"
)

type Service interface {
	GetBySku(ctx context.Context, sku string) (*ProdutDetailes, error)
	OtherSimilars(ctx context.Context, sku string, embedding []float32, limit int) ([]ProdutDetailes, error)
	GetByEmbedding(ctx context.Context, embedding []float32, limit int) ([]ProdutDetailes, error)
}

type service struct {
	prepo Repository
	arepo attribute.Repository
}

func NewService(pr Repository, ar attribute.Repository) Service {
	return &service{
		pr,
		ar,
	}
}

func (s *service) GetBySku(ctx context.Context, sku string) (*ProdutDetailes, error) {

	attributesCh := make(chan []string)
	errorCh := make(chan error)

	go func() {
		attributesArray, err := s.arepo.GetBySKU(ctx, sku)
		if err != nil {
			errorCh <- err
			return
		}
		attributes := []string{}
		for _, attribute := range attributesArray {
			attributes = append(attributes, attribute.Information)
		}
		attributesCh <- attributes
	}()

	product, err := s.prepo.GetBySku(ctx, sku)
	if err != nil {
		return nil, err
	}

	attributes, ok := <-attributesCh
	if !ok {
		err := <-errorCh
		return nil, err
	}
	return &ProdutDetailes{
		Id:         product.Id,
		Sku:        product.Sku,
		Name:       product.Name,
		Attributes: attributes,
	}, nil
}

func (s *service) GetByEmbedding(ctx context.Context, embedding []float32, limit int) ([]ProdutDetailes, error) {
	products, err := s.prepo.MostSimilarVectors(ctx, embedding, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve most similar vectors: %w", err)
	}
	productCount := len(products)
	productsIds := make([]int, 0, productCount)
	for _, product := range products {
		productsIds = append(productsIds, product.Id)
	}
	attributesMap, err := s.arepo.GetByProducts(ctx, productsIds)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve attributes by products: %w", err)
	}
	productsDetails := make([]ProdutDetailes, 0, productCount)
	for _, product := range products {
		productAttributes, found := attributesMap[product.Id]
		if !found {
			continue
		}
		productsDetails = append(productsDetails, ProdutDetailes{
			Id:         product.Id,
			Sku:        product.Sku,
			Name:       product.Name,
			Price:      product.Price,
			UrlPath:    product.UrlPath,
			Attributes: productAttributes,
		})
	}
	return productsDetails, nil
}

func (s *service) OtherSimilars(ctx context.Context, sku string, embedding []float32, limit int) ([]ProdutDetailes, error) {
	products, err := s.prepo.MostSimilarVectorsExeptProductBySku(ctx, embedding, limit, sku)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve most similar vectors: %w", err)
	}

	productCount := len(products)
	productsIds := make([]int, 0, productCount)
	for _, product := range products {
		productsIds = append(productsIds, product.Id)
	}

	attributesMap, err := s.arepo.GetByProducts(ctx, productsIds)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve attributes by products: %w", err)
	}

	productsDetails := make([]ProdutDetailes, 0, productCount)
	for _, product := range products {
		productAttributes, found := attributesMap[product.Id]
		if !found {
			continue
		}
		productDetail := ProdutDetailes{
			Id:         product.Id,
			Sku:        product.Sku,
			Name:       product.Name,
			Attributes: productAttributes,
		}
		productsDetails = append(productsDetails, productDetail)
	}

	return productsDetails, nil
}
