package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

// UploadImageHandler handles image uploads and returns an error only if something goes wrong. Otherwise, it returns the new file name.
func UploadImageHandler(w http.ResponseWriter, r *http.Request, staticRoute string, fileName string) (previousFileName string, err error) {

	// Parse the incoming multipart form (image upload)
	if err = r.ParseMultipartForm(10 << 20); err != nil {
		return "", fmt.Errorf("failed to parse form: %v", err)
	}

	// Get the file from the form (with the key "image")
	file, handler, err := r.FormFile("image")
	if err != nil {
		return "", fmt.Errorf("failed to retrieve file: %v", err)
	}
	defer file.Close()

	// Get the file extension (e.g., ".jpg", ".png")
	fileExtension := path.Ext(handler.Filename)
	if fileExtension == "" {
		return "", fmt.Errorf("file has no extension")
	}

	// Create the directory where we will store the uploaded image
	fileDir := path.Join("./static", staticRoute)
	if err = os.MkdirAll(fileDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	// Construct the file path to store the new image
	newFilePath := path.Join(fileDir, fmt.Sprintf("%s%s", fileName, fileExtension))

	// Ensure the directory for the file exists (in case fileName contains subdirectories)
	if err = os.MkdirAll(path.Dir(newFilePath), os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create subdirectory: %v", err)
	}
	newFileName := ""

	// Check if the file already exists
	if _, err := os.Stat(newFilePath); err == nil {
		// File exists, rename it with "_old_{timestamp}" suffix
		timestamp := time.Now().Format("20060102_150405")
		newFileName = fmt.Sprintf("%s_old_%s%s", fileName, timestamp, fileExtension)
		oldFilePath := path.Join(fileDir, newFileName)

		// Rename the existing file
		log.Printf("Renaming existing file: %s -> %s\n", newFilePath, oldFilePath)
		if err := os.Rename(newFilePath, oldFilePath); err != nil {
			return "", fmt.Errorf("failed to rename existing file: %v", err)
		}
	}

	// Create the file on disk for the new image
	log.Printf("Saving new file to: %s\n", newFilePath)
	destFile, err := os.Create(newFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer destFile.Close()

	// Copy the uploaded file into the newly created file
	if _, err = io.Copy(destFile, file); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	log.Printf("File uploaded successfully for %s: %s\n", fileName, newFilePath)
	return newFileName, nil
}
