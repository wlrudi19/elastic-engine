package model

func ToProductResponse(product Product) ProductResponse {
	return ProductResponse{
		Id:          product.Id,
		Name:        product.Name,
		Description: product.Description,
		Amount:      product.Amount,
		Stok:        product.Stok,
	}
}
