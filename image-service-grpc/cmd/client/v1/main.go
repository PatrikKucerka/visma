package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: client <input-file> <output-file>")
	}

	filePath := os.Args[1]
	outputPath := os.Args[2]

	content, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	response, err := http.Post("http://localhost:8080/image", "image/jpeg", content)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	fmt.Println(response.Status, response.Request.Method, response.Request.URL)
	fmt.Println("Body len:", len(body))

	out, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	_, err = out.Write(body)
	if err != nil {
		log.Fatal(err)
	}
}
