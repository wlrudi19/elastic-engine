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

type ProductResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Amount      string `json:"amount"`
	Stok        int    `json:"stok"`
}
