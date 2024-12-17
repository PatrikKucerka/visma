package main

import (
	"log"
	"net/http"

	"github.com/vismaml/hiring/image-service/pkg/api"
)

func main() {
	http.HandleFunc("/image", api.ImageEndpoint)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
