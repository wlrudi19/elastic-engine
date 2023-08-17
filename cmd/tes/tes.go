package main

import (
	"log"

	"github.com/wlrudi19/elastic-engine/infrastructure/jwt"
)

func main() {
	jwtGet := jwt.NewJWT()
	id := 7
	email := "joni@gmail.com"
	generateNih, err := jwtGet.GenerateAccessToken(id, email)

	if err != nil {
		log.Printf("ini error %v", err)
	}

	log.Printf("ini token %s", generateNih)
}
