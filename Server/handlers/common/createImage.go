package handlers

import (
	"MegaMedia/internal/config"
	"context"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateImage creates a new image object in MongoDB and uploads the file
func CreateImage(w http.ResponseWriter, r *http.Request, databaseName string, collectionName string, staticRoute string) error {
	// Parse the incoming multipart form (image upload)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return fmt.Errorf("failed to parse form: %v", err)
	}

	// Get the file from the form (with the key "image")
	file, handler, err := r.FormFile("image")
	if err != nil {
		return fmt.Errorf("failed to retrieve file: %v", err)
	}
	defer file.Close()

	// Get the image name from the form
	imageName := r.FormValue("name")
	if imageName == "" {
		return fmt.Errorf("image name is required")
	}

	// Get the file extension (e.g., ".jpg", ".png")
	fileExtension := path.Ext(handler.Filename)
	if fileExtension == "" {
		return fmt.Errorf("file has no extension")
	}

	// Create the directory where we will store the uploaded image
	fileDir := path.Join("./static", staticRoute)
	if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Generate a unique filename, adding timestamp only if file already exists
	originalFileName := handler.Filename
	fileNameWithoutExt := strings.TrimSuffix(originalFileName, fileExtension)
	fileName := fmt.Sprintf("%s%s", fileNameWithoutExt, fileExtension)
	filePath := path.Join(fileDir, fileName)

	// Check if file already exists, if so add timestamp
	if _, err := os.Stat(filePath); err == nil {
		// File exists, add timestamp
		timestamp := time.Now().Format("20060102_150405")
		fileName = fmt.Sprintf("%s_%s%s", fileNameWithoutExt, timestamp, fileExtension)
		filePath = path.Join(fileDir, fileName)
	}

	// Create the file on disk for the new image
	destFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer destFile.Close()

	// Copy the uploaded file into the newly created file
	if _, err = io.Copy(destFile, file); err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	// Get image dimensions
	file.Seek(0, 0) // Reset file pointer
	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return fmt.Errorf("failed to decode image: %v", err)
	}

	// Create MongoDB client
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.Cfg.MongoDBURI))
	if err != nil {
		return fmt.Errorf("database connection error: %v", err)
	}
	defer client.Disconnect(context.TODO())

	// Access database and collection
	collection := client.Database(databaseName).Collection(collectionName)

	// Create new image document
	newImage := Image{
		ID:       primitive.NewObjectID(),
		Name:     imageName,
		URL:      fileName,
		HeightPX: img.Height,
		WidthPX:  img.Width,
		Previous: []string{},
	}

	// Insert the new image into MongoDB
	_, err = collection.InsertOne(context.TODO(), newImage)
	if err != nil {
		// If database insert fails, clean up the uploaded file
		os.Remove(filePath)
		return fmt.Errorf("failed to insert image into database: %v", err)
	}

	// Success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Image created successfully: %s", imageName),
		"image":   newImage,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %v", err)
	}

	w.Write(jsonResponse)
	return nil
}
