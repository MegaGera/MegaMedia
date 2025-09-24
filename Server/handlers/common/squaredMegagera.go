package handlers

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

// SquaredMegageraImageHandler handles image squared feature for MegaGera
// fileName already includes the extension (e.g., "image.png")
func SquaredMegageraImageHandler(w http.ResponseWriter, r *http.Request, staticRoute string, fileName string) (previousFileName string, err error) {
	// Construct the file path to retrieve the original image
	fileDir := path.Join("./static", staticRoute)
	newFilePath := path.Join(fileDir, fileName)

	// Read the existing file
	file, err := os.Open(newFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open existing file: %v", err)
	}
	defer file.Close()

	// Read the file content into a variable
	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to read existing file: %v", err)
	}

	// Process the image to make it square
	bgColor := color.Transparent
	squaredImage := squareImage(img, bgColor)

	// File exists, rename it with "_old_{timestamp}" suffix
	timestamp := time.Now().Format("20060102_150405")

	// Extract filename without extension and add timestamp
	fileNameWithoutExt := strings.TrimSuffix(fileName, path.Ext(fileName))
	newFileName := fmt.Sprintf("%s_old_%s%s", fileNameWithoutExt, timestamp, path.Ext(fileName))
	oldFilePath := path.Join(fileDir, newFileName)

	// Rename the existing file
	log.Printf("Renaming existing file: %s -> %s\n", newFilePath, oldFilePath)
	if err := os.Rename(newFilePath, oldFilePath); err != nil {
		return "", fmt.Errorf("failed to rename existing file: %v", err)
	}

	// Create the file on disk for the new image
	log.Printf("Saving new file to: %s\n", newFilePath)
	destFile, err := os.Create(newFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer destFile.Close()

	if err = png.Encode(destFile, squaredImage); err != nil {
		return "", fmt.Errorf("failed to save / encode file: %v", err)
	}

	log.Printf("File squared successfully for %s: %s\n", fileName, newFilePath)
	return newFileName, nil
}
