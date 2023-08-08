package main

import (
	"fmt"

	"github.com/wlrudi19/elastic-engine/config"
	"github.com/wlrudi19/elastic-engine/test"
)

func main() {
	config := config.LoadConfig()
	db, err := test.TestConnection(config.Database)
	if err != nil {
		panic(err)
	}

	fmt.Println("ELASTIC ENGINE PROJECT - RUDI LESMANA")
	fmt.Println(db)
	defer db.Close()
}
