package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatal("Usage: client <mode> <input> <output-file>")
	}

	mode := os.Args[1]
	input := os.Args[2]
	outputPath := os.Args[3]

	// Determine if input is an URL or an image based on the mode
	switch mode {
	case "url":
		// The input is a URL (send the URL as plain text)
		requestBody := []byte(input)
		response, err := http.Post("http://localhost:8080/image", "text/plain", bytes.NewReader(requestBody))
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(response.Status, response.Request.Method, response.Request.URL)
		fmt.Println("Body len:", len(body))

		// Save the response body to the output file
		out, err := os.Create(outputPath)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		_, err = out.Write(body)
		if err != nil {
			log.Fatal(err)
		}

	case "image":
		//The input is image
		content, err := os.Open(input)
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

	default:
		log.Fatal("Invalid mode. Use 'string' for URL or 'image' for image path.")
	}
}
