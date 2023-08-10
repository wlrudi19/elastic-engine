package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/wlrudi19/elastic-engine/app/product/api"
	"github.com/wlrudi19/elastic-engine/app/product/repository"
	"github.com/wlrudi19/elastic-engine/app/product/service"
	"github.com/wlrudi19/elastic-engine/config"
	"github.com/wlrudi19/elastic-engine/helper"
	"github.com/wlrudi19/elastic-engine/infrastructure/middleware"
)

func main() {
	dbConfig := config.NewLoad()
	config := dbConfig.LoadConfig()
	db, err := dbConfig.TestConnection(config.Database)

	helper.PanicIfError(err)

	validate := validator.New()
	productRepository := repository.NewProductRepository()
	productService := service.NewProductService(productRepository, db, validate)
	productController := api.NewProductController(productService)
	router := api.NewProductRouter(productController)

	fmt.Println("ELASTIC ENGINE PROJECT - RUDI LESMANA")
	defer db.Close()

	server := http.Server{
		Addr:    "localhost:7654",
		Handler: middleware.NewAuthMiddleware(router),
	}

	fmt.Println("Starting server on port 7654...")

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// port := "7654" // Change this to your desired port
	// fmt.Printf("Server started on port %s\n", port)
	// http.ListenAndServe(":"+port, nil)
}

//TEST CONNECTION
// //test insert
// tx, err := db.Begin()
// if err != nil {
// 	fmt.Printf("Failed to begin transaction: %v\n", err)
// 	return
// }
// defer tx.Rollback() // Jangan lupa rollback jika terjadi error atau commit jika berhasil

// createProductRepo := repository.NewProductRepository()

// // Membuat data product yang akan diinsert
// product := model.Product{
// 	Name:        "Sample Product 2",
// 	Description: "Sample Description 2",
// 	Amount:      "",
// 	Stok:        5,
// }

// // Memanggil fungsi CreateProduct dengan objek tx
// createdProduct, err := createProductRepo.CreateProduct(context.Background(), tx, product)
// if err != nil {
// 	fmt.Printf("Failed to create product: %v\n", err)
// } else {
// 	fmt.Printf("Product created successfully with ID: %d\n", createdProduct.Id)
// 	// Jika tidak terjadi error, commit transaksi
// 	if err := tx.Commit(); err != nil {
// 		fmt.Printf("Failed to commit transaction: %v\n", err)
// 	}
// }
