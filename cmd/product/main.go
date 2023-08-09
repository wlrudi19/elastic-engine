package main

import (
	"fmt"
	"net/http"

	"github.com/wlrudi19/elastic-engine/config"
)

func main() {
	dbConfig := config.NewLoad()
	config := dbConfig.LoadConfig()
	db, err := dbConfig.TestConnection(config.Database)

	if err != nil {
		panic(err)
	}

	fmt.Println("ELASTIC ENGINE PROJECT - RUDI LESMANA")

	defer db.Close()

	port := "7654" // Change this to your desired port
	fmt.Printf("Server started on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
