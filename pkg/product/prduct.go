package product

type Product struct {
	Id        int       `json:"id"`
	Sku       string    `json:"sku"`
	Name      string    `json:"name"`
	UrlPath   string    `json:"url_path"`
	Price     string    `json:"price"`
	Embedding []float32 `json:"embedding"`
}

type ProdutDetailes struct {
	Id         int      `json:"id"`
	Sku        string   `json:"sku"`
	Name       string   `json:"name"`
	UrlPath    string   `json:"url_path"`
	Price      string   `json:"price"`
	Attributes []string `json:"attributes"`
}
