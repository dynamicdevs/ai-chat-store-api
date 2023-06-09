package product

type Product struct {
	Id        int       `json:"id"`
	Sku       string    `json:"sku"`
	Name      string    `json:"name"`
	Embedding []float32 `json:"embedding"`
}
