package main

import (
	"log"

	"github.com/wlrudi19/elastic-engine/infrastructure/middlewares"
)

func main() {
	jwtGenerate := middlewares.GenerateAccessToken
	jwtValidation := middlewares.ValidateToken
	id := 7
	email := "joni@gmail.com"
	generateNih, err := jwtGenerate(id, email)

	if err != nil {
		log.Printf("ini error %v", err)
	}

	log.Printf("ini token %s", generateNih)

	claims, err := jwtValidation(generateNih)

	if err != nil {
		log.Println("token invalid")
	}

	email, ok := claims["Email"].(string)

	if !ok {
		log.Printf("email kosong: %s", email)
	}
}
