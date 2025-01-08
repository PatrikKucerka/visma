package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	pb "github.com/vismaml/hiring/image-service-grpc/proto"
	"golang.org/x/image/draw"
)

var (
	rdb *redis.Client
)

type ImageProcessingServer struct {
	pb.UnimplementedImageProcessingServiceServer
}

// Initialize the Redis client
func InitRedisClient() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis-svc:6379", // Redis service address
		Password: "",               // No password
		DB:       0,                // Default DB
	})
}

// ProcessImage processes the image based on the request parameters
func (s *ImageProcessingServer) ProcessImage(ctx context.Context, req *pb.Image) (*pb.ProcessImageResponse, error) {
	// If the request is a plain text, it is assumend to be an URL (otherwise it would be e.g. image/jpeg)
	if req.ContentType == "text/plain" {
		response, err := http.Get(req.Url) // Fetch the image from the URL
		if err != nil {
			return &pb.ProcessImageResponse{Status: "failed"}, err
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body) // Read the image data
		if err != nil {
			return &pb.ProcessImageResponse{Status: "failed"}, err
		}
		req.ImageData = body // Set the image data to the fetched image
	}

	// Generate the hash of the image data to create a unique key
	hash := sha256.New()
	hash.Write(req.ImageData)
	cacheKey := hex.EncodeToString(hash.Sum(nil)) // Create the cache key from the hash

	// Check if the image is already cached in Redis
	cachedImage, err := rdb.Get(ctx, cacheKey).Bytes()
	if err == nil {
		// Return the cached image if it exists
		return &pb.ProcessImageResponse{
			Status:         "cached",
			ProcessedImage: cachedImage,
		}, nil
	}
	// Decode the image data
	img, _, err := decodeImage(req.ImageData)
	if err != nil {
		return &pb.ProcessImageResponse{Status: "failed"}, err
	}
	//If the grayscale flag is set, convert the image to grayscale
	if req.Grayscale {
		img = toGrayscale(img)
	}
	// If the scale flag is set, resize the image
	if req.Scale {
		img = resize(img, int(req.Width), int(req.Height))
	}
	// Encode the processed image
	processedImage, err := encodeImage(img, req.ContentType)
	if err != nil {
		return &pb.ProcessImageResponse{Status: "failed"}, err
	}

	// Cache the processed image in Redis with an expiration time (e.g., 1 hour)
	err = rdb.Set(ctx, cacheKey, processedImage, time.Hour).Err()
	if err != nil {
		return &pb.ProcessImageResponse{Status: "failed"}, err
	}

	return &pb.ProcessImageResponse{
		Status:         "success",
		ProcessedImage: processedImage,
	}, nil
}

// decodeImage decodes the image data
func decodeImage(data []byte) (image.Image, string, error) {
	reader := bytes.NewReader(data)
	img, format, err := image.Decode(reader)
	return img, format, err
}

// encodeImage encodes the image data
func encodeImage(img image.Image, contentType string) ([]byte, error) {
	var buf bytes.Buffer
	var err error

	switch contentType {
	case "image/jpeg":
		err = jpeg.Encode(&buf, img, nil)
	case "image/png":
		err = png.Encode(&buf, img)
	default:
		err = jpeg.Encode(&buf, img, nil) // Default to JPEG
	}
	return buf.Bytes(), err
}

// toGrayscale converts the image to grayscale
func toGrayscale(img image.Image) image.Image {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			grayImg.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
		}
	}
	return grayImg
}

// Resize resizes the image to the desired resolution
func resize(oldImg image.Image, width, height int) image.Image {
	// Create a new blank image with the target resolution
	resizedImg := image.NewRGBA(image.Rect(0, 0, width, height))

	// Perform the rescaling
	draw.CatmullRom.Scale(resizedImg, resizedImg.Rect, oldImg, oldImg.Bounds(), draw.Over, nil)

	return resizedImg
}
