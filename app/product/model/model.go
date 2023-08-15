package model

type Product struct {
	Id          int
	Name        string
	Description string
	Amount      string
	Stok        int
}

type CreateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Amount      string `json:"amount"`
	Stok        int    `json:"stok"`
}

type FindProductRequest struct {
	Id int `json:"id"`
}

type FindProductResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Amount      string `json:"amount"`
	Stok        int    `json:"stok"`
}
