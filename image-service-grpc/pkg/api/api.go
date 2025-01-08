package api

import (
	"io"
	"net/http"
	"os"
	"strconv"

	pb "github.com/vismaml/hiring/image-service-grpc/proto"
	"google.golang.org/grpc"
)

func ImageEndpoint(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new Image (message) instance
	req := &pb.Image{
		ImageData:   getImageData(r.Header.Get("Content-Type"), body),
		Scale:       os.Getenv("SCALE") == "true",
		Grayscale:   os.Getenv("GRAYSCALE") == "true",
		ContentType: r.Header.Get("Content-Type"),
		Width:       parseEnvToInt32("WIDTH"),
		Height:      parseEnvToInt32("HEIGHT"),
		Url:         getUrl(r.Header.Get("Content-Type"), body),
	}

	// Connect to the gRPC server
	conn, err := grpc.Dial("image-imgservice-svc:50051", grpc.WithInsecure())
	if err != nil {
		http.Error(w, "Failed to connect to gRPC server: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// Create a client for the ImageProcessingService
	service := pb.NewImageProcessingServiceClient(conn)

	// Call the ProcessImage method
	ctx := r.Context()
	processedImageResponse, err := service.ProcessImage(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(processedImageResponse.ProcessedImage)
}

// Converts the environment variable to int32
func parseEnvToInt32(key string) int32 {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return 0
	}
	return int32(value)
}

// If the content type is not text/plain, the body an image
func getImageData(contentType string, body []byte) []byte {
	if contentType != "text/plain" {
		return body
	}
	return nil
}

// If the content type is text/plain, the body is an URL
func getUrl(contentType string, body []byte) string {
	if contentType == "text/plain" {
		return string(body)
	}
	return ""
}
