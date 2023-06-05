package attribute

type Attribute struct {
	Id          int       `json:"id"`
	Product     int       `json:"product"`
	Information string    `json:"information"` //ejemplo "Marca:	APPLE"
	Embedding   []float32 `json:"embedding"`
}
