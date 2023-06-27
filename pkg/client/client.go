package client

type Client struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	SystemPromt string `json:"system_promt"`
}
