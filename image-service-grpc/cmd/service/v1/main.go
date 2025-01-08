package main

import (
	"log"
	"net/http"

	"github.com/lpernett/godotenv"
	"github.com/vismaml/hiring/image-service-grpc/pkg/api"
)

func main() {
	err := godotenv.Load("/app/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/image", api.ImageEndpoint)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
